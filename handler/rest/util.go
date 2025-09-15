package rest

import (
	"banking/core"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// Maps domain (core) errors to HTTP ones and
// sends them to the client
//
// If an error can't be mapped,
// it sends 500 status code by default and logs the error
//
// A prefix (prfx) can be passed
// to form a better error stack trace
func sendHttpError(
	w http.ResponseWriter,
	req *http.Request,
	prfx string,
	err error,
) {
	switch {
	case errors.Is(err, core.ErrResourceNotFound):
		http.Error(w, "", http.StatusNotFound)
	case errors.Is(err, core.ErrInvalidAccess):
		http.Error(w, "", http.StatusForbidden)
	case errors.Is(err, core.ErrInvalidCredentials):
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
	default:
		http.Error(w, "", http.StatusInternalServerError)
		log.Printf(
			"ERROR: %s - %s %s %s, status: %d, %s\n",
			req.RemoteAddr,
			req.Method,
			req.RequestURI,
			req.Proto,
			http.StatusInternalServerError,
			fmt.Sprintf("%s: %s", prfx, err.Error()),
		)
	}
}
