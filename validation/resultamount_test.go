package validation

import "testing"

// func ValidateResultAmount(resultAmount string) (int64, error) {
// 	strippedResultAmount := strings.TrimSpace(resultAmount)
// 	if strippedResultAmount == "" {
// 		return -1, fmt.Errorf("you must provide a result amount")
// 	}
//
// 	resultAmountNumber, err := strconv.ParseInt(strippedResultAmount, 10, 64)
// 	if err != nil {
// 		return -1, fmt.Errorf("result amount must be a number")
// 	} else if resultAmountNumber <= 0 {
// 		return -1, fmt.Errorf("result amount must be greater than 0")
// 	} else if resultAmountNumber >= 50 {
// 		return -1, fmt.Errorf("result amount must be less than 50")
// 	}
//
// 	return resultAmountNumber, nil
// }

func Test_ValidateResultAmount_StringToNumber(t *testing.T) {
	inputs := []string{"1", "20", "50"}
	expected := []int64{1, 20, 50}

	for i := range inputs {
		inp := inputs[i]
		exp := expected[i]

		act, err := ValidateResultAmount(inp)
		if err != nil {
			t.Errorf(`ValidateResultAmount("%s"); unexpected error: %v`, inp, err)
		} else if exp != act {
			t.Errorf(`ValidateResultAmount("%s") = "%d"; expected "%v"`, inp, act, exp)
		}
	}
}

func Test_ValidateResultAmount_StringToNumber_IgnoresWhitespace(t *testing.T) {
	inputs := []string{"  1 ", "      20 ", " 50          "}
	expected := []int64{1, 20, 50}

	for i := range inputs {
		inp := inputs[i]
		exp := expected[i]

		act, err := ValidateResultAmount(inp)
		if err != nil {
			t.Errorf(`ValidateResultAmount("%s"); unexpected error: %v`, inp, err)
		} else if exp != act {
			t.Errorf(`ValidateResultAmount("%s") = "%d"; expected "%v"`, inp, act, exp)
		}
	}
}

func Test_ValidateResultAmount_LessThan1ReturnsError(t *testing.T) {
	input := "0"

	actual, err := ValidateResultAmount(input)

	if err == nil {
		t.Errorf(`ValidatePageNumber("%s") = %d; expected error`, input, actual)
	}
}

func Test_ValidateResultAmount_GreaterThan50ReturnsError(t *testing.T) {
	input := "51"

	actual, err := ValidateResultAmount(input)

	if err == nil {
		t.Errorf(`ValidatePageNumber("%s") = %d; expected error`, input, actual)
	}
}

func Test_ValidateResultAmount_1DoesNotError(t *testing.T) {
	input := "1"
	expected := int64(1)

	actual, err := ValidateResultAmount(input)

	if err != nil {
		t.Errorf(`ValidatePageNumber("%s"); unexpected error: %v`, input, err)
	} else if expected != actual {
		t.Errorf(`ValidatePageNumber("%s") = %d; expected %d`, input, actual, expected)
	}
}

func Test_ValidateResultAmount_50DoesNotError(t *testing.T) {
	input := "50"
	expected := int64(50)

	actual, err := ValidateResultAmount(input)

	if err != nil {
		t.Errorf(`ValidatePageNumber("%s"); unexpected error: %v`, input, err)
	} else if expected != actual {
		t.Errorf(`ValidatePageNumber("%s") = %d; expected %d`, input, actual, expected)
	}
}
