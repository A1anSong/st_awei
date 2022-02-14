package Models

import "gorm.io/gorm"

type Balance struct {
	gorm.Model

	// FromStripe
	Available int `json:"available"`
	Pending   int `json:"pending"`

	// ForeignKeys
	AccountID uint
}
