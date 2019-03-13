// Code generated by go-swagger; DO NOT EDIT.

package payments

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/nvloff/f3_payments_service/gen/models"
)

// GetPaymentsOKCode is the HTTP code returned for type GetPaymentsOK
const GetPaymentsOKCode int = 200

/*GetPaymentsOK List of payment details

swagger:response getPaymentsOK
*/
type GetPaymentsOK struct {

	/*
	  In: Body
	*/
	Payload *models.PaymentDetailsListResponse `json:"body,omitempty"`
}

// NewGetPaymentsOK creates GetPaymentsOK with default headers values
func NewGetPaymentsOK() *GetPaymentsOK {

	return &GetPaymentsOK{}
}

// WithPayload adds the payload to the get payments o k response
func (o *GetPaymentsOK) WithPayload(payload *models.PaymentDetailsListResponse) *GetPaymentsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get payments o k response
func (o *GetPaymentsOK) SetPayload(payload *models.PaymentDetailsListResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPaymentsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetPaymentsDefault Unexpected error

swagger:response getPaymentsDefault
*/
type GetPaymentsDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.APIError `json:"body,omitempty"`
}

// NewGetPaymentsDefault creates GetPaymentsDefault with default headers values
func NewGetPaymentsDefault(code int) *GetPaymentsDefault {
	if code <= 0 {
		code = 500
	}

	return &GetPaymentsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get payments default response
func (o *GetPaymentsDefault) WithStatusCode(code int) *GetPaymentsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get payments default response
func (o *GetPaymentsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get payments default response
func (o *GetPaymentsDefault) WithPayload(payload *models.APIError) *GetPaymentsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get payments default response
func (o *GetPaymentsDefault) SetPayload(payload *models.APIError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPaymentsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
