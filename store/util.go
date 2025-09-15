package store

import (
	"banking/core"
	"database/sql"
	"fmt"
)

// Maps SQL errors to domain ones
//
// If an error can't be mapped,
// returns a custom error
//
// A prefix (prfx) can be passed
// to form a better error stack trace
func newCoreError(prfx string, err error) error {
	if err == nil {
		return nil
	}
	if err == sql.ErrNoRows {
		return fmt.Errorf("%s: %w", prfx, core.ErrResourceNotFound)
	}
	return fmt.Errorf("%s: %w", prfx, err)
}
