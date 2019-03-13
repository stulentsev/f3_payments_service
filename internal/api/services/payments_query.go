package services

import (
	"context"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	eh "github.com/looplab/eventhorizon"
	repo "github.com/looplab/eventhorizon/repo/mongodb"
	"github.com/nvloff/f3_payments_service/gen/models"
	"github.com/nvloff/f3_payments_service/internal/conversion"
)

type QueryObject interface {
	Execute() ([]*models.Payment, error)
}

type PaymentsQuery struct {
	RepoRepo eh.ReadRepo
	Query    bson.M
}

func (p *PaymentsQuery) Execute() ([]*models.Payment, error) {
	mongoRepo := p.RepoRepo.Parent().(*repo.Repo)

	result, err := mongoRepo.FindCustom(context.Background(), func(c *mgo.Collection) *mgo.Query {
		return c.Find(p.Query)
	})

	if err != nil {
		return nil, err
	}

	return p.mapQueryResultToApiResult(result)
}

func (p *PaymentsQuery) mapQueryResultToApiResult(result []interface{}) ([]*models.Payment, error) {
	apiPayments := make([]*models.Payment, len(result))
	for i, dbPayment := range result {
		apiPayment := &models.Payment{}
		err := conversion.Map(apiPayment, dbPayment)

		if err != nil {
			return nil, err
		}

		apiPayments[i] = apiPayment
	}

	return apiPayments, nil
}
