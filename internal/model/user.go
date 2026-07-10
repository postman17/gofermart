package model

type AuthUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
