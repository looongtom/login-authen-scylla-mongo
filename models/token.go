package models

type Token struct {
	Token     string `json:"token"`
	User      string `json:"user"`
	CreatedAt string `json:"created_at"`
	ExpiredAt string `json:"expired_at"`
}
