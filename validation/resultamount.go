package validation

import (
	"fmt"
	"strconv"
	"strings"
)

// ValidateResultAmount validates that the resultAmount query is good. This function
// will strip all spaces from resultAmount, and return the string as an integer if
// it is valid. Otherwise, it returns an -1 and an error describing how it was incorrect.
func ValidateResultAmount(resultAmount string) (int64, error) {
	strippedResultAmount := strings.TrimSpace(resultAmount) 
	if strippedResultAmount == "" {
		return -1, fmt.Errorf("you must provide a result amount")
	}

	resultAmountNumber, err := strconv.ParseInt(strippedResultAmount, 10, 64)
	if err != nil {
		return -1, fmt.Errorf("result amount must be a number")
	} else if resultAmountNumber <= 0 {
		return -1, fmt.Errorf("result amount must be greater than 0")
	} else if resultAmountNumber > 50 {
		return -1, fmt.Errorf("result amount must be less than 50")
	}

	return resultAmountNumber, nil
}
