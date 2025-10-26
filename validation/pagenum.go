package validation

import (
	"fmt"
	"strconv"
	"strings"
)

// ValidatePageNumber validates that the pageNumber query is good. This function
// will strip all spaces from pageNumber, and return the string as an integer if
// it is valid. Otherwise, it returns an -1 and an error describing how it was incorrect.
func ValidatePageNumber(pageNumber string) (int64, error) {
	pageNumberQuery := strings.TrimSpace(pageNumber)
	if pageNumberQuery == "" {
		return -1, fmt.Errorf("page number cannot be empty")
	}

	pageNumberInt, err := strconv.ParseInt(pageNumberQuery, 10, 64)
	if err != nil {
		return -1, fmt.Errorf("page number must be an integer")
	} else if pageNumberInt < 1 {
		return -1, fmt.Errorf("page number cannot be bellow 1")
	}

	return pageNumberInt, nil
}
