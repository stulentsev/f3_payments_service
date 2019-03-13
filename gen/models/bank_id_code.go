// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/validate"
)

// BankIDCode bank Id code
// swagger:model BankIdCode
type BankIDCode string

const (

	// BankIDCodeGBDSC captures enum value "GBDSC"
	BankIDCodeGBDSC BankIDCode = "GBDSC"
)

// for schema
var bankIdCodeEnum []interface{}

func init() {
	var res []BankIDCode
	if err := json.Unmarshal([]byte(`["GBDSC"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		bankIdCodeEnum = append(bankIdCodeEnum, v)
	}
}

func (m BankIDCode) validateBankIDCodeEnum(path, location string, value BankIDCode) error {
	if err := validate.Enum(path, location, value, bankIdCodeEnum); err != nil {
		return err
	}
	return nil
}

// Validate validates this bank Id code
func (m BankIDCode) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateBankIDCodeEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
