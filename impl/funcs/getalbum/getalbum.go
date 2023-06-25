package getalbum

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
		return nil, fmt.Errorf("getalbum: failed to deduce injector")
	}

	return &injector, nil
}

// Description:
//
//	The router handler for getting an album.
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

	store := store.NewMongoStore[models.TrackInfo](injector.MongoInstance, "gostream", "albums")

	filter := query.Filter{
		Root: query.FilterOperatorEq{
			Key:   "_id",
			Value: request.PathParameters["id"],
		},
		Limit: 10,
	}

	items, err := store.FindItems(&filter)

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
	return &api.APIResponse{
		StatusCode: http.StatusOK,
		Body:       resultItem,
	}
}
