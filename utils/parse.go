package utils

import (
	"log"
	"strconv"
)

func ParseStringToInt64(s string) int64 {
	result, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Fatal("Error parsing string to int64 failed")
	}
	return result
}

func ParseStringToInt(s string) int {
	result, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal("Error parsing string to int64 failed")
	}
	return result
}
