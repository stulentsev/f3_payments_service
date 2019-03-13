package api

import (
	"context"

	"github.com/go-openapi/runtime/middleware"
	"github.com/nvloff/f3_payments_service/gen/models"
	"github.com/nvloff/f3_payments_service/gen/restapi/operations/payments"
	"github.com/nvloff/f3_payments_service/internal/conversion"
	"github.com/nvloff/f3_payments_service/internal/domain"
)

// Create a new Payment
func (api *PaymentsAPI) PostPayments(ctx context.Context, params payments.PostPaymentsParams) middleware.Responder {
	if params.PaymentCreationRequest == nil {
		return payments.NewPostPaymentsBadRequest()
	}

	// convert the API request to a domain command
	command := &domain.CreatePayment{}
	err := conversion.Map(command, params.PaymentCreationRequest.Data)
	if err != nil {
		return payments.NewPostPaymentsDefault(500).WithPayload(NewApiError(err))
	}

	// Execute the domain command. Pass different context than the request one
	// as the request will terminate before we finish the event handling
	// and we may interrupt the projection
	err = api.CommandBus.HandleCommand(context.Background(), command)

	switch err := err.(type) {
	case nil:
		res := &models.PaymentCreationResponse{}
		res.Data = params.PaymentCreationRequest.Data

		return payments.NewPostPaymentsCreated().WithPayload(res)
	case *domain.ErrAlreadyCreated:
		return payments.NewPostPaymentsBadRequest().WithPayload(NewApiError(err))
	case *domain.ErrCommandValidation:
		return payments.NewPostPaymentsBadRequest().WithPayload(NewApiError(err))
	default:
		return payments.NewPostPaymentsDefault(500).WithPayload(NewApiError(err))
	}
}
