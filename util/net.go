package util

import (
	"regexp"
	"strconv"
)

func IsValidPortStr(port string) bool {
	val, err := strconv.Atoi(port)
	if err != nil {
		return false
	}
	return val >= 1 && val <= 65535
}

const emailRx = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

func IsValidEmail(email string) bool {
	return regexp.MustCompile(emailRx).MatchString(email)
}
