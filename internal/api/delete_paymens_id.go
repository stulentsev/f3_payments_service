package api

import (
	"context"

	"github.com/go-openapi/runtime/middleware"
	"github.com/nvloff/f3_payments_service/gen/restapi/operations/payments"
	"github.com/nvloff/f3_payments_service/internal/domain"
)

// Delete Payment by ID
func (api *PaymentsAPI) DeletePaymentsID(ctx context.Context, params payments.DeletePaymentsIDParams) middleware.Responder {
	command := &domain.DeletePayment{
		ID: params.ID.String(),
	}

	// Execute the domain command. Pass different context than the request one
	// as the request will terminate before we finish the event handling
	// and we may interrupt the projection
	err := api.CommandBus.HandleCommand(context.Background(), command)

	switch err := err.(type) {
	case nil:
		return payments.NewDeletePaymentsIDNoContent()
	case *domain.ErrNotFound:
		return payments.NewDeletePaymentsIDNotFound().WithPayload(NewApiError(err))
	default:
		return payments.NewDeletePaymentsIDDefault(500).WithPayload(NewApiError(err))
	}
}
