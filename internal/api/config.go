package api

import (
	eh "github.com/looplab/eventhorizon"
)

type Config struct {
	// Store domain info like EventBus, CommandBus, etc.
	CommandBus  eh.CommandHandler
	PaymentRepo eh.ReadRepo
}

func New(c Config) *PaymentsAPI {
	return &PaymentsAPI{
		Config: c,
	}
}

type PaymentsAPI struct {
	Config
}
