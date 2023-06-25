package store

import (
	"context"

	"github.com/gostream-official/albums/pkg/store/query"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Description:
//
//	A MongoDB instance.
//	A wrapper around the MongoDB client.
type MongoInstance struct {

	// The MongoDB client.
	Client *mongo.Client
}

// Description:
//
//	A MongoDB store.
//	A wrapper around the MongoDB collection.
type MongoStore[T interface{}] struct {

	// The MongoDB collection.
	Collection *mongo.Collection
}

// Description:
//
//	Creates a new mongo instance.
//	Connects instantly to the given URI.
//
// Parameters:
//
//	uri The MongoDB connection URI.
//
// Returns:
//
//	The created mongo instance, or an error, if the connection fails.
func NewMongoInstance(uri string) (*MongoInstance, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	ctx := context.Background()
	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &MongoInstance{
		Client: client,
	}, nil
}

// Description:
//
//	Creates a new mongo store.
//
// Parameters:
//
//	instance 	The mongo instance which is referred to.
//	database 	The mongo database name referring to.
//	collection 	The mongo collection name referring to.
//
// Type Parameters:
//
//	T The type of document stored in the mongo store to create.
//
// Returns:
//
//	The created mongo store.
func NewMongoStore[T interface{}](instance *MongoInstance, database string, collection string) *MongoStore[T] {
	databaseRef := instance.Client.Database(database)
	collectionRef := databaseRef.Collection(collection)

	return &MongoStore[T]{
		Collection: collectionRef,
	}
}

// Description:
//
//	Creates a new item.
//
// Parameters:
//
//	item The item to create.
//
// Returns:
//
//	An error if creation fails.
func (store *MongoStore[T]) CreateItem(item interface{}) error {
	ctx := context.Background()
	_, err := store.Collection.InsertOne(ctx, item)

	if err != nil {
		return err
	}

	return nil
}

// Description:
//
//	Updates a single item.
//
// Parameters:
//
//	filter The filter used for searching the documents to update.
//	update The update operator used for updating the filtered documents.
//
// Returns:
//
//	The number of modified documents.
//	An error if the update fails.
func (store *MongoStore[T]) UpdateItem(filter *query.Filter, update *query.Update) (int64, error) {
	var query bson.M
	var updateQuery bson.M

	if filter.Root == nil {
		query = bson.M{}
	} else {
		query = filter.Root.Compile()
	}

	if update.Root == nil {
		updateQuery = bson.M{}
	} else {
		updateQuery = update.Root.Compile()
	}

	ctx := context.Background()
	result, err := store.Collection.UpdateOne(ctx, query, updateQuery)

	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}

// Description:
//
//	Queries items in the store.
//
// Parameters:
//
//	The query filter to use.
//
// Returns:
//
//	An array of all items matching the given query filter.
//	An error if the query fails.
func (store *MongoStore[T]) FindItems(filter *query.Filter) ([]T, error) {
	items := make([]T, 0)

	var query bson.M

	if filter.Root == nil {
		query = bson.M{}
	} else {
		query = filter.Root.Compile()
	}

	ctx := context.Background()
	options := options.Find().SetLimit(int64(filter.Limit))

	cursor, err := store.Collection.Find(ctx, query, options)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var item T
		err := cursor.Decode(&item)

		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

// Description:
//
//	Deletes an item by its ID.
//
// Parameters:
//
//	The ID of the document to delete.
//
// Returns:
//
//	The number of deleted documents.
//	An error if the request fails.
func (store MongoStore[T]) DeleteItem(id string) (int64, error) {
	ctx := context.Background()

	result, err := store.Collection.DeleteOne(ctx, bson.M{
		"_id": id,
	})

	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
