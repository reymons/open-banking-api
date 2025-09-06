package rest

import (
	"banking/service"
	"banking/util"
	"net/http"
)

type AuthHandler struct {
	authService service.Auth
}

func NewAuthHandler(as service.Auth) *AuthHandler {
	return &AuthHandler{as}
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, req *http.Request) {
	body, ok := util.DecodeBody[signInReq](w, req)
	if !ok {
		return
	}

	cli, err := h.authService.SignIn(
		req.Context(),
		body.Email,
		body.Password,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	accessToken, err := util.CreateJwtToken(util.AccessTokenDuration, cli.ID, cli.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.SetJwtTokenCookie(w, util.AccessTokenCookie, accessToken, util.AccessTokenDuration)

	res := signInRes{
		ID:        cli.ID,
		FirstName: cli.FirstName,
		LastName:  cli.LastName,
		Email:     cli.Email,
		Phone:     cli.Phone,
	}
	util.EncodeBody(w, http.StatusOK, &res)
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, req *http.Request) {
	body, ok := util.DecodeBody[signUpReq](w, req)
	if !ok {
		return
	}

	cli, err := h.authService.SignUp(
		req.Context(),
		body.FirstName,
		body.LastName,
		body.BirthDate,
		body.Email,
		body.Phone,
		body.Password,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	accessToken, err := util.CreateJwtToken(util.AccessTokenDuration, cli.ID, cli.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.SetJwtTokenCookie(w, util.AccessTokenCookie, accessToken, util.AccessTokenDuration)

	res := signUpRes{
		ID:        cli.ID,
		FirstName: cli.FirstName,
		LastName:  cli.LastName,
		Email:     cli.Email,
		Phone:     cli.Phone,
	}
	util.EncodeBody(w, http.StatusOK, &res)
}
