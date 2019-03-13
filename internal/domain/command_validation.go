package domain

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// ErrCommandValidation signifies a domain error in the command parameters
type ErrCommandValidation struct {
	original error
}

func (e *ErrCommandValidation) Error() string {
	return e.original.Error()
}

// Validate validates a command to the domain
func (c CreatePayment) Validate() error {
	err := c.internalValidate()
	if err != nil {
		return &ErrCommandValidation{err}
	}

	return nil
}

// Example Command validation
// The validation is deliberately separated from the structure definition
// We reuse the same structures for create/update but they have slightly different validations
func (c CreatePayment) internalValidate() error {
	var err error

	// Validate base PaymentData
	err = validation.ValidateStruct(&c,
		validation.Field(&c.ID, validation.Required, is.UUID),
		validation.Field(&c.OrganisationID, validation.Required, is.UUID),
		validation.Field(&c.Type, validation.Required, validation.Match(regexp.MustCompile("^[A-Za-z_]*$")), validation.In("Payment")),
		validation.Field(&c.Version, validation.Min(0)),
	)

	if err != nil {
		return err
	}

	// Validate Attributes
	err = validation.ValidateStruct(&c.Data,
		validation.Field(&c.Data.PaymentID, validation.Required),
		validation.Field(&c.Data.PaymentType, validation.Required),
		validation.Field(&c.Data.Amount, validation.Required, validation.Match(regexp.MustCompile("^[0-9.]{0,20}$"))),
		validation.Field(&c.Data.Currency, validation.Required),
		validation.Field(&c.Data.ChargesInformation, validation.Required),
	)

	if err != nil {
		return err
	}

	// Example of deeply nested validation
	err = validation.ValidateStruct(&c.Data.ChargesInformation,
		validation.Field(&c.Data.ChargesInformation.SenderCharges, validation.Required, validation.Length(1, 100)),
	)

	if err != nil {
		return err
	}
	for _, charges := range c.Data.ChargesInformation.SenderCharges {
		err = validation.ValidateStruct(&charges,
			validation.Field(&charges.Amount, validation.Required, validation.Match(regexp.MustCompile("^[0-9.]{0,20}$"))),
			validation.Field(&charges.Currency, validation.Required),
		)

		if err != nil {
			return err
		}
	}

	return err
}
