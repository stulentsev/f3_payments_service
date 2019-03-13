package persistence

import "github.com/google/uuid"

// Payment read model
type Payment struct {
	ID             uuid.UUID  `json:"id" bson:"_id"`
	Type           string     `json:"type" bson:"type"`
	Version        int        `json:"version" bson:"version"`
	OrganisationID uuid.UUID  `json:"organisation_id" bson:"organisation_id"`
	Attributes     Attributes `json:"attributes" bson:"attributes"`
}

// EntityID implements EventHorizon Entity Interface
func (p *Payment) EntityID() uuid.UUID { return p.ID }

// AggregateVersion implements EventHorizon Versioned Entity Interface
func (p *Payment) AggregateVersion() int { return p.Version }

// Attributes embedded document in Payment
type Attributes struct {
	Money                `bson:",inline"`
	BeneficiaryParty     Party              `json:"beneficiary_party" bson:"beneficiary_party"`
	DebtorParty          Party              `json:"debtor_party" bson:"debtor_party"`
	ChargesInformation   ChargesInformation `json:"charges_information" bson:"charges_information"`
	EndToEndReference    string             `json:"end_to_end_reference" bson:"end_to_end_reference"`
	FX                   FX                 `json:"fx" bson:"fx"`
	NumericReference     string             `json:"numeric_reference" bson:"numeric_reference"`
	PaymentID            string             `json:"payment_id" bson:"payment_id"`
	PaymentPurpose       string             `json:"payment_purpose" bson:"payment_purpose"`
	PaymentScheme        string             `json:"payment_scheme" bson:"payment_scheme"`
	PaymentType          string             `json:"payment_type" bson:"payment_type"`
	ProcessingDate       string             `json:"processing_date" bson:"processing_date"`
	Reference            string             `json:"reference" bson:"reference"`
	SchemePaymentSubType string             `json:"scheme_payment_sub_type" bson:"scheme_payment_sub_type"`
	SchemePaymentType    string             `json:"scheme_payment_type" bson:"scheme_payment_type"`
	SponsorParty         SponsorParty       `json:"sponsor_party" bson:"sponsor_party"`
}

// Party embedded document
type Party struct {
	AccountName       string `json:"account_name" bson:"account_name"`
	AccountNumber     string `json:"account_number" bson:"account_number"`
	AccountNumberCode string `json:"account_number_code" bson:"account_number_code"`
	AccountType       int    `json:"account_type" bson:"account_type"`
	Address           string `json:"address" bson:"address"`
	BankID            string `json:"bank_id" bson:"bank_id"`
	BankIDCode        string `json:"bank_id_code" bson:"bank_id_code"`
	Name              string `json:"name" bson:"name"`
}

// ChargesInformation embedded document
type ChargesInformation struct {
	BearerCode              string  `json:"bearer_code" bson:"bearer_code"`
	SenderCharges           []Money `json:"sender_charges" bson:"sender_charges"`
	ReceiverChargesAmount   string  `json:"receiver_charges_amount" bson:"receiver_charges_amount"`
	ReceiverChargesCurrency string  `json:"receiver_charges_currency" bson:"receiver_charges_currency"`
}

// FX embedded document
type FX struct {
	ContractReference string `json:"contract_reference" bson:"contract_reference"`
	ExchangeRate      string `json:"exchange_rate" bson:"exchange_rate"`
	OriginalAmount    string `json:"original_amount" bson:"original_amount"`
	OriginalCurrency  string `json:"original_currency" bson:"original_currency"`
}

// SponsorParty embedded document
type SponsorParty struct {
	AccountNumber string `json:"account_number" bson:"account_number"`
	BankID        string `json:"bank_id" bson:"bank_id"`
	BankIDCode    string `json:"bank_id_code" bson:"bank_id_code"`
}

// Money embedded document
type Money struct {
	Amount   string `json:"amount" bson:"amount"`
	Currency string `json:"currency" bson:"currency"`
}
