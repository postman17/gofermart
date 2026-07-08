package model

import "time"

type Withdraw struct {
	Order    string `json:"order"`
	TotalSum int64  `json:"sum"`
}

type WithdrawalsModel struct {
	ID        int64     `json:"-"`
	UserID    int64     `json:"-"`
	Order     string    `json:"order"`
	TotalSum  int64     `json:"sum"`
	CreatedAt time.Time `json:"processed_at"`
}
