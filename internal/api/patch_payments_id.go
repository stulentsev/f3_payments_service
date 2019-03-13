package api

import (
	"context"

	"github.com/go-openapi/runtime/middleware"
	"github.com/nvloff/f3_payments_service/gen/restapi/operations/payments"
	"github.com/nvloff/f3_payments_service/internal/conversion"
	"github.com/nvloff/f3_payments_service/internal/domain"
)

// Partial Update of Payment with PATCH
func (api *PaymentsAPI) PatchPaymentsID(ctx context.Context, params payments.PatchPaymentsIDParams) middleware.Responder {
	if params.PaymentUpdateRequest == nil {
		return payments.NewPatchPaymentsIDBadRequest()
	}

	// convert the API request to a domain command
	commandData := &domain.PaymentCommandData{}
	err := conversion.Map(commandData, params.PaymentUpdateRequest.Data.Attributes)
	if err != nil {
		return payments.NewPatchPaymentsIDDefault(500)
	}

	command := &domain.UpdatePayment{
		ID:   params.ID.String(),
		Data: *commandData,
	}

	// Execute the domain command. Pass different context than the request one
	// as the request will terminate before we finish the event handling
	// and we may interrupt the projection
	err = api.CommandBus.HandleCommand(context.Background(), command)

	switch err := err.(type) {
	case nil:
		return payments.NewPatchPaymentsIDAccepted()
	case *domain.ErrNotFound:
		return payments.NewPatchPaymentsIDNotFound().WithPayload(NewApiError(err))
	default:
		return payments.NewPatchPaymentsIDDefault(500).WithPayload(NewApiError(err))
	}

}
