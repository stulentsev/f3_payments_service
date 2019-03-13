package api

import (
	"context"

	"github.com/nvloff/f3_payments_service/internal/api/services"

	"github.com/go-openapi/runtime/middleware"
	"github.com/nvloff/f3_payments_service/gen/models"
	"github.com/nvloff/f3_payments_service/gen/restapi/operations/payments"
	"github.com/nvloff/f3_payments_service/internal/persistence"
)

func (api *PaymentsAPI) GetPayments(ctx context.Context, params payments.GetPaymentsParams) middleware.Responder {
	queryBuilder := persistence.MongoQueryBuilder{}
	queryBuilder.Currency(params.FilterCurrency)

	queryObject := services.PaymentsQuery{
		RepoRepo: api.PaymentRepo,
		Query:    queryBuilder.Build(),
	}

	apiPayments, err := queryObject.Execute()

	if err != nil {
		return payments.NewGetPaymentsIDDefault(500).WithPayload(NewApiError(err))
	}

	return payments.NewGetPaymentsOK().
		WithPayload(&models.PaymentDetailsListResponse{Data: apiPayments})
}
