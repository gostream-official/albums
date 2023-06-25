package router

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gostream-official/albums/pkg/api"
)

// Description:
//
//	Implementation of the Router interface for gin.
type GinRouter struct {

	// The gin engine.
	engine *gin.Engine
}

// Description:
//
//	Package initializer.
//	Sets gin to release mode.
func init() {
	gin.SetMode(gin.ReleaseMode)
}

// Description:
//
//	Creates a new gin router.
//
// Returns:
//
//	The created gin router.
func NewGinRouter() *GinRouter {
	engine := gin.New()

	engine.RedirectTrailingSlash = true
	engine.RedirectFixedPath = true

	return &GinRouter{
		engine: engine,
	}
}

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
func (router *GinRouter) Handle(method string, path string, handler RouterHandlerFunc) {
	router.engine.Handle(method, path, func(context *gin.Context) {
		internalRouteHandler(path, context, handler)
	})
}

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
func (router *GinRouter) HandleWith(method string, path string, handler RouterInjectionHandlerFunc) *RouterInjector {
	injector := &RouterInjector{}

	router.engine.Handle(method, path, func(context *gin.Context) {
		internalRouteInjectionHandler(path, context, handler, injector)
	})

	return injector
}

// Description:
//
//	Starts the HTTP server for this router and listens to all registered routes.
//
// Returns:
//
//	An error if serving the router fails.
func (router *GinRouter) Run(port uint16) error {
	portFmt := fmt.Sprintf(":%d", port)

	// server := &http.Server{
	// 	Addr:    portFmt,
	// 	Handler: router.engine,
	// }

	// // Setting this to false apparently reduces memory usage.
	// // However, setting this to true apparently is the standard and improves performance.
	// server.SetKeepAlivesEnabled(true)

	// return server.ListenAndServe()

	return router.engine.Run(portFmt)
}

// Description:
//
//	Internal handler method for incoming requests.
//	Triggered by the gin framework.
//
// Parameters:
//
//	pathHandle 	The registered path handle.
//	context 	The internal gin context.
//	handler 	The registered handler function.
func internalRouteHandler(pathHandle string, context *gin.Context, handler RouterHandlerFunc) {
	request := context.Request

	internalRequest, err := transformRequest(pathHandle, request)

	if err != nil {
		panic("router: cannot transform request")
	}

	internalResponse := handler(internalRequest)
	applyResponse(internalResponse, context)
}

// Description:
//
//	Internal handler method for incoming requests.
//	Allows object injection for the route handler.
//	Triggered by the gin framework.
//
// Parameters:
//
//	pathHandle 	The registered path handle.
//	context 	The internal gin context.
//	handler 	The registered handler function.
//	injector	The router injector.
func internalRouteInjectionHandler(pathHandle string, context *gin.Context, handler RouterInjectionHandlerFunc, injector *RouterInjector) {
	request := context.Request

	internalRequest, err := transformRequest(pathHandle, request)

	if err != nil {
		panic("router: cannot transform request")
	}

	internalResponse := handler(internalRequest, injector.Injector)
	applyResponse(internalResponse, context)
}

// Description:
//
//	Transforms an incoming HTTP request to a router request.
//
// Parameters:
//
//	pathHandle 	The registered path handle.
//	request		The request to transform.
//
// Returns:
//
//	The transformed request, or an error, if the request could not be transformed.
func transformRequest(pathHandle string, request *http.Request) (*api.APIRequest, error) {
	result := api.APIRequest{
		Url:             request.URL.String(),
		Path:            request.URL.Path,
		Method:          request.Method,
		Headers:         make(map[string]string),
		PathParameters:  make(map[string]string),
		QueryParameters: make(map[string]string),
	}

	for key, values := range request.Header {
		result.Headers[key] = strings.Join(values, ",")
	}

	pathParameters, err := extractPathParameters(pathHandle, request.URL.Path)
	if err != nil {
		return nil, err
	}

	result.PathParameters = pathParameters

	queryParameters, err := extractQueryParameters(request.URL.String())
	if err != nil {
		return nil, err
	}

	result.QueryParameters = queryParameters

	defer request.Body.Close()

	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	result.Body = string(body)
	return &result, nil
}

// Description:
//
//	Extracts all path parameters using the registered path handle and the actual request path.
//
// Example:
//   - handle: 	/some/path/:variable
//   - path:	/some/path/128
//
// Parameters:
//
//	handle The registered path handle.
//	path The actual request path.
//
// Returns:
//
//	A key-value map of the extracted path parameters.
func extractPathParameters(handle string, path string) (map[string]string, error) {
	result := make(map[string]string)

	path = strings.TrimPrefix(path, "/")
	handle = strings.TrimPrefix(handle, "/")

	pathSegments := strings.Split(path, "/")
	handleSegments := strings.Split(handle, "/")

	if len(handleSegments) != len(pathSegments) {
		return nil, fmt.Errorf("router: number of url segments does not match number of path segments")
	}

	for index, segment := range handleSegments {
		if !strings.HasPrefix(segment, ":") {
			continue
		}

		paramName := strings.TrimPrefix(segment, ":")
		result[paramName] = pathSegments[index]
	}

	return result, nil
}

// Description:
//
//	Extracts all query parameters from the request path.
//
// Parameters:
//
//	path The actual request path.
//
// Returns:
//
//	A key-value map of the extracted query parameters.
func extractQueryParameters(path string) (map[string]string, error) {
	parameters := make(map[string]string)

	parsedURL, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	query := parsedURL.Query()

	for key, values := range query {
		parameters[key] = strings.Join(values, ",")
	}

	return parameters, nil
}

// Description:
//
//	Applies a router response to the internal gin context.
//
// Parameters:
//
//	response 	The response to apply.
//	context 	The gin context.
func applyResponse(response *api.APIResponse, context *gin.Context) {
	for key, value := range response.Headers {
		context.Header(key, value)
	}

	if response.Body == nil {
		context.Status(response.StatusCode)
		return
	}

	context.JSON(response.StatusCode, response.Body)
}
