package api

// Description:
//
//	A representation of a HTTP response.
type APIResponse struct {

	// The response status code.
	StatusCode int `json:"statusCode"`

	// The response headers.
	Headers map[string]string `json:"headers"`

	// The response body, represented as an object.
	Body interface{} `json:"body"`
}
