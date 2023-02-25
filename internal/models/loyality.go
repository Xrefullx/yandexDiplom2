package models

import (
	"fmt"
)

type Loyalty struct {
	Status      string  `json:"status"`
	Accrual     Accrual `json:"accrual,omitempty"`
	NumberOrder string  `json:"order"`
}

type Accrual float64

func (ws Accrual) String() string {
	return fmt.Sprintf("%.2f", ws)
}
