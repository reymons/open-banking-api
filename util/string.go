package util

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
)

func HashString(str string) string {
	hash := sha512.New()
	hash.Write([]byte(str))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func Base64URLString(len int) (string, error) {
	bytes := make([]byte, len)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
