package domain

import (
	eh "github.com/looplab/eventhorizon"
)

const (
	// PaymentCreatedEvent is the event after a payment is created
	PaymentCreatedEvent = eh.EventType("payment:created")
	// PaymentUpdatedEvent is the event after a payment is updated
	PaymentUpdatedEvent = eh.EventType("payment:updated")
	// PaymentDeletedEvent is the event after a payment is deleted
	PaymentDeletedEvent = eh.EventType("payment:deleted")
)

// Register the events with EventHorizon
func init() {
	eh.RegisterEventData(PaymentCreatedEvent, func() eh.EventData {
		return &PaymentCreatedData{}
	})
	eh.RegisterEventData(PaymentUpdatedEvent, func() eh.EventData {
		return &PaymentUpdatedData{}
	})
	eh.RegisterEventData(PaymentDeletedEvent, func() eh.EventData {
		return &PaymentDeletedData{}
	})
}

// PaymentCreatedData stores event data for the Created Event
type PaymentCreatedData struct {
	ID             string              `json:"id" bson:"_id"`
	Type           string              `json:"type" bson:"type"`
	Version        int                 `json:"version" bson:"version"`
	OrganisationID string              `json:"organisation_id" bson:"organisation_id"`
	Attributes     AttributesEventData `json:"attributes" bson:"attributes"`
}

// PaymentUpdatedData stores event data for the Updated Event
type PaymentUpdatedData struct {
	ID         string              `json:"id" bson:"_id"`
	Attributes AttributesEventData `json:"attributes" bson:"attributes"`
}

// PaymentDeletedData stores event data for the Deleted Event
type PaymentDeletedData struct {
	ID string `json:"id" bson:"_id"`
}

// AttributesEventData stores shared event data for payment events
type AttributesEventData struct {
	MoneyEventData       `bson:",inline"`
	BeneficiaryParty     PartyEventData              `json:"beneficiary_party" bson:"beneficiary_party"`
	DebtorParty          PartyEventData              `json:"debtor_party" bson:"debtor_party"`
	ChargesInformation   ChargesInformationEventData `json:"charges_information" bson:"charges_information"`
	EndToEndReference    string                      `json:"end_to_end_reference" bson:"end_to_end_reference"`
	FX                   FXEventData                 `json:"fx" bson:"fx"`
	NumericReference     string                      `json:"numeric_reference" bson:"numeric_reference"`
	PaymentID            string                      `json:"payment_id" bson:"payment_id"`
	PaymentPurpose       string                      `json:"payment_purpose" bson:"payment_purpose"`
	PaymentScheme        string                      `json:"payment_scheme" bson:"payment_scheme"`
	PaymentType          string                      `json:"payment_type" bson:"payment_type"`
	ProcessingDate       string                      `json:"processing_date" bson:"processing_date"`
	Reference            string                      `json:"reference" bson:"reference"`
	SchemePaymentSubType string                      `json:"scheme_payment_sub_type" bson:"scheme_payment_sub_type"`
	SchemePaymentType    string                      `json:"scheme_payment_type" bson:"scheme_payment_type"`
	SponsorParty         SponsorPartyEventData       `json:"sponsor_party" bson:"sponsor_party"`
}

// PartyEventData stores shared event data for payment events
type PartyEventData struct {
	AccountName       string `json:"account_name" bson:"account_name"`
	AccountNumber     string `json:"account_number" bson:"account_number"`
	AccountNumberCode string `json:"account_number_code" bson:"account_number_code"`
	AccountType       int    `json:"account_type" bson:"account_type"`
	Address           string `json:"address" bson:"address"`
	BankID            string `json:"bank_id" bson:"bank_id"`
	BankIDCode        string `json:"bank_id_code" bson:"bank_id_code"`
	Name              string `json:"name" bson:"name"`
}

// ChargesInformationEventData stores shared event data for payment events
type ChargesInformationEventData struct {
	BearerCode              string           `json:"bearer_code" bson:"bearer_code"`
	SenderCharges           []MoneyEventData `json:"sender_charges" bson:"sender_charges"`
	ReceiverChargesAmount   string           `json:"receiver_charges_amount" bson:"receiver_charges_amount"`
	ReceiverChargesCurrency string           `json:"receiver_charges_currency" bson:"receiver_charges_currency"`
}

// FXEventData stores shared event data for payment events
type FXEventData struct {
	ContractReference string `json:"contract_reference" bson:"contract_reference"`
	ExchangeRate      string `json:"exchange_rate" bson:"exchange_rate"`
	OriginalAmount    string `json:"original_amount" bson:"original_amount"`
	OriginalCurrency  string `json:"original_currency" bson:"original_currency"`
}

// SponsorPartyEventData stores shared event data for payment events
type SponsorPartyEventData struct {
	AccountNumber string `json:"account_number" bson:"account_number"`
	BankID        string `json:"bank_id" bson:"bank_id"`
	BankIDCode    string `json:"bank_id_code" bson:"bank_id_code"`
}

// MoneyEventData stores shared event data for payment events
type MoneyEventData struct {
	Amount   string `json:"amount" bson:"amount"`
	Currency string `json:"currency" bson:"currency"`
}
