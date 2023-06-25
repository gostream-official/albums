package arrays

// Description:
//
//	Applies a function on every element of the array, returning the modified array.
//	Does not modify the original array.
//
// Parameters:
//
//	array 	The input array.
//	mapFunc The mapping function.
//
// Returns:
//
//	The modified array.
func Map[TArray, UResult any](array []TArray, mapFunc func(TArray) UResult) []UResult {
	result := make([]UResult, len(array))

	for index := range array {
		result[index] = mapFunc(array[index])
	}

	return result
}
