package domain

import (
	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
)

func init() {
	eh.RegisterCommand(func() eh.Command { return &CreatePayment{} })
	eh.RegisterCommand(func() eh.Command { return &UpdatePayment{} })
	eh.RegisterCommand(func() eh.Command { return &DeletePayment{} })
}

const (
	// CreatePaymentCommand creates a new payment.
	CreatePaymentCommand = eh.CommandType("payments:create")
	// UpdatePaymentCommand updates existing payment
	UpdatePaymentCommand = eh.CommandType("payments:update")
	// DeletePaymentCommand deletes a payment
	DeletePaymentCommand = eh.CommandType("payments:delete")
)

// Static type check that the `eventhorizon.Command` interface is implemented.
var _ = eh.Command(&CreatePayment{})
var _ = eh.Command(&UpdatePayment{})
var _ = eh.Command(&DeletePayment{})

// CreatePayment command data
type CreatePayment struct {
	Type           string             `json:"type"`
	ID             string             `json:"id"`
	Version        int                `json:"version"`
	OrganisationID string             `json:"organisation_id"`
	Data           PaymentCommandData `json:"attributes"`
}

// UpdatePayment command data
type UpdatePayment struct {
	ID   string             `json:"id"`
	Data PaymentCommandData `json:"attributes"`
}

// DeletePayment command data
type DeletePayment struct {
	ID string `json:"id"`
}

//  Implement EventHorizon Command Interface
func (c *CreatePayment) AggregateType() eh.AggregateType { return PaymentAggregateType }
func (c *CreatePayment) AggregateID() uuid.UUID          { return uuid.MustParse(c.ID) }
func (c *CreatePayment) CommandType() eh.CommandType     { return CreatePaymentCommand }

func (c *UpdatePayment) AggregateType() eh.AggregateType { return PaymentAggregateType }
func (c *UpdatePayment) AggregateID() uuid.UUID          { return uuid.MustParse(c.ID) }
func (c *UpdatePayment) CommandType() eh.CommandType     { return UpdatePaymentCommand }

func (c *DeletePayment) AggregateType() eh.AggregateType { return PaymentAggregateType }
func (c *DeletePayment) AggregateID() uuid.UUID          { return uuid.MustParse(c.ID) }
func (c *DeletePayment) CommandType() eh.CommandType     { return DeletePaymentCommand }

// PaymentCommandData contains details about a payment command
type PaymentCommandData struct {
	Amount               string                        `json:"amount"`
	BeneficiaryParty     BeneficiaryPartyCommandData   `json:"beneficiary_party"`
	ChargesInformation   ChargesInformationCommandData `json:"charges_information"`
	Currency             string                        `json:"currency"`
	DebtorParty          DebtorPartyCommandData        `json:"debtor_party"`
	EndToEndReference    string                        `json:"end_to_end_reference"`
	Fx                   FXCommandData                 `json:"fx"`
	NumericReference     string                        `json:"numeric_reference"`
	PaymentID            string                        `json:"payment_id"`
	Purpose              string                        `json:"payment_purpose"`
	Scheme               string                        `json:"payment_scheme"`
	PaymentType          string                        `json:"payment_type"`
	ProcessingDate       string                        `json:"processing_date"`
	Reference            string                        `json:"reference"`
	SchemePaymentSubType string                        `json:"scheme_payment_sub_type"`
	SchemePaymentType    string                        `json:"scheme_payment_type"`
	SponsorParty         SponsorPartyCommandData       `json:"sponsor_party"`
}

// ChargesInformationCommandData contains details about payment command attributes
type ChargesInformationCommandData struct {
	BearerCode              string                     `json:"bearer_code"`
	SenderCharges           []SenderChargesCommandData `json:"sender_charges"`
	ReceiverChargesAmount   string                     `json:"receiver_charges_amount"`
	ReceiverChargesCurrency string                     `json:"receiver_charges_currency"`
}

// SenderChargesCommandData contains details about payment command attributes
type SenderChargesCommandData struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

// BeneficiaryPartyCommandData contains details about payment command attributes
type BeneficiaryPartyCommandData struct {
	AccountName       string `json:"account_name"`
	AccountNumber     string `json:"account_number"`
	AccountNumberCode string `json:"account_number_code"`
	AccountType       int    `json:"account_type"`
	Address           string `json:"address"`
	BankID            string `json:"bank_id"`
	BankIDCode        string `json:"bank_id_code"`
	Name              string `json:"name"`
}

// DebtorPartyCommandData contains details about payment command attributes
type DebtorPartyCommandData struct {
	AccountName       string `json:"account_name"`
	AccountNumber     string `json:"account_number"`
	AccountNumberCode string `json:"account_number_code"`
	Address           string `json:"address"`
	BankID            string `json:"bank_id"`
	BankIDCode        string `json:"bank_id_code"`
	Name              string `json:"name"`
}

// FXCommandData contains details about payment command attributes
type FXCommandData struct {
	ContractReference string `json:"contract_reference"`
	ExchangeRate      string `json:"exchange_rate"`
	OriginalAmount    string `json:"original_amount"`
	OriginalCurrency  string `json:"original_currency"`
}

// SponsorPartyCommandData contains details about payment command attributes
type SponsorPartyCommandData struct {
	AccountNumber string `json:"account_number"`
	BankID        string `json:"bank_id"`
	BankIDCode    string `json:"bank_id_code"`
}
