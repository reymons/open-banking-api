package httputils 

import (
    "encoding/json"
    "banking/lib/validation"
)

// Decodes an http body and validates it.
//
// If the decoding process fails OR the body is invalid,
// automatically sends a 400-status response
//
// You must check the boolean flag that's returned to learn
// if the response has been sent (decoding failed)
//
// Example:
//
// body, ok := DecodeBody[MyBody](w, req)
// if !ok {
//     return
// }
// Do something with the body, etc.
func DecodeBody[T validation.Validator](w http.ResponseWriter, req *http.Request) (T, error) {
    var v T
	if err := json.NewDecoder(req.Body).Decode(&v); err != nil {
        return v, err
	}
	if problems := v.Valid(); len(problems) > 0 {
        return v, fmt.Errorf("invalid payload")
	}
	return v, nil
}

// JSON-Encodes a body and automatically writes its contents to the response
// If en error occurs, closes the connection and returns the error to the function user
// The user must not call any actions on the response writer afterwards
func EncodeBody[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

