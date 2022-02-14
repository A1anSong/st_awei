package Models

import (
	"gorm.io/gorm"
	"time"
)

type Invoice struct {
	gorm.Model

	// FromStripe
	In               string    `json:"in"`
	AutoAdvance      bool      `json:"autoAdvance"`
	CollectionMethod string    `json:"collectionMethod"`
	Currency         string    `json:"currency"`
	Description      string    `json:"description"`
	HostedInvoiceURL string    `json:"hostedInvoiceURL"`
	PeriodEnd        time.Time `json:"periodEnd"`
	PeriodStart      time.Time `json:"periodStart"`
	Status           string    `json:"status"`
	Total            int       `json:"total"`

	// ForeignKeys
	AccountID uint
}
