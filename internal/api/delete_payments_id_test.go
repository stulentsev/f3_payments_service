package api_test

import (
	"net/http/httptest"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DELETE payment", func() {
	const fixturePaymentId = "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43"

	var (
		recorder *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		// Create a payment
		createFixturePayment()
	})

	Context("Valid", func() {
		JustBeforeEach(func() {
			recorder = deletePayment(fixturePaymentId)
		})

		It("returns HTTP Deleted", func() {
			Expect(recorder.Code).To(Equal(204))
		})

		It("returns error when fetching the payment", func() {
			waitForPaymentUpdate()
			recorder = getPayment(fixturePaymentId)

			Expect(recorder.Code).To(Equal(404))
		})
	})

	Context("Missing", func() {
		It("returns 404", func() {
			recorder = deletePayment(uuid.New().String())

			Expect(recorder.Code).To(Equal(404))
		})
	})

	Context("Invalid Request", func() {
		It("returns 422", func() {
			recorder = deletePayment("123")

			Expect(recorder.Code).To(Equal(422))
		})
	})

})
