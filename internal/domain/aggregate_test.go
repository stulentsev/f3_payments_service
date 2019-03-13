package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/aggregatestore/events"
	"github.com/looplab/eventhorizon/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Aggregate", func() {
	var (
		agg            *Aggregate
		cmd            eh.Command
		id             uuid.UUID
		err            error
		expectedEvents []eh.Event
	)

	BeforeEach(func() {
		TimeNow = func() time.Time {
			return time.Date(2017, time.July, 10, 23, 0, 0, 0, time.Local)
		}
		id = uuid.New()

		agg = &Aggregate{
			AggregateBase: events.NewAggregateBase(PaymentAggregateType, id),

			status: initial,
		}
	})

	Describe("Handle command", func() {
		JustBeforeEach(func() {
			err = agg.HandleCommand(context.Background(), cmd)
		})

		Context("unknown command", func() {
			BeforeEach(func() {
				cmd = &mocks.Command{ID: id}
				agg.status = created
			})

			It("returns an error", func() {
				Expect(err).To(MatchError("could not handle command: Command"))
			})
		})

		Context("Create Command", func() {
			BeforeEach(func() {
				cmd = &CreatePayment{ID: id.String()}

				expectedEvents = []eh.Event{eh.NewEventForAggregate(
					PaymentCreatedEvent,
					&PaymentCreatedData{ID: id.String()},
					TimeNow(), PaymentAggregateType, id, 1)}
			})

			It("records the event", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(agg.Events()).To(Equal(expectedEvents))
			})
		})

		Context("Create Command, already created", func() {
			BeforeEach(func() {
				cmd = &CreatePayment{ID: id.String()}
				agg.status = created
			})

			It("returns an error", func() {
				Expect(err).To(MatchError("payment already created"))
			})
		})

		Context("Create Command, already deleted", func() {
			BeforeEach(func() {
				cmd = &CreatePayment{ID: id.String()}
				agg.status = deleted
			})

			It("returns an error", func() {
				Expect(err).To(MatchError("payment already created"))
			})
		})

		Context("Delete Command", func() {
			BeforeEach(func() {
				cmd = &DeletePayment{ID: id.String()}
				agg.status = created

				expectedEvents = []eh.Event{eh.NewEventForAggregate(
					PaymentDeletedEvent,
					PaymentDeletedData{ID: id.String()},
					TimeNow(), PaymentAggregateType, id, 1)}
			})

			It("records the event", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(agg.Events()).To(BeEquivalentTo(expectedEvents))
			})
		})

		Context("Delete Command, already deleted", func() {
			BeforeEach(func() {
				cmd = &DeletePayment{ID: id.String()}
				agg.status = deleted
			})

			It("returns an error", func() {
				Expect(err).To(MatchError("payment not found"))
			})
		})

		Context("Update Command", func() {
			BeforeEach(func() {
				cmd = &UpdatePayment{ID: id.String()}
				agg.status = created

				expectedEvents = []eh.Event{eh.NewEventForAggregate(
					PaymentUpdatedEvent,
					&PaymentUpdatedData{ID: id.String()},
					TimeNow(), PaymentAggregateType, id, 1)}
			})

			It("records the event", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(agg.Events()).To(BeEquivalentTo(expectedEvents))
			})
		})

		Context("Update Command, not created", func() {
			BeforeEach(func() {
				cmd = &UpdatePayment{ID: id.String()}
				agg.status = initial
			})

			It("returns an error", func() {
				Expect(err).To(MatchError("payment not found"))
			})
		})

		Context("Update Command, deleted", func() {
			BeforeEach(func() {
				cmd = &UpdatePayment{ID: id.String()}
				agg.status = deleted
			})

			It("returns an error", func() {
				Expect(err).To(MatchError("payment not found"))
			})
		})
	})

	Describe("Apply Event ", func() {
		var (
			initialAgg  *Aggregate
			expectedAgg *Aggregate
			event       eh.Event
		)

		JustBeforeEach(func() {
			initialAgg = &Aggregate{
				AggregateBase: events.NewAggregateBase(PaymentAggregateType, id),
			}
			err = initialAgg.ApplyEvent(context.Background(), event)
		})

		Context("Unknown event", func() {
			BeforeEach(func() {
				event = eh.NewEventForAggregate(
					eh.EventType("unknown"),
					nil,
					TimeNow(), PaymentAggregateType, id, 1)
			})

			It("returns an error", func() {
				Expect(err).To(MatchError("could not apply event: unknown"))
			})
		})

		Context("Create Event", func() {
			BeforeEach(func() {
				event = eh.NewEventForAggregate(
					PaymentCreatedEvent,
					&PaymentCreatedData{ID: id.String()},
					TimeNow(), PaymentAggregateType, id, 1)
			})

			It("changes the status to created", func() {
				expectedAgg = &Aggregate{
					AggregateBase: events.NewAggregateBase(PaymentAggregateType, id),
					status:        created,
				}

				Expect(initialAgg).To(Equal(expectedAgg))
			})
		})

		Context("Delete Event", func() {
			BeforeEach(func() {
				event = eh.NewEventForAggregate(
					PaymentDeletedEvent,
					&PaymentDeletedData{ID: id.String()},
					TimeNow(), PaymentAggregateType, id, 1)
			})

			It("changes the status to deleted", func() {
				expectedAgg = &Aggregate{
					AggregateBase: events.NewAggregateBase(PaymentAggregateType, id),
					status:        deleted,
				}

				Expect(initialAgg).To(Equal(expectedAgg))
			})
		})
	})
})
