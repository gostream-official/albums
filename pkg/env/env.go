package env

import (
	"fmt"
	"os"
)

// Description:
//
//	Gets the value of an environment variable, if it exists.
//
// Parameters:
//
//	name The name of the environment variable to get.
//
// Returns:
//
//	The value of the environment variable, or an error.
func GetEnvironmentVariable(name string) (string, error) {
	variable, exists := os.LookupEnv(name)

	if !exists {
		return "", fmt.Errorf("%s: %s", "env: cannot find environment variable", name)
	}

	return variable, nil
}

// Description:
//
//	Gets the value of an environment variable, if it exists.
//	If it does not exist, the given fallback value is returned.
//
// Parameters:
//
//	name 		The name of the environment variable to get.
//	fallback 	The fallback value to use if the environment variable does not exist.
//
// Returns:
//
//	The value of the environment variable, if it exists, the fallback value otherwise.
func GetEnvironmentVariableWithFallback(name string, fallback string) string {
	variable, exists := os.LookupEnv(name)

	if !exists {
		return fallback
	}

	return variable
}
