package models

// Description:
//
//	The data model definition for an album.
//	This is a direct reference to the database data model.
type AlbumInfo struct {

	// The id of the album (primary key).
	ID string `json:"id" bson:"_id"`

	// The title of the album.
	Title string `json:"title" bson:"title"`

	// The tracks contained in the album.
	TrackIDs []string `json:"trackIds" bson:"trackIds"`

	// Some album statistics.
	Stats AlbumStats `json:"stats" bson:"stats"`
}

// Description:
//
//	Some album statistics.
type AlbumStats struct {

	// The popularity factor.
	Popularity float32 `json:"popularity" bson:"popularity,truncate"`
}
