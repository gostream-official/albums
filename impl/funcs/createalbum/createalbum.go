package createalbum

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
//	The request body for the create track endpoint.
type CreateAlbumRequestBody struct {

	// The title of the album.
	Title string `json:"title"`

	// The tracks contained in the album.
	TrackIDs []string `json:"trackIds"`

	// Some album statistics.
	Stats CreateAlbumStatsRequestBody `json:"stats"`
}

// Description:
//
//	The request body for the track statistics.
type CreateAlbumStatsRequestBody struct {

	// The popularity factor.
	Popularity float32 `json:"popularity"`
}

// Description:
//
//	The error response body for the create track endpoint.
type CreateAlbumErrorResponseBody struct {

	// The error message.
	Message string `json:"message"`
}

// Description:
//
//	Describes a validation error.
type CreateAlbumValidationError struct {

	// The JSON field which is referenced by the error message.
	FieldRef string `json:"ref"`

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
		return nil, fmt.Errorf("createalbum: failed to deduce injector")
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
func ExtractRequestBody(request *api.APIRequest) (*CreateAlbumRequestBody, error) {
	body := &CreateAlbumRequestBody{}

	bytes := []byte(request.Body)
	err := json.Unmarshal(bytes, body)

	if err != nil {
		return nil, err
	}

	return body, nil
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
func ValidateRequestBody(request *CreateAlbumRequestBody) *CreateAlbumValidationError {
	title := strings.TrimSpace(request.Title)
	trackIds := arrays.Map(request.TrackIDs, func(artist string) string {
		return strings.TrimSpace(artist)
	})

	if len(title) == 0 {
		return &CreateAlbumValidationError{
			FieldRef:     "title",
			ErrorMessage: "value must not be empty",
		}
	}

	for _, trackID := range trackIds {
		_, err := uuid.Parse(trackID)
		if err != nil {
			return &CreateAlbumValidationError{
				FieldRef:     "trackIds",
				ErrorMessage: "array value must be valid uuid",
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
func ValidateAlbumStatsRequestBody(request *CreateAlbumStatsRequestBody) *CreateAlbumValidationError {
	return nil
}

// Description:
//
//	Checks whether the given track id exists in the mongo store.
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

	requestBody, err := ExtractRequestBody(request)
	if err != nil {
		log.Warnf("[%s] failed to extract request body: %s", context.ID, err)
		return &api.APIResponse{
			StatusCode: http.StatusBadRequest,
			Body: CreateAlbumErrorResponseBody{
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

	trackStore := store.NewMongoStore[models.TrackInfo](injector.MongoInstance, "gostream", "tracks")
	albumStore := store.NewMongoStore[models.AlbumInfo](injector.MongoInstance, "gostream", "albums")

	for _, trackID := range requestBody.TrackIDs {
		err = CheckIfTrackExists(trackStore, trackID)
		if err != nil {
			log.Warnf("[%s] track does not exist: %s", context.ID, err)
			return &api.APIResponse{
				StatusCode: http.StatusBadRequest,
				Body: CreateAlbumErrorResponseBody{
					Message: "track does not exist",
				},
			}
		}
	}

	albumInfo := models.AlbumInfo{
		ID:       uuid.New().String(),
		Title:    requestBody.Title,
		TrackIDs: requestBody.TrackIDs,
		Stats: models.AlbumStats{
			Popularity: requestBody.Stats.Popularity,
		},
	}

	log.Tracef("[%s] attempting to create database item ...", context.ID)
	err = albumStore.CreateItem(albumInfo)

	if err != nil {
		log.Errorf("[%s] failed to create database item: %s", context.ID, err)
		return &api.APIResponse{
			StatusCode: http.StatusInternalServerError,
		}
	}

	log.Tracef("[%s] successfully completed request", context.ID)
	return &api.APIResponse{
		StatusCode: http.StatusOK,
		Body:       albumInfo,
	}
}
