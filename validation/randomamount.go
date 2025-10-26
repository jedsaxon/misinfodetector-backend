package validation

import (
	"fmt"
)

// ValidateRandomAmount validates that the randomAmount value is good. This function
// will return an error if the amount is greater than 20,000, or less than 1.
func ValidateRandomAmount(randomAmount int) error {
	if randomAmount > 20000 {
		return fmt.Errorf("random amount cannot be greater than 20,000")
	} else if randomAmount < 1 {
		return fmt.Errorf("random amount cannot be less than 1")
	}

	return nil
}
