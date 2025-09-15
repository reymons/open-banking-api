package store

import (
	"banking/core"
	"database/sql"
	"errors"
	"fmt"
	"testing"
)

func TestNewCoreError(t *testing.T) {
	type testCase struct {
		prfx            string
		err             error
		expectedErr     error
		expectedFullErr error
	}

	cases := []testCase{
		testCase{
			prfx:            "",
			err:             nil,
			expectedErr:     nil,
			expectedFullErr: nil,
		},
		testCase{
			prfx:            "prfx",
			err:             sql.ErrNoRows,
			expectedErr:     core.ErrResourceNotFound,
			expectedFullErr: fmt.Errorf("prfx: %w", core.ErrResourceNotFound),
		},
	}

	for _, c := range cases {
		got := newCoreError(c.prfx, c.err)

		if !errors.Is(got, c.expectedErr) {
			t.Errorf("create core error: received error is not of core errors")
		}

		// Since we already did the errors.Is() check,
		// it means if one of the errors is nil, then another one is also nil,
		// thus we don't need to compare their exact message strings
		if got == nil {
			continue
		}

		if got.Error() != c.expectedFullErr.Error() {
			t.Errorf("expected error message = %s, got = %s", c.expectedFullErr.Error(), got.Error())
		}
	}
}
