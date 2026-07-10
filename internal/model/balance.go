package model

type Balance struct {
	ID        int64   `json:"id"`
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
