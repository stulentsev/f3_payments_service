package api

import (
	"context"

	"github.com/go-openapi/runtime/middleware"
	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	"github.com/nvloff/f3_payments_service/gen/models"
	"github.com/nvloff/f3_payments_service/gen/restapi/operations/payments"
	"github.com/nvloff/f3_payments_service/internal/conversion"
)

// Get Payment by ID
func (api *PaymentsAPI) GetPaymentsID(ctx context.Context, params payments.GetPaymentsIDParams) middleware.Responder {
	id := uuid.MustParse(params.ID.String())

	entity, err := api.PaymentRepo.Find(context.Background(), id)
	// 404 on not found
	if rrErr, ok := err.(eh.RepoError); ok && rrErr.Err == eh.ErrEntityNotFound {
		return payments.NewGetPaymentsIDNotFound().WithPayload(NewApiError(err))
	}

	if err != nil {
		// 500 on any other error
		return payments.NewGetPaymentsIDDefault(500).WithPayload(NewApiError(err))
	}

	// Convert the Repo Entity to HTTP Response data
	apiModel := &models.Payment{}
	err = conversion.Map(apiModel, entity)

	if err != nil {
		return payments.NewGetPaymentsIDDefault(500).WithPayload(NewApiError(err))
	}

	res := &models.PaymentDetailsResponse{}
	res.Data = apiModel

	return payments.NewGetPaymentsIDOK().WithPayload(res)
}
