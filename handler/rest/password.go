package rest

import (
	"banking/service"
	"banking/util"
	"net/http"
)

type PasswordHandler struct {
	passwordService service.PasswordService
}

func NewPasswordHandler(ps service.PasswordService) *PasswordHandler {
	return &PasswordHandler{ps}
}

func (h *PasswordHandler) RequestPasswordReset(w http.ResponseWriter, req *http.Request) {
	body, ok := util.DecodeBody[resetPswdReq](w, req)
	if !ok {
		return
	}

	err := h.passwordService.RequestPasswordReset(req.Context(), body.Email)
	if err != nil {
		logError(err, "reset password", req)
	}

	// Do not send an error to the client
	// to prevent a user from knowning if an email exists
	w.WriteHeader(http.StatusCreated)
}

func (h *PasswordHandler) ResetPassword(w http.ResponseWriter, req *http.Request) {
	body, ok := util.DecodeBody[submitResetTokenReq](w, req)
	if !ok {
		return
	}

	err := h.passwordService.ResetPassword(req.Context(), body.Token, body.Password)
	if err != nil {
		sendHttpError(w, req, "submit token", err)
	}
}
