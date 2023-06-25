package api

// Description:
//
//	The representation of a HTTP request.
type APIRequest struct {

	// The requested URL.
	Url string

	// The requested path.
	Path string

	// The request method used.

	Method string

	// The request headers
	Headers map[string]string `json:"headers"`

	// A key-value mapping of path parameters.
	PathParameters map[string]string `json:"pathParameters"`

	// A key-value mapping of query parameters.
	QueryParameters map[string]string `json:"queryParameters"`

	// The request body.
	Body string `json:"body"`
}
