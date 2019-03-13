package api_test

import (
	"bytes"
	"net/http/httptest"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PATCH payment", func() {
	const fixturePaymentId = "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43"

	var (
		recorder *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		createFixturePayment()
	})

	Context("Valid", func() {
		JustBeforeEach(func() {
			recorder = updatePayment(fixturePaymentId, bytes.NewBufferString(`{"data":{"attributes":{"currency":"EUR"}}}`))
		})

		It("returns HTTP Accepted", func() {
			Expect(recorder.Code).To(Equal(202))
		})

		It("updates the payment version", func() {
			recorder = getPayment(fixturePaymentId)

			Expect(recorder.Code).To(Equal(200))
			Expect(recorder.Body.String()).To(ContainSubstring(`"version":2`))
		})

		It("updates the payment attributes", func() {
			recorder = getPayment(fixturePaymentId)

			Expect(recorder.Code).To(Equal(200))
			Expect(recorder.Body.String()).To(ContainSubstring(`"currency":"EUR"`))
		})
	})

	Context("Missing", func() {
		It("returns 404", func() {
			recorder = updatePayment(uuid.New().String(), bytes.NewBufferString(`{"data":{"attributes":{"currency":"EUR"}}}`))

			Expect(recorder.Code).To(Equal(404))
		})
	})

	Context("Empty Request", func() {
		It("returns 400", func() {
			recorder = updatePayment(uuid.New().String(), nil)

			Expect(recorder.Code).To(Equal(400))
		})
	})

	Context("Invalid Request", func() {
		It("returns 422", func() {
			recorder = updatePayment(uuid.New().String(), bytes.NewBufferString(`{"data":{"whatever":{"else":"EUR"}}}`))

			Expect(recorder.Code).To(Equal(422))
		})
	})

})
