package models

import (
	"fmt"
	"time"
)

type Order struct {
	UserLogin   string
	Status      string     `json:"status"`
	Accrual     AccrualSum `json:"accrual,omitempty"`
	Uploaded    time.Time  `json:"uploaded_at"`
	NumberOrder string     `json:"number"`
}

type AccrualSum float64

func (ws AccrualSum) String() string {
	return fmt.Sprintf("%.2f", ws)
}
