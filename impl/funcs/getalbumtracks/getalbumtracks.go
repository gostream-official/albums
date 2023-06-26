package getalbumtracks

import (
	"fmt"
	"net/http"

	"github.com/gostream-official/albums/impl/inject"
	"github.com/gostream-official/albums/impl/models"
	"github.com/gostream-official/albums/pkg/api"
	"github.com/gostream-official/albums/pkg/marshal"
	"github.com/gostream-official/albums/pkg/parallel"
	"github.com/gostream-official/albums/pkg/store"
	"github.com/gostream-official/albums/pkg/store/query"
	"github.com/revx-official/output/log"
)

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
		return nil, fmt.Errorf("getalbumtracks: failed to deduce injector")
	}

	return &injector, nil
}

// Description:
//
//	Finds all tracks contained in an album.
//
// Parameters:
//
//	store The track store.
//	album The album to search.
//
// Returns:
//
//	All tracks contained in the given album.
//	An error if the database query fails.
func FindTracksForAlbum(store *store.MongoStore[models.TrackInfo], album *models.AlbumInfo) ([]models.TrackInfo, error) {
	filter := query.Filter{}

	filters := make([]query.IQuery, 0)
	for _, trackID := range album.TrackIDs {
		eq := query.FilterOperatorEq{
			Key:   "_id",
			Value: trackID,
		}

		filters = append(filters, eq)
	}

	or := query.FilterOperatorOr{
		Or: filters,
	}

	filter.Root = or

	tracks, err := store.FindItems(&filter)
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

// Description:
//
//	The router handler for getting all tracks in an album.
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
		log.Errorf("[%s] failed to get endpoint injector: %s", context.ID, err)
		return &api.APIResponse{
			StatusCode: http.StatusInternalServerError,
		}
	}

	albumStore := store.NewMongoStore[models.AlbumInfo](injector.MongoInstance, "gostream", "albums")
	trackStore := store.NewMongoStore[models.TrackInfo](injector.MongoInstance, "gostream", "tracks")

	filter := query.Filter{
		Root: query.FilterOperatorEq{
			Key:   "_id",
			Value: request.PathParameters["id"],
		},
		Limit: 10,
	}

	items, err := albumStore.FindItems(&filter)

	if err != nil {
		log.Errorf("[%s] failed to retrieve database items: %s", context.ID, err)
		return &api.APIResponse{
			StatusCode: http.StatusInternalServerError,
		}
	}

	if len(items) == 0 {
		return &api.APIResponse{
			StatusCode: http.StatusNotFound,
		}
	}

	resultItem := items[0]
	if len(resultItem.TrackIDs) == 0 {
		log.Warnf("[%s] album does not contain tracks", context.ID)
		return &api.APIResponse{
			StatusCode: http.StatusOK,
			Body:       []models.TrackInfo{},
		}
	}

	tracks, err := FindTracksForAlbum(trackStore, &resultItem)
	if err != nil {
		log.Errorf("[%s] failed to find album tracks: %s", context.ID, err)
		return &api.APIResponse{
			StatusCode: http.StatusInternalServerError,
		}
	}

	return &api.APIResponse{
		StatusCode: http.StatusOK,
		Body:       tracks,
	}
}
