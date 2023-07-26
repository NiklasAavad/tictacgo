package gamesocket

import "errors"

// Removes the element at index i from slice s
// Will return the slice and should be assigned to the original slice, as for append
func RemoveAt[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// Removes a specific element from a slice. Will return the original slice if the element was not found, but will return an errar
// Will return the slice and should be assigned to the original slice, as for append
func RemoveElement[T comparable](s []T, e T) ([]T, error) {
	for i, v := range s {
		if v == e {
			return RemoveAt(s, i), nil
		}
	}
	return s, errors.New("Element not found")
}
