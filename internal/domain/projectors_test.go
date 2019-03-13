package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	"github.com/nvloff/f3_payments_service/internal/persistence"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Dummy test entity
type genericEntity struct{}

// Implement EventHorizon Interface
func (p *genericEntity) EntityID() uuid.UUID   { return uuid.New() }
func (p *genericEntity) AggregateVersion() int { return 0 }

var _ = Describe("Aggregate", func() {
	var (
		model     eh.Entity
		initial   eh.Entity
		event     eh.Event
		id        uuid.UUID
		err       error
		projector *PaymentProjector
	)

	BeforeEach(func() {
		TimeNow = func() time.Time {
			return time.Date(2017, time.July, 10, 23, 0, 0, 0, time.Local)
		}
		id = uuid.New()
		projector = &PaymentProjector{}

	})

	JustBeforeEach(func() {
		model, err = projector.Project(context.Background(), event, initial)
	})

	Describe("Unhandled Event", func() {
		BeforeEach(func() {
			initial = &persistence.Payment{}

			event = eh.NewEventForAggregate(
				eh.EventType("unknown"),
				nil,
				TimeNow(), PaymentAggregateType, id, 1)
		})

		It("returns error", func() {
			Expect(err).To(MatchError("could not project event: unknown"))
		})

		It("returns the model", func() {
			Expect(model).To(Equal(initial))
		})
	})

	Describe("Wrong model type", func() {
		BeforeEach(func() {
			initial = &genericEntity{}

			event = eh.NewEventForAggregate(
				eh.EventType("unknown"),
				nil,
				TimeNow(), PaymentAggregateType, id, 1)
		})

		It("returns error", func() {
			Expect(err).To(MatchError("model is of incorrect type"))
		})

		It("returns the model", func() {
			Expect(model).To(Equal(initial))
		})
	})

	Describe("Creation", func() {
		BeforeEach(func() {
			initial = &persistence.Payment{}

			validEvent := eh.NewEventForAggregate(
				PaymentCreatedEvent,
				&PaymentCreatedData{
					ID:             id.String(),
					OrganisationID: id.String(),
				},
				TimeNow(), PaymentAggregateType, id, 1)

			event = validEvent
		})

		Context("valid", func() {
			It("projects", func() {
				Expect(model).To(Equal(&persistence.Payment{
					ID:             id,
					OrganisationID: id,
					Version:        1,
				}))
			})

			It("doesnt return error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("Invalid Event Data", func() {
			BeforeEach(func() {
				event = eh.NewEventForAggregate(
					PaymentCreatedEvent,
					&PaymentCreatedData{},
					TimeNow(), PaymentAggregateType, id, 1)

			})

			It("returns error", func() {
				Expect(err).To(MatchError("could not map event to model: invalid UUID length: 0"))
			})

			It("returns nil model", func() {
				Expect(model).To(BeNil())
			})
		})

		Context("Invalid Event Data Structure", func() {
			BeforeEach(func() {
				eventData := struct{}{}

				event = eh.NewEventForAggregate(
					PaymentCreatedEvent,
					eventData,
					TimeNow(), PaymentAggregateType, id, 1)

			})

			It("returns error", func() {
				Expect(err).To(MatchError("invalid event data type: {}"))
			})

			It("returns nil model", func() {
				Expect(model).To(BeNil())
			})
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			initial = &persistence.Payment{
				ID:             id,
				OrganisationID: id,
				Version:        1,
				Attributes: persistence.Attributes{
					Money: persistence.Money{
						Amount:   "3.0",
						Currency: "EUR",
					},
				},
			}

			validEvent := eh.NewEventForAggregate(
				PaymentUpdatedEvent,
				&PaymentUpdatedData{
					ID: id.String(),
					Attributes: AttributesEventData{
						MoneyEventData: MoneyEventData{
							Currency: "USD",
						},
					},
				},
				TimeNow(), PaymentAggregateType, id, 1)

			event = validEvent
		})

		Context("valid", func() {

			It("projects", func() {
				Expect(model).To(Equal(&persistence.Payment{
					ID:             id,
					OrganisationID: id,
					Version:        2,
					Attributes: persistence.Attributes{
						Money: persistence.Money{
							Amount:   "3.0",
							Currency: "USD",
						},
					},
				}))
			})

			It("doesnt return error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("Invalid Event Data", func() {
			BeforeEach(func() {
				event = eh.NewEventForAggregate(
					PaymentUpdatedEvent,
					&PaymentUpdatedData{},
					TimeNow(), PaymentAggregateType, id, 1)

			})

			It("returns error", func() {
				Expect(err).To(MatchError("could not map event to model: invalid UUID length: 0"))
			})

			It("returns the model", func() {
				Expect(model).To(Equal(initial))
			})
		})

		Context("Invalid Event Data Structure", func() {
			BeforeEach(func() {
				eventData := struct{}{}

				event = eh.NewEventForAggregate(
					PaymentUpdatedEvent,
					eventData,
					TimeNow(), PaymentAggregateType, id, 1)

			})

			It("returns error", func() {
				Expect(err).To(MatchError("invalid event data type: {}"))
			})

			It("returns the model", func() {
				Expect(model).To(Equal(initial))
			})
		})
	})

	Describe("Deletion", func() {
		BeforeEach(func() {
			initial = &persistence.Payment{}

			validEvent := eh.NewEventForAggregate(
				PaymentDeletedEvent,
				&PaymentDeletedData{ID: id.String()},
				TimeNow(), PaymentAggregateType, id, 1)

			event = validEvent
		})

		It("doesnt return error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("doesnt return the model", func() {
			Expect(model).To(BeNil())
		})
	})
})
