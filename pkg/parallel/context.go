package parallel

import (
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"time"
)

// Description:
//
//	A parallel context.
//	Parallel contexts are used to identify the current thread.
//	This can be useful for logging in endpoints. If an endpoint is called multiple times,
//	asynchronous logging can distort the order of the logs, so that the correct logs might not be identified easily.
//
//	Using a parallel context, every log can be assigned an identifier, so that the behavior can be comprehended.
type Context struct {

	// The id of the context. Equals the short id.
	ID string

	// The short id of the context. Unique in most of the cases.
	ShortID string

	// The long id of the context. Unique in pretty much every case.
	LongID string
}

// Description:
//
//	Creates a new parallel context.
//
// Returns:
//
//	The created context.
func NewContext() *Context {
	sha := sha256.New()

	currentTime := time.Now().UnixNano()

	input := strconv.FormatInt(currentTime, 10)
	bytes := []byte(input)

	sha.Write(bytes)

	result := sha.Sum(nil)
	decoded := base64.URLEncoding.EncodeToString(result)
	shortend := decoded[0:7]

	return &Context{
		ID:      shortend,
		ShortID: shortend,
		LongID:  decoded,
	}
}
