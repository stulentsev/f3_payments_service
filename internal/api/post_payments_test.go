package api_test

import (
	"bytes"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("POST payments", func() {
	const fixtureFile = "create_payment.json"
	const invalidFixtureFile = "create_payment_invalid.json"
	const fixturePaymentId = "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43"

	var (
		recorder *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
	})

	Context("Empty Request", func() {
		JustBeforeEach(func() {
			recorder = postPayment(nil)
		})

		It("returns 400", func() {
			Expect(recorder.Code).To(Equal(400))
		})

	})

	Context("Invalid Request Structure", func() {
		JustBeforeEach(func() {
			recorder = postPayment(bytes.NewBufferString(`{"data":{}}`))
		})

		It("returns 422", func() {
			Expect(recorder.Code).To(Equal(422))
		})

		It("returns json errors", func() {
			errorsJson := mapFromJSON(recorder.Body.Bytes())

			Expect(errorsJson["code"]).To(Equal(float64(602)))
			Expect(errorsJson["message"]).To(Equal("attributes in body is required"))
		})
	})

	Context("Invalid Request Data", func() {
		JustBeforeEach(func() {
			recorder = postPayment(bytes.NewReader(helperLoadBytes(invalidFixtureFile)))
		})

		It("returns 422", func() {
			Expect(recorder.Code).To(Equal(400))
		})

		It("returns json errors", func() {
			errorsJson := mapFromJSON(recorder.Body.Bytes())

			Expect(errorsJson["error_message"]).To(Equal("type: must be a valid value."))
		})
	})

	Context("Valid Request", func() {
		JustBeforeEach(func() {
			recorder = postPayment(bytes.NewReader(helperLoadBytes(fixtureFile)))
		})

		JustAfterEach(func() {
			waitForPaymentProjection(fixturePaymentId)
		})

		It("returns 201", func() {
			Expect(recorder.Code).To(Equal(201))

		})

		It("returns the created payment", func() {
			Expect(recorder.Code).To(Equal(201))

			Expect(recorder.Body.String()).To(MatchJSON(helperLoadBytes(fixtureFile)))
		})

	})

	Context("Invalid domain", func() {
		var (
			recorder1 *httptest.ResponseRecorder
			recorder2 *httptest.ResponseRecorder
		)

		JustBeforeEach(func() {
			recorder1 = postPayment(bytes.NewReader(helperLoadBytes(fixtureFile)))
			recorder2 = postPayment(bytes.NewReader(helperLoadBytes(fixtureFile)))
		})

		It("Wraps the domain error", func() {
			Expect(recorder1.Code).To(Equal(201))
			Expect(recorder2.Code).To(Equal(400))

			errorsJson := mapFromJSON(recorder2.Body.Bytes())
			Expect(errorsJson["error_message"]).To(Equal("payment already created"))
		})
	})
})
