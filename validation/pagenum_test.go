package validation

import "testing"

func Test_ValidatePageNum_StringToNumber(t *testing.T) {
	inputs := []string{"1", "100", "500"}
	expected := []int64{1, 100, 500}

	for i := range inputs {
		inp := inputs[i]
		exp := expected[i]

		act, err := ValidatePageNumber(inp)
		if err != nil {
			t.Errorf(`ValidatePageNumber("%s"); unexpected error: %v`, inp, err)
		} else if exp != act {
			t.Errorf(`ValidatePageNumber("%s") = "%d"; expected "%v"`, inp, act, exp)
		}
	}
}

func Test_ValidatePageNum_StringToNumber_IgnoresWhitespace(t *testing.T) {
	inputs := []string{"    1 ", "  100  ", " 500 "}
	expected := []int64{1, 100, 500}

	for i := range inputs {
		inp := inputs[i]
		exp := expected[i]

		act, err := ValidatePageNumber(inp)
		if err != nil {
			t.Errorf(`ValidatePageNumber("%s"); unexpected error: %v`, inp, err)
		} else if exp != act {
			t.Errorf(`ValidatePageNumber("%s") = "%d"; expected "%v"`, inp, act, exp)
		}
	}
}

func Test_ValidatePageNum_Page0_ReturnsError(t *testing.T) {
	input := "0"

	actual, err := ValidatePageNumber(input)

	if err == nil {
		t.Errorf(`ValidatePageNumber("%s") = "%d"; expected error`, input, actual)
	}
}

func Test_ValidatePgaeNum_Page1_DoesNotError(t *testing.T) {
	input := "1"

	_, err := ValidatePageNumber(input)

	if err != nil {
		t.Errorf(`ValidatePageNumber("%s"); unexpected error: %v`, input, err)
	}
}
