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
		sendHttpError(w, req, "sign in", err)
		return
	}

	accessToken, err := util.CreateJwtToken(util.AccessTokenDuration, cli.ID, cli.Role)
	if err != nil {
		sendHttpError(w, req, "create jwt token", err)
		return
	}

	util.SetJwtTokenCookie(w, util.AccessTokenCookie, accessToken, util.AccessTokenDuration)

	res := signInRes{
		ID:        cli.ID,
		FirstName: cli.FirstName,
		LastName:  cli.LastName,
		Email:     cli.Email,
	}
	util.EncodeBody(w, http.StatusOK, &res)
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, req *http.Request) {
	body, ok := util.DecodeBody[signUpReq](w, req)
	if !ok {
		return
	}

	err := h.authService.SignUp(
		req.Context(),
		body.FirstName,
		body.LastName,
		body.BirthDate,
		body.Email,
		body.Password,
	)
	if err != nil {
		sendHttpError(w, req, "sign up", err)
		return
	}
}

func (h *AuthHandler) SendVerificationCode(w http.ResponseWriter, req *http.Request) {
	body, ok := util.DecodeBody[sendVerificationCodeReq](w, req)
	if !ok {
		return
	}

	err := h.authService.SendVerificationCode(req.Context(), body.Email)
	if err != nil {
		sendHttpError(w, req, "send verification code", err)
	}
}

func (h *AuthHandler) SubmitVerification(w http.ResponseWriter, req *http.Request) {
	body, ok := util.DecodeBody[submitVerificationReq](w, req)
	if !ok {
		return
	}

	cli, err := h.authService.SubmitVerification(
		req.Context(),
		body.Email,
		body.Code,
	)
	if err != nil {
		sendHttpError(w, req, "submit verification", err)
		return
	}

	accessToken, err := util.CreateJwtToken(util.AccessTokenDuration, cli.ID, cli.Role)
	if err != nil {
		sendHttpError(w, req, "create jwt token", err)
		return
	}

	util.SetJwtTokenCookie(w, util.AccessTokenCookie, accessToken, util.AccessTokenDuration)

	res := submitVerificationRes{
		ID:        cli.ID,
		FirstName: cli.FirstName,
		LastName:  cli.LastName,
		Email:     cli.Email,
		BirthDate: cli.BirthDate,
	}
	util.EncodeBody(w, http.StatusOK, &res)
}
