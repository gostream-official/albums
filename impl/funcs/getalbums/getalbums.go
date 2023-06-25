package getalbums

import (
	"fmt"
	"net/http"
	"strconv"

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
		return nil, fmt.Errorf("getalbums: failed to deduce injector")
	}

	return &injector, nil
}

// Description:
//
//	Creates a query filter from the incoming API request.
//
// Parameters:
//
//	request The incoming API request.
//
// Returns:
//
//	The created query filter.
func CreateFilterFromQueryParameters(request *api.APIRequest) query.Filter {
	andFilter := query.FilterOperatorAnd{
		And: make([]query.IQuery, 0),
	}

	var realLimit int
	var realLimitErr error

	limit, limitOk := request.QueryParameters["limit"]
	if limitOk {
		realLimit, realLimitErr = strconv.Atoi(limit)
	}

	resultFilter := query.Filter{}

	if limitOk && realLimitErr == nil {
		resultFilter.Limit = uint32(realLimit)
	}

	if len(andFilter.And) > 0 {
		resultFilter.Root = andFilter
	}

	return resultFilter
}

// Description:
//
//	The router handler for: Get Track By ID
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
	filter := CreateFilterFromQueryParameters(request)

	items, err := store.FindItems(&filter)

	if err != nil {
		log.Errorf("[%s] failed to retrieve database items: %s", context.ID, err)
		return &api.APIResponse{
			StatusCode: http.StatusInternalServerError,
		}
	}

	return &api.APIResponse{
		StatusCode: http.StatusOK,
		Body:       items,
	}
}
