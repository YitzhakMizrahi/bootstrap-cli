// Package utils provides utility functions for common operations across the bootstrap-cli.
// This includes type conversion helpers and other shared functionality.
package utils

// ConvertToInterfaceSlice converts a slice of any type to []interface{}
func ConvertToInterfaceSlice[T any](slice []T) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
} 