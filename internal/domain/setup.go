package domain

import (
	"log"

	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/aggregatestore/events"
	"github.com/looplab/eventhorizon/commandhandler/aggregate"
	"github.com/looplab/eventhorizon/commandhandler/bus"
	"github.com/looplab/eventhorizon/eventhandler/projector"
	"github.com/looplab/eventhorizon/middleware/commandhandler/validator"
	"github.com/nvloff/f3_payments_service/internal/persistence"
)

// Config is the domain config structure
type Config struct {
	EventStore  eh.EventStore
	EventBus    eh.EventBus
	CommandBus  *bus.CommandHandler
	PaymentRepo eh.ReadWriteRepo
	EventID     uuid.UUID
}

// PaymentsDomain is a central domain config used to interact with the domain
type PaymentsDomain struct {
	eventBus eh.EventBus
	eventID  uuid.UUID

	EventStore  eh.EventStore
	PaymentRepo eh.ReadWriteRepo
	CommandBus  *bus.CommandHandler
}

// NewPaymentsDomain initialises payments domain with required
// EventHandlers, EventStore, etc.
func NewPaymentsDomain(c *Config) *PaymentsDomain {
	return &PaymentsDomain{
		eventBus: c.EventBus,
		eventID:  c.EventID,

		EventStore:  c.EventStore,
		CommandBus:  c.CommandBus,
		PaymentRepo: c.PaymentRepo,
	}
}

// Setup all the event horizon infrastructure
func (p *PaymentsDomain) Setup() {
	p.setupLogger()
	p.setupCommandHandler()
	p.setupProjector()
}

func (p *PaymentsDomain) setupLogger() {
	// log all events
	p.eventBus.AddObserver(eh.MatchAny(), &Logger{})
}

func (p *PaymentsDomain) setupCommandHandler() {
	// Create the aggregate repository.
	aggregateStore, err := events.NewAggregateStore(p.EventStore, p.eventBus)
	if err != nil {
		log.Fatalf("could not create aggregate store: %s", err)
	}

	// Create the aggregate command handler and register the commands it handles.
	paymentHandler, err := aggregate.NewCommandHandler(PaymentAggregateType, aggregateStore)
	if err != nil {
		log.Fatalf("could not create command handler: %s", err)
	}

	// log all commands
	validationMiddleware := validator.NewMiddleware()
	commandHandler := eh.UseCommandHandlerMiddleware(paymentHandler, LoggingMiddleware, validationMiddleware)

	err = p.CommandBus.SetHandler(commandHandler, CreatePaymentCommand)
	err = p.CommandBus.SetHandler(commandHandler, UpdatePaymentCommand)
	err = p.CommandBus.SetHandler(commandHandler, DeletePaymentCommand)

	if err != nil {
		log.Fatalf("could not set command handler: %s", err)
	}
}

func (p *PaymentsDomain) setupProjector() {
	paymentsProjector := projector.NewEventHandler(
		NewPaymentProjector(), p.PaymentRepo)

	paymentsProjector.SetEntityFactory(func() eh.Entity { return &persistence.Payment{} })

	// trigger the projector on those events
	p.eventBus.AddHandler(eh.MatchAnyEventOf(
		PaymentCreatedEvent,
		PaymentUpdatedEvent,
		PaymentDeletedEvent,
	), paymentsProjector)
}
