package util

import (
	"banking/core/security"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

const (
	AccessTokenCookie   = "accessToken"
	AccessTokenDuration = 1 * 24 * time.Hour // 1 day
)

var jwtSecret string

func SetJwtSecret(secret string) {
	jwtSecret = secret
}

type UserClaims struct {
	jwt.RegisteredClaims
	ID   int64         `json:"id"`
	Role security.Role `json:"role"`
}

type JwtUser struct {
	ID   int64
	Role security.Role
}

func CreateJwtToken(duration time.Duration, userID int64, role security.Role) (string, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("create jwt token: %w", err)
	}

	claims := UserClaims{
		ID:   userID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Subject:   strconv.FormatInt(userID, 10),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	return token.SignedString([]byte(jwtSecret))
}

func VerifyJwtToken(token string) (*UserClaims, error) {
	tkn, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if userClaims, ok := tkn.Claims.(*UserClaims); !ok {
		return nil, fmt.Errorf("invalid token claims")
	} else {
		return userClaims, nil
	}
}

func SetJwtTokenCookie(
	w http.ResponseWriter,
	name string,
	token string,
	duration time.Duration,
) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    token,
		MaxAge:   int(duration / time.Second),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
}
