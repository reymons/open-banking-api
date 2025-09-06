package validation

import (
    "fmt"
)

type Validator interface {
    // Validates an object and returns a problems' map
    // If len(map) > 0, the object is invalid
    Valid() (problems map[string]string)
}

type err struct {}

var Err err

func (err) Min(n int) string {
    return fmt.Sprintf("shouldn't be less than %d", n)
}
func (err) Max(n int) string {
    return fmt.Sprintf("shouldn't be greater than %d", n)
}
func (err) Required() string {
    return "required field"
}
func (err) InvalidEmail() string {
    return "invalid email"
}

type valid struct{}

var Valid valid

func (valid) Email(v string) bool {
    _ = v
    return true
}


