package model

import "time"

type ResetPasswordReq struct {
	ClientID  int
	Token     string
	ExpiresAt time.Time
}

func (m *ResetPasswordReq) Expired() bool {
	return time.Now().After(m.ExpiresAt)
}
