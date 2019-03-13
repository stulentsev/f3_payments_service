package api_test

import (
	"net/http/httptest"

	"github.com/nvloff/f3_payments_service/gen/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GET payments", func() {
	const payment1Id = "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43"
	const payment2Id = "1f10bc16-ad79-4243-9700-a9f3b6ce6084"

	var (
		recorder *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		createPayment("payment1.json", payment1Id)
		createPayment("payment2.json", payment2Id)
	})

	Context("Base", func() {

		It("returns payments", func() {
			recorder = getPayments("")

			Expect(recorder.Code).To(Equal(200))

			response := &models.PaymentDetailsListResponse{}

			err := response.UnmarshalBinary(recorder.Body.Bytes())
			Expect(err).NotTo(HaveOccurred())

			Expect(len(response.Data)).To(Equal(2))
			Expect(response.Data[0].ID.String()).To(Equal(payment1Id))
			Expect(response.Data[1].ID.String()).To(Equal(payment2Id))
		})
	})

	Context("Filter Currency", func() {

		It("returns payments", func() {
			recorder = getPayments("?filter[currency]=USD")

			Expect(recorder.Code).To(Equal(200))

			response := &models.PaymentDetailsListResponse{}

			err := response.UnmarshalBinary(recorder.Body.Bytes())
			Expect(err).NotTo(HaveOccurred())

			Expect(len(response.Data)).To(Equal(1))
			Expect(response.Data[0].ID.String()).To(Equal(payment2Id))
		})
	})
	Context("Filter Empty Result", func() {

		It("returns empty response", func() {
			recorder = getPayments("?filter[currency]=XXX")

			Expect(recorder.Code).To(Equal(200))

			response := &models.PaymentDetailsListResponse{}

			err := response.UnmarshalBinary(recorder.Body.Bytes())
			Expect(err).NotTo(HaveOccurred())

			Expect(len(response.Data)).To(Equal(0))
		})
	})
})
