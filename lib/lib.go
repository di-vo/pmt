package lib

import (
	"fmt"
	"os"
)

func Contains[T comparable](slice []T, element T) bool {
	for _, val := range slice {
		if val == element {
			return true
		}
	}
	return false
}

func WriteToLog(s string) {
	file, err := os.Create("debug.log")
	if err != nil {
		fmt.Println("fatal: ", err)
		os.Exit(1)
	}
	defer file.Close()

	_, err = file.WriteString(s)
	if err != nil {
		fmt.Println("fatal: ", err)
		os.Exit(1)
	}
}
