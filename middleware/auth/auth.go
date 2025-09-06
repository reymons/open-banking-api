package auth

import (
	"banking/util"
	"context"
	"net/http"
)

const ctxKey = "claims"

func Middleware(
	handler func(http.ResponseWriter, *http.Request),
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		tokenCookie, err := req.Cookie(util.AccessTokenCookie)
		if err != nil {
			http.Error(w, "no access token cookie", http.StatusUnauthorized)
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

func GetJwtUser(w http.ResponseWriter, req *http.Request) (util.JwtUser, bool) {
	user, ok := req.Context().Value(ctxKey).(util.JwtUser)
	if !ok {
		http.Error(w, "claims missing", http.StatusUnauthorized)
	}
	return user, ok
}
