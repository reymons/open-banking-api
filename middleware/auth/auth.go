package auth

import (
	"banking/util"
	"context"
	"errors"
	"net/http"
)

const ctxKey = "claims"

func Middleware(
	handler func(http.ResponseWriter, *http.Request),
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		tokenCookie, err := req.Cookie(util.AccessTokenCookie)
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		claims, err := util.VerifyJwtToken(tokenCookie.Value)
		if err != nil {
			http.Error(w, "verify token: "+err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(
			req.Context(),
			ctxKey,
			util.JwtUser{
				ID:   claims.ID,
				Role: claims.Role,
			},
		)
		req = req.WithContext(ctx)
		handler(w, req)
	}
}

func GetJwtUserFromReq(req *http.Request) (util.JwtUser, error) {
	if u, ok := req.Context().Value(ctxKey).(util.JwtUser); !ok {
		return util.JwtUser{}, errors.New("claims missing")
	} else {
		return u, nil
	}
}

func GetJwtUser(w http.ResponseWriter, req *http.Request) (util.JwtUser, bool) {
	u, err := GetJwtUserFromReq(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
	return u, err == nil
}
