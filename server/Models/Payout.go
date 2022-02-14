package Models

import (
	"gorm.io/gorm"
	"time"
)

type Payout struct {
	gorm.Model

	// FromStripe
	Po                  string    `json:"po"`
	Amount              int       `json:"amountL"`
	ArrivalDate         time.Time `json:"arrivalDate"`
	Currency            string    `json:"currency"`
	Description         string    `json:"description"`
	StatementDescriptor string    `json:"statementDescriptor"`
	Status              string    `json:"status"`

	// ForeignKeys
	AccountID uint
}
