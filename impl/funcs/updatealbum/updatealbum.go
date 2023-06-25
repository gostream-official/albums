package updatealbum

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gostream-official/albums/impl/inject"
	"github.com/gostream-official/albums/impl/models"
	"github.com/gostream-official/albums/pkg/api"
	"github.com/gostream-official/albums/pkg/arrays"
	"github.com/gostream-official/albums/pkg/marshal"
	"github.com/gostream-official/albums/pkg/parallel"
	"github.com/gostream-official/albums/pkg/store"
	"github.com/gostream-official/albums/pkg/store/query"
	"github.com/revx-official/output/log"

	"github.com/google/uuid"
)

// Description:
//
//	The request body for the update album endpoint.
type UpdateAlbumRequestBody struct {

	// The title of the album.
	Title string `json:"title,omitempty"`

	// The tracks contained in the album.
	TrackIDs []string `json:"trackIds,omitempty"`

	// Some album statistics.
	Stats UpdateAlbumStatsRequestBody `json:"stats,omitempty"`
}

// Description:
//
//	The request body for the track statistics.
type UpdateAlbumStatsRequestBody struct {

	// The popularity factor.
	Popularity float32 `json:"popularity,omitempty"`
}

// Description:
//
//	The error response body for the create track endpoint.
type UpdateAlbumErrorResponseBody struct {

	// The error message.
	Message string `json:"message"`
}

// Description:
//
//	Describes a validation error.
type UpdateAlbumValidationError struct {

	// The JSON field which is referenced by the error message.
	FieldRef string `json:"ref"`

	// The error message.
	ErrorMessage string `json:"error"`
}

// Description:
//
//	Describes a path parameter validation error.
type UpdateAlbumPathValidationError struct {

	// The JSON field which is referenced by the error message.
	PathRef string `json:"pathRef"`

	// The error message.
	ErrorMessage string `json:"error"`
}

// Description:
//
//	Attempts to cast the input object to the endpoint injector.
//	If this cast fails, we cannot proceed to process this request.
//
// Parameters:
//
//	object 	The injector object.
//
// Returns:
//
//	The injector if the cast is successful, an error otherwise.
func GetSafeInjector(object interface{}) (*inject.Injector, error) {
	injector, ok := object.(inject.Injector)

	if !ok {
		return nil, fmt.Errorf("updatealbum: failed to deduce injector")
	}

	return &injector, nil
}

// Description:
//
//	Unmarshals the request body for this endpoint.
//
// Parameters:
//
//	request The original request.
//
// Returns:
//
//	The unmarshalled request body, or an error when unmarshalling fails.
func ExtractRequestBody(request *api.APIRequest) (*UpdateAlbumRequestBody, error) {
	body := &UpdateAlbumRequestBody{}

	bytes := []byte(request.Body)
	err := json.Unmarshal(bytes, body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

// Description:
//
//	Gets and validates the id path parameter.
//
// Parameters:
//
//	request The http request.
//
// Returns:
//
//	The id path parameter.
//	A validatior error if the id is not a valid uuid.
func GetAndValidateID(request *api.APIRequest) (string, *UpdateAlbumPathValidationError) {
	id := request.PathParameters["id"]

	_, err := uuid.Parse(id)
	if err != nil {
		return "", &UpdateAlbumPathValidationError{
			PathRef:      ":id",
			ErrorMessage: "value is not a valid uuid",
		}
	}

	return id, nil
}

// Description:
//
//	Validates the request body for this endpoint.
//
// Parameters:
//
//	request The request body.
//
// Returns:
//
//	An error if the validation fails.
func ValidateRequestBody(request *UpdateAlbumRequestBody) *UpdateAlbumValidationError {

	if request.Title != "" {
		title := strings.TrimSpace(request.Title)

		if len(title) == 0 {
			return &UpdateAlbumValidationError{
				FieldRef:     "title",
				ErrorMessage: "value is not a valid uuid",
			}
		}
	}

	if len(request.TrackIDs) > 0 {
		trackIDs := arrays.Map(request.TrackIDs, func(artist string) string {
			return strings.TrimSpace(artist)
		})

		for _, trackID := range trackIDs {
			_, err := uuid.Parse(trackID)
			if err != nil {
				return &UpdateAlbumValidationError{
					FieldRef:     "trackIDs",
					ErrorMessage: "array contains invalid uuid",
				}
			}
		}
	}

	return ValidateAlbumStatsRequestBody(&request.Stats)
}

// Description:
//
//	Validates the statistics request body for this endpoint.
//
// Parameters:
//
//	request The statistics request body.
//
// Returns:
//
//	An error if the validation fails.
func ValidateAlbumStatsRequestBody(request *UpdateAlbumStatsRequestBody) *UpdateAlbumValidationError {
	return nil
}

// Description:
//
//	Searches an album with the given id in the database.
//
// Parameters:
//
//	store 	The store to search through.
//	id 		The id to search for.
//
// Returns:
//
//	The first matched album.
//	An error if the query fails.
func FindAlbumByID(store *store.MongoStore[models.AlbumInfo], id string) (*models.AlbumInfo, error) {
	filter := query.Filter{
		Root: query.FilterOperatorEq{
			Key:   "_id",
			Value: id,
		},
		Limit: 1,
	}

	items, err := store.FindItems(&filter)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("updatealbum: album not found")
	}

	return &items[0], nil
}

// Description:
//
//	Checks whether the given artist id exists in the mongo store.
//
// Parameters:
//
//	store 	The mongo store to search.
//	trackID The artist id to search.
//
// Returns:
//
//	An error, if the artist could not be found or an error,
//	if the database request failed, nothing if successful.
func CheckIfTrackExists(store *store.MongoStore[models.TrackInfo], trackID string) error {
	filter := query.Filter{
		Root: query.FilterOperatorEq{
			Key:   "_id",
			Value: trackID,
		},
		Limit: 1,
	}

	items, err := store.FindItems(&filter)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return fmt.Errorf("track not found")
	}

	return nil
}

// Description:
//
//	The router handler for track creation.
//
// Parameters:
//
//	request The incoming request.
//	object 	The injector. Contains injected dependencies.
//
// Returns:
//
//	An API response object.
func Handler(request *api.APIRequest, object interface{}) *api.APIResponse {
	context := parallel.NewContext()

	log.Infof("[%s] %s: %s", context.ID, request.Method, request.Path)
	log.Tracef("[%s] request: %s", context.ID, marshal.Quick(request))

	injector, err := GetSafeInjector(object)
	if err != nil {
		log.Warnf("[%s] failed to get endpoint injector: %s", context.ID, err)
		return &api.APIResponse{
			StatusCode: http.StatusInternalServerError,
		}
	}

	id, validationErr := GetAndValidateID(request)
	if validationErr != nil {
		log.Warnf("[%s] failed path parameter validation: %s", context.ID, validationErr.ErrorMessage)
		return &api.APIResponse{
			StatusCode: http.StatusBadRequest,
			Body:       validationErr,
		}
	}

	trackStore := store.NewMongoStore[models.TrackInfo](injector.MongoInstance, "gostream", "tracks")
	albumStore := store.NewMongoStore[models.AlbumInfo](injector.MongoInstance, "gostream", "albums")

	albumInfo, err := FindAlbumByID(albumStore, id)
	if err != nil {
		log.Warnf("[%s] could not find album: %s", context.ID, err)
		return &api.APIResponse{
			StatusCode: http.StatusNotFound,
		}
	}

	requestBody, err := ExtractRequestBody(request)
	if err != nil {
		log.Warnf("[%s] failed to extract request body: %s", context.ID, err)
		return &api.APIResponse{
			StatusCode: http.StatusBadRequest,
			Body: UpdateAlbumErrorResponseBody{
				Message: "invalid request body",
			},
		}
	}

	validationError := ValidateRequestBody(requestBody)
	if validationError != nil {
		log.Warnf("[%s] failed request body validation: %s", context.ID, validationError.ErrorMessage)
		return &api.APIResponse{
			StatusCode: http.StatusBadRequest,
			Body:       validationError,
		}
	}

	if len(requestBody.TrackIDs) > 0 {
		for _, trackIDs := range requestBody.TrackIDs {
			err = CheckIfTrackExists(trackStore, trackIDs)
			if err != nil {
				log.Warnf("[%s] track does not exist: %s", context.ID, err)
				return &api.APIResponse{
					StatusCode: http.StatusBadRequest,
					Body: UpdateAlbumErrorResponseBody{
						Message: "track does not exist",
					},
				}
			}
		}
	}

	if requestBody.Title != "" {
		albumInfo.Title = requestBody.Title
	}

	if len(requestBody.TrackIDs) > 0 {
		albumInfo.TrackIDs = requestBody.TrackIDs
	}

	if requestBody.Stats.Popularity != 0 {
		albumInfo.Stats.Popularity = requestBody.Stats.Popularity
	}

	updateFilter := query.Filter{
		Root: query.FilterOperatorEq{
			Key:   "_id",
			Value: id,
		},
	}

	updateOperator := query.Update{
		Root: query.UpdateOperatorSet{
			Set: map[string]interface{}{
				"title":            albumInfo.Title,
				"trackIds":         albumInfo.TrackIDs,
				"stats.popularity": albumInfo.Stats.Popularity,
			},
		},
	}

	log.Tracef("[%s] attempting to update database item ...", context.ID)
	count, err := albumStore.UpdateItem(&updateFilter, &updateOperator)

	if err != nil {
		log.Errorf("[%s] failed to update database item: %s", context.ID, err)
		return &api.APIResponse{
			StatusCode: http.StatusInternalServerError,
		}
	}

	if count == 0 {
		log.Warnf("[%s] zero modified items", context.ID)
		return &api.APIResponse{
			StatusCode: http.StatusNoContent,
		}
	}

	log.Tracef("[%s] successfully completed request", context.ID)
	return &api.APIResponse{
		StatusCode: http.StatusNoContent,
	}
}
