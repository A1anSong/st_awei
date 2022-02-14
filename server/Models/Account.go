package Models

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model

	An       string `json:"an"`
	Email    string `json:"email"`
	Password string `json:"password"`
	//BIC            string `json:"bic"`
	BankCard string `json:"bankCard"`
	//PcardEmail     string `json:"pcardEmail"`
	Type           string `json:"type"`
	Region         string `json:"region"`
	Auth2fa        string `json:"auth2fa"`
	BackupCode     string `json:"backupCode"`
	PublishableKey string `json:"publishableKey"`
	SecretKey      string `json:"secretKey"`

	// FromStripe
	Acc            string `json:"acc"`
	ChargesEnabled bool   `json:"chargesEnabled"`
	PayoutsEnabled bool   `json:"payoutsEnabled"`

	// ForeignKeys
	Balance  Balance
	Invoices []Invoice
	Payouts  []Payout
}
