package marshal

import (
	"encoding/json"
)

const (

	// The default marshalling prefix.
	DefaultPrefix = ""

	// The default marshalling indent.
	DefaultIndent = "  "
)

// Description:
//
//	Marshals an object into JSON.
//	Does not indent the JSON string.
//
//	If an error occurs, an empty string is returned instead.
//
// Parameters:
//
//	object The object to marshal.
//
// Returns:
//
//	The object as a JSON string.
func Quick(object interface{}) string {
	bytes, err := json.Marshal(object)

	if err != nil {
		return ""
	}

	return string(bytes)
}

// Description:
//
//	Marshals an object into JSON.
//	Does indent the JSON string.
//
// Parameters:
//
//	object The object to marshal.
//
// Returns:
//
//	The object as a JSON string, or an error if marshalling fails.
func WithIndent(object interface{}) (string, error) {
	bytes, err := json.MarshalIndent(object, DefaultPrefix, DefaultIndent)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// Description:
//
//	Marshals an object into JSON.
//	Does indent the JSON string.
//
//	If an error occurs, an empty string is returned instead.
//
// Parameters:
//
//	object The object to marshal.
//
// Returns:
//
//	The object as a JSON string.
func QuickWithIndent(object interface{}) string {
	bytes, err := json.MarshalIndent(object, DefaultPrefix, DefaultIndent)

	if err != nil {
		return ""
	}

	return string(bytes)
}
