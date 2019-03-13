package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/nvloff/f3_payments_service/internal/conversion"
	"github.com/pkg/errors"

	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/aggregatestore/events"
)

func init() {
	// Register the PaymentAggregate with EventHorizon
	eh.RegisterAggregate(func(id uuid.UUID) eh.Aggregate {
		return &Aggregate{
			AggregateBase: events.NewAggregateBase(PaymentAggregateType, id),
		}
	})
}

// PaymentAggregateType is the aggregate type for the payment.
const PaymentAggregateType = eh.AggregateType("payment")

// ErrAlreadyCreated when we try to create an already created event
type ErrAlreadyCreated struct{}

func (*ErrAlreadyCreated) Error() string {
	return "payment already created"
}

// ErrNotFound
// Domain error if we try to do anything on a payment that is not created
type ErrNotFound struct{}

func (*ErrNotFound) Error() string {
	return "payment not found"
}

type aggregateStatus int

const (
	initial aggregateStatus = iota
	created
	deleted
)

// Aggregate is an aggregate for a payment
type Aggregate struct {
	*events.AggregateBase

	// Only domain value we track is the persistence status of the Payment
	status aggregateStatus
}

// TimeNow is a mockable version of time.Now.
var TimeNow = time.Now

// HandleCommand implements the HandleCommand method of the
// eventhorizon.CommandHandler interface.
func (a *Aggregate) HandleCommand(ctx context.Context, cmd eh.Command) error {
	switch cmd.(type) {
	case *CreatePayment:
		// An aggregate can only be created once.
		if a.status != initial {
			return &ErrAlreadyCreated{}
		}

	default:
		// All other events require the aggregate to be created.
		if a.status != created {
			return &ErrNotFound{}
		}
	}

	switch cmd := cmd.(type) {
	case *CreatePayment:
		// Map the CommandData to EventData and store the event
		eventData := &PaymentCreatedData{}
		err := conversion.Map(eventData, cmd)

		if err != nil {
			return errors.Wrap(err, "Cant convert command data to event data")
		}

		a.StoreEvent(PaymentCreatedEvent, eventData, TimeNow())
	case *UpdatePayment:
		// Map the CommandData to EventData and store the event
		eventData := &PaymentUpdatedData{}
		err := conversion.Map(eventData, cmd)

		if err != nil {
			return errors.Wrap(err, "Cant convert command data to event data")
		}

		eventData.ID = cmd.ID

		a.StoreEvent(PaymentUpdatedEvent, eventData, TimeNow())
	case *DeletePayment:
		eventData := PaymentDeletedData{ID: cmd.ID}

		a.StoreEvent(PaymentDeletedEvent, eventData, TimeNow())
	default:
		return fmt.Errorf("could not handle command: %s", cmd.CommandType())
	}
	return nil
}

// ApplyEvent implements the ApplyEvent method of the
// eventhorizon.Aggregate interface.
func (a *Aggregate) ApplyEvent(ctx context.Context, event eh.Event) error {
	switch event.EventType() {
	case PaymentCreatedEvent:
		a.status = created
	case PaymentUpdatedEvent:
	case PaymentDeletedEvent:
		a.status = deleted
	default:
		return fmt.Errorf("could not apply event: %s", event.EventType())
	}
	return nil
}
