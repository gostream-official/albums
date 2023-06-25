package query

import "go.mongodb.org/mongo-driver/bson"

// Description:
//
//	Used to update documents in a store.
type Update struct {

	// The root query.
	Root IQuery
}

// Description:
//
//	Updates a specific field by updating a specific field.
type UpdateOperatorSet struct {

	// The query interface implementation.
	IQuery

	// The key-value mappings to set.
	Set map[string]interface{}
}

// Description:
//
//	Compiles the update operator and potential sub operators into a MongoDB BSON document.
//
// Returns:
//
//	A MongoDB bson document representing this update operator.
func (update UpdateOperatorSet) Compile() bson.M {
	return bson.M{"$set": update.Set}
}
