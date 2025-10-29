package config

import (
	"fmt"
	"log"
)

func fatalEmpty(prfx string, envVar string) {
	log.Fatal(fmt.Sprintf("%s is empty. Specify %s env variable", prfx, envVar))
}

func fatalInvalidPort(prfx string, envVar string) {
	log.Fatal(fmt.Sprintf("%s is invalid. Specify %s env variable with a value between 1 and 65535", prfx, envVar))
}

func fatalInvalidEmail(prfx string, envVar string) {
	log.Fatal(fmt.Sprintf("%s email is invalid. Specify %s env variable with a valid email address", prfx, envVar))
}
