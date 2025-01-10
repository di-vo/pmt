package lib

func Contains[T comparable](slice []T, element T) bool {
	for _, val := range slice {
		if val == element {
			return true
		}
	}
	return false
}
