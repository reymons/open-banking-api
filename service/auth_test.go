package service

import (
	"strconv"
	"testing"
)

func TestGenerateCode_GeneratesNDigitCode(t *testing.T) {
	lengths := []int{6, 10, 1, 2, 7}

	for _, length := range lengths {
		code, err := generateCode(length)
		if err != nil {
			t.Errorf("Did not expect an error, got: %s", err.Error())
		}
		if len(code) != length {
			t.Errorf("Expcted the length %d, got %d", length, len(code))
		}
		for _, c := range code {
			if _, err := strconv.Atoi(string(c)); err != nil {
				t.Errorf("Expected a number as a code character, got %d", c)
			}
		}
	}
}
