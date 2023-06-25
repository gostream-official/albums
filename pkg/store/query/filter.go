package query

import (
	"go.mongodb.org/mongo-driver/bson"
)

// Description:
//
//	The root filter object.
//	The top level descriptor for any database query.
type Filter struct {

	// The root filter.
	Root IQuery

	// The query result limit.
	Limit uint32
}

// Description:
//
//	The 'and' filter. Allows a combination of multiple filter conditions.
type FilterOperatorAnd struct {

	// The filter interface implementation.
	IQuery

	// All filter conditions checked by the 'and' operator.
	And []IQuery
}

// Description:
//
//	The 'or' filter. Allows a combination of multiple filter conditions.
type FilterOperatorOr struct {

	// The filter interface implementation.
	IQuery

	// All filter conditions checked by the 'or' operator.
	Or []IQuery
}

// Description:
//
//	The 'equals' filter.
//	Allows to filter documents which match a specific field value.
type FilterOperatorEq struct {

	// The filter interface implementation.
	IQuery

	// The document key to refer to.
	Key string

	// The document value which should be matched.
	Value interface{}
}

// Description:
//
//	The 'not equals' filter.
//	Allows to filter documents which do not match a specific field value.
type FilterOperatorNeq struct {

	// The filter interface implementation.
	IQuery

	// The document key to refer to.
	Key string

	// The document value which should not be matched.
	Value interface{}
}

// Description:
//
//	The 'less than' filter.
//	Allows to filter documents for fields with a value less than the given value.
type FilterOperatorLt struct {

	// The filter interface implementation.
	IQuery

	// The document key to refer to.
	Key string

	// The document value should be less than this value.
	Value interface{}
}

// Description:
//
//	The 'less than equals' filter.
//	Allows to filter documents for fields with a value which is less than or equals the given value.
type FilterOperatorLte struct {

	// The filter interface implementation.
	IQuery

	// The document key to refer to.
	Key string

	// The document value should be less than or equal this value.
	Value interface{}
}

// Description:
//
//	The 'greater than' filter.
//	Allows to filter documents for fields with a value greater than the given value.
type FilterOperatorGt struct {

	// The filter interface implementation.
	IQuery

	// The document key to refer to.
	Key string

	// The document value should be greater than this value.
	Value interface{}
}

// Description:
//
//	The 'greater than equals' filter.
//	Allows to filter documents for fields with a value which is greater than or equals the given value.
type FilterOperatorGte struct {

	// The filter interface implementation.
	IQuery

	// The document key to refer to.
	Key string

	// The document value should be greater than or equal this value.
	Value interface{}
}

// Description:
//
//	Compiles the filter and potential sub filters into a MongoDB BSON document.
//
// Returns:
//
//	A MongoDB bson document representing this filter.
func (filter FilterOperatorAnd) Compile() bson.M {
	andArray := make([]bson.M, 0)

	for _, and := range filter.And {
		andArray = append(andArray, and.Compile())
	}

	return bson.M{"$and": andArray}
}

// Description:
//
//	Compiles the filter and potential sub filters into a MongoDB BSON document.
//
// Returns:
//
//	A MongoDB bson document representing this filter.
func (filter FilterOperatorOr) Compile() bson.M {
	orArray := make([]bson.M, 0)

	for _, or := range filter.Or {
		orArray = append(orArray, or.Compile())
	}

	return bson.M{"$or": orArray}
}

// Description:
//
//	Compiles the filter and potential sub filters into a MongoDB BSON document.
//
// Returns:
//
//	A MongoDB bson document representing this filter.
func (filter FilterOperatorEq) Compile() bson.M {
	return bson.M{filter.Key: filter.Value}
}

// Description:
//
//	Compiles the filter and potential sub filters into a MongoDB BSON document.
//
// Returns:
//
//	A MongoDB bson document representing this filter.
func (filter FilterOperatorNeq) Compile() bson.M {
	return bson.M{filter.Key: bson.M{"$ne": filter.Value}}
}

// Description:
//
//	Compiles the filter and potential sub filters into a MongoDB BSON document.
//
// Returns:
//
//	A MongoDB bson document representing this filter.
func (filter FilterOperatorLt) Compile() bson.M {
	return bson.M{filter.Key: bson.M{"$lt": filter.Value}}
}

// Description:
//
//	Compiles the filter and potential sub filters into a MongoDB BSON document.
//
// Returns:
//
//	A MongoDB bson document representing this filter.
func (filter FilterOperatorLte) Compile() bson.M {
	return bson.M{filter.Key: bson.M{"$lte": filter.Value}}
}

// Description:
//
//	Compiles the filter and potential sub filters into a MongoDB BSON document.
//
// Returns:
//
//	A MongoDB bson document representing this filter.
func (filter FilterOperatorGt) Compile() bson.M {
	return bson.M{filter.Key: bson.M{"$gt": filter.Value}}
}

// Description:
//
//	Compiles the filter and potential sub filters into a MongoDB BSON document.
//
// Returns:
//
//	A MongoDB bson document representing this filter.
func (filter FilterOperatorGte) Compile() bson.M {
	return bson.M{filter.Key: bson.M{"$gte": filter.Value}}
}
