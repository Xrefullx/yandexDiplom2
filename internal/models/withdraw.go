package models

import (
	"fmt"
	"time"
)

type Withdraw struct {
	UserLogin   string
	NumberOrder string      `json:"order"`
	Sum         WithdrawSum `json:"sum"`
	ProcessedAT time.Time   `json:"processed_at"`
}

type WithdrawSum float64

func (ws WithdrawSum) String() string {
	return fmt.Sprintf("%.2f", ws)
}
