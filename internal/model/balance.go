package model

type Balance struct {
	ID        int64 `json:"id"`
	Current   int64 `json:"current"`
	Withdrawn int64 `json:"withdrawn"`
}
