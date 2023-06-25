package query

import "go.mongodb.org/mongo-driver/bson"

// Description:
//
//	The filter interface.
type IQuery interface {

	// Description:
	//
	//	Compiles the filter into a bson document for MongoDB.
	//
	// Returns:
	//
	//	The filter represented as a MongoDB bson document.
	Compile() bson.M
}
