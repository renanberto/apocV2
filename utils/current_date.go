package utils

import (
	"time"
)

//GetCurrentDate returns a formated date for the moment it was executed
func GetCurrentDate() string {
	t := time.Now().UTC().Format(time.RFC3339)
	return t
}
