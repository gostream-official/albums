package inject

import "github.com/gostream-official/albums/pkg/store"

// Description:
//
//	The injector object for this service.
//	This object is used for endpoint dependency injection.
type Injector struct {

	// The MongoDB store instance.
	MongoInstance *store.MongoInstance
}
