package domain

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/imdario/mergo"
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/eventhandler/projector"
	"github.com/nvloff/f3_payments_service/internal/conversion"
	"github.com/nvloff/f3_payments_service/internal/persistence"
	"github.com/pkg/errors"
)

// PaymentProjector can project domain events to read model entities
type PaymentProjector struct{}

// Static validations of EventHorizon interface
var _ = eh.Entity(&persistence.Payment{})
var _ = eh.Versionable(&persistence.Payment{})

// NewPaymentProjector builds a new projector
func NewPaymentProjector() *PaymentProjector {
	return &PaymentProjector{}
}

// ProjectorType implements the ProjectorType method of the
// eventhorizon.Projector interface.
func (p *PaymentProjector) ProjectorType() projector.Type {
	return projector.Type(string(PaymentAggregateType) + "_projector")
}

// Project implements the Project method of the eventhorizon.Projector interface.
func (p *PaymentProjector) Project(ctx context.Context,
	event eh.Event, entity eh.Entity) (eh.Entity, error) {

	payment, ok := entity.(*persistence.Payment)
	if !ok {
		// return the original Entity(whatever it is), as we don't want to delete it
		return entity, errors.New("model is of incorrect type")
	}

	switch event.EventType() {
	case PaymentCreatedEvent:
		eventData, ok := event.Data().(*PaymentCreatedData)

		if !ok {
			// we don't have enough data to create the projection so return nil
			return nil, fmt.Errorf("invalid event data type: %v", event.Data())
		}

		modelFromEvent := &persistence.Payment{}
		// cast the event data to a DB model, so that we can assign the attributes with mergo
		err := conversion.Map(modelFromEvent, eventData)

		if err != nil {
			// we don't have enough data to create the projection so return nil
			return nil, errors.Wrap(err, "could not map event to model")
		}

		// now that we have both the event data and the entity casted to the same struct
		// we can update the entity with the event data
		if err := mergo.Merge(&payment.Attributes, modelFromEvent.Attributes); err != nil {
			// we don't have enough data to create the projection so return nil
			return nil, errors.Wrap(err, "could not assign attributes from event")
		}

		payment.ID = uuid.MustParse(eventData.ID)
		payment.OrganisationID = uuid.MustParse(eventData.OrganisationID)
		payment.Type = eventData.Type

		// Set the ID when first created.
		payment.ID = event.AggregateID()
	case PaymentUpdatedEvent:
		eventData, ok := event.Data().(*PaymentUpdatedData)

		if !ok {
			return payment, fmt.Errorf("invalid event data type: %v", event.Data())
		}

		// cast the event data to a DB model, so that we can assign the attributes with mergo
		modelFromEvent := &persistence.Payment{}
		err := conversion.Map(modelFromEvent, eventData)

		if err != nil {
			return payment, errors.Wrap(err, "could not map event to model")
		}

		// now that we have both the event data and the entity casted to the same struct
		// we can update the entity with the event data
		if err := mergo.Merge(&payment.Attributes, modelFromEvent.Attributes, mergo.WithOverride); err != nil {
			return payment, errors.Wrap(err, "could not assign attributes from event")
		}
	case PaymentDeletedEvent:
		// returning nil for the model will delete it from the projection
		return nil, nil
	default:
		// Also return the model here to not delete it.
		return payment, fmt.Errorf("could not project event: %s", event.EventType())
	}

	// Always increment the version and set update time on successful updates.
	payment.Version++
	return payment, nil
}
