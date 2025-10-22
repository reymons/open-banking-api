package rest

import (
	"banking/middleware/auth"
	"banking/store"
	"banking/util"
	"net/http"
)

type UserHandler struct {
	clientStore store.Client
}

func NewUserHandler(cs store.Client) *UserHandler {
	return &UserHandler{cs}
}

func (h *UserHandler) GetMe(w http.ResponseWriter, req *http.Request) {
	jwtUser, ok := auth.GetJwtUser(w, req)
	if !ok {
		return
	}

	user, err := h.clientStore.GetByID(req.Context(), jwtUser.ID)
	if err != nil {
		sendHttpError(w, req, "get user by id", err)
		return
	}

	res := getUserMeRes{
		ID:        user.ID,
		Email:     user.Email,
		BirthDate: user.BirthDate,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	util.EncodeBody(w, http.StatusOK, &res)
}
