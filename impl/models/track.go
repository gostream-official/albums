package models

import "time"

// Description:
//
//	The data model definition for a track.
//	This is a direct reference to the database data model.
type TrackInfo struct {

	// The id of the track (primary key).
	ID string `json:"id" bson:"_id"`

	// The id of the track artist.
	ArtistID string `json:"artistId" bson:"artistId"`

	// Additional ids of featuring artists.
	FeaturedArtistIDs []string `json:"featuredArtistIds" bson:"featuredArtistIds"`

	// The title of the track.
	Title string `json:"title" bson:"title"`

	// The label that published the track.
	Label string `json:"label" bson:"label"`

	// The release date of the track.
	ReleaseDate time.Time `json:"releaseDate" bson:"releaseDate"`

	// Some track statistics.
	TrackStats TrackStats `json:"trackStats" bson:"trackStats"`

	// Some audio features of the track.
	AudioFeatures AudioFeatures `json:"audioFeatures" bson:"audioFeatures"`
}

// Description:
//
//	Represents some track statistics.
type TrackStats struct {

	// The amount of streams of the track.
	Streams uint32 `json:"streams" bson:"streams"`

	// The amount of likes of the track.
	Likes uint32 `json:"likes" bson:"likes"`
}

// Descriptions:
//
//	Represents some audio features of a track.
type AudioFeatures struct {

	// The key of the track.
	Key string `json:"key" bson:"key"`

	// The tempo of the track.
	Tempo float32 `json:"tempo" bson:"tempo"`

	// The duration of the track.
	Duration float32 `json:"duration" bson:"duration"`

	// The energy level of the track.
	Energy float32 `json:"energy" bson:"energy"`

	// The danceability level of the track.
	Danceability float32 `json:"danceability" bson:"danceability"`

	// The accousticness level of the track.
	Accousticness float32 `json:"accousticness" bson:"accousticness"`

	// The instrumentalness level of the track.
	Instrumentalness float32 `json:"instrumentalness" bson:"instrumentalness"`

	// The liveness level of the track.
	Liveness float32 `json:"liveness" bson:"liveness"`

	// The loudness of the track (in LUFS).
	Loudness float32 `json:"loudness" bson:"loudness"`

	// The time signature of the track.
	TimeSignature int `json:"timeSignature" bson:"timeSignature"`
}

const (

	// The A Minor key.
	AudioKeyAmin = "A Minor"

	// The B Minor key.
	AudioKeyBmin = "B Minor"

	// The C Minor key.
	AudioKeyCmin = "C Minor"

	// The D Minor key.
	AudioKeyDmin = "D Minor"

	// The E Minor key.
	AudioKeyEmin = "E Minor"

	// The F Minor key.
	AudioKeyFmin = "F Minor"

	// The G Minor key.
	AudioKeyGmin = "G Minor"

	// The A Major key.
	AudioKeyAmaj = "A Major"

	// The B Major key.
	AudioKeyBmaj = "B Major"

	// The C Major key.
	AudioKeyCmaj = "C Major"

	// The D Major key.
	AudioKeyDmaj = "D Major"

	// The E Major key.
	AudioKeyEmaj = "E Major"

	// The F Major key.
	AudioKeyFmaj = "F Major"

	// The G Major key.
	AudioKeyGmaj = "G Major"

	// The A# Minor key.
	AudioKeyAsharpmin = "A# Minor"

	// The B# Minor key.
	AudioKeyBsharpmin = "B# Minor"

	// The C# Minor key.
	AudioKeyCsharpmin = "C# Minor"

	// The D# Minor key.
	AudioKeyDsharpmin = "D# Minor"

	// The E# Minor key.
	AudioKeyEsharpmin = "E# Minor"

	// The F# Minor key.
	AudioKeyFsharpmin = "F# Minor"

	// The G# Minor key.
	AudioKeyGsharpmin = "G# Minor"

	// The Ab Major key.
	AudioKeyAflatmaj = "Ab Major"

	// The Bb Major key.
	AudioKeyBflatmaj = "Bb Major"

	// The Cb Major key.
	AudioKeyCflatmaj = "Cb Major"

	// The Db Major key.
	AudioKeyDflatmaj = "Db Major"

	// The Eb Major key.
	AudioKeyEflatmaj = "Eb Major"

	// The Fb Major key.
	AudioKeyFflatmaj = "Fb Major"

	// The Gb Major key.
	AudioKeyGflatmaj = "Gb Major"
)
