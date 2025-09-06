package util_test

import (
	"testing"
    "banking/util"
)

func TestVerifyPassword(t *testing.T) {
    password := "123qwe"
    hashedPassword, err := util.HashPassword(password)
    if err != nil {
        t.Fatal(err)
    }
    if !util.VerifyPassword(password, hashedPassword) {
        t.Error("VerifyPassword() = false; want true")
    }
}
