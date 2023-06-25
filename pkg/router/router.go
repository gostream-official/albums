package router

import "github.com/gostream-official/albums/pkg/api"

// Description:
//
//	Function definition for router endpoint handlers.
type RouterHandlerFunc = func(request *api.APIRequest) *api.APIResponse

// Description:
//
//	Function definition for router endpoint handlers, which support object injection.
type RouterInjectionHandlerFunc = func(request *api.APIRequest, injector interface{}) *api.APIResponse

// Description:
//
//	The router interface.
type Router interface {

	// Description:
	//
	//	Registers a new HTTP handler function for the given method and path.
	//	Paths can include wildcards and path variables.
	//
	// Parameters:
	//
	//	method 	The http method to handle.
	//	path   	The path to handle.
	//	handler	The handler responsible for handling the request.
	Handle(method string, path string, handler RouterHandlerFunc)

	// Description:
	//
	//	Registers a new HTTP handler function for the given method and path.
	//	Paths can include wildcards and path variables.
	//
	//	This method allows object injection for the router handler.
	//
	// Parameters:
	//
	//	method 	The http method to handle.
	//	path   	The path to handle.
	//	handler	The handler responsible for handling the request.
	//
	// Returns:
	//
	//	The router injector which allows object injection for the registered endpoint.
	HandleWith(method string, path string, handler RouterInjectionHandlerFunc) *RouterInjector

	// Description:
	//
	//	Starts the HTTP server for this router and listens to all registered routes.
	//
	// Returns:
	//
	//	An error if serving the router fails.
	Run(port uint16) error
}

// Description:
//
//	The router injector.
//	Responsible for object injection for route handlers.
type RouterInjector struct {

	// The object to inject.
	Injector interface{}
}

// Description:
//
//	Creates the default router.
//
// Returns:
//
//	The default router.
func Default() Router {
	return NewGinRouter()
}

// Description:
//
//	Injects the given object to the endpoint this method is called on.
//
// Parameters:
//
//	object The object to inject.
func (handler *RouterInjector) Inject(object interface{}) {
	handler.Injector = object
}
