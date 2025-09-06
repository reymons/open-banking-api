package rest

import (
	"net/http"
    "banking/middleware/auth"
)

type ClientHandler struct{}

func NewClientHandler() *ClientHandler {
	return &ClientHandler{}
}

func (h *ClientHandler) GetAll(w http.ResponseWriter, req *http.Request) {
    user, ok := auth.GetJwtUser(w, req)
    if !ok {
        return
    }
    _ = user
    w.Write([]byte("OK"))
}
