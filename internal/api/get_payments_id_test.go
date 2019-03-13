package api_test

import (
	"net/http/httptest"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GET payment", func() {
	const fixtureFileResponse = "create_payment_response.json"
	const fixturePaymentID = "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43"

	var (
		recorder *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		createFixturePayment()
	})

	Context("Existing", func() {
		It("returns the payment", func() {
			recorder = getPayment(fixturePaymentID)

			Expect(recorder.Code).To(Equal(200))
			Expect(recorder.Body.String()).To(MatchJSON(helperLoadBytes(fixtureFileResponse)))
		})
	})

	Context("Missing", func() {
		It("returns 404", func() {
			recorder = getPayment(uuid.New().String())

			Expect(recorder.Code).To(Equal(404))
		})
	})

	Context("Invalid Request", func() {
		It("returns 422", func() {
			recorder = getPayment("123")

			Expect(recorder.Code).To(Equal(422))
		})
	})

})
