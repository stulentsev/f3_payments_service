package domain_test

import (
	"github.com/nvloff/f3_payments_service/internal/domain"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Internal/Domain/Commands", func() {
	Describe("CreatePaymentCommand", func() {
		var validCommand *domain.CreatePayment
		var command *domain.CreatePayment
		BeforeEach(func() {
			validCommand = &domain.CreatePayment{
				Type:           "Payment",
				ID:             "ee320377-5f8d-4236-b5d5-9373c135e9e7",
				Version:        0,
				OrganisationID: "03ee559d-1ffe-4a49-82b9-2cf8ed31ecd8",
				Data: domain.PaymentCommandData{
					Amount: "100.21",
					BeneficiaryParty: domain.BeneficiaryPartyCommandData{
						AccountName:       "W Owens",
						AccountNumber:     "31926819",
						AccountNumberCode: "BBAN",
						AccountType:       0,
						Address:           "1 The Beneficiary Localtown SE2",
						BankID:            "403000",
						BankIDCode:        "GBDSC",
						Name:              "Wilfred Jeremiah Owens",
					},
					ChargesInformation: domain.ChargesInformationCommandData{
						BearerCode: "SHAR",
						SenderCharges: []domain.SenderChargesCommandData{
							{Amount: "5.00", Currency: "GBP"},
							{Amount: "10.00", Currency: "USD"},
						},
						ReceiverChargesAmount:   "1.00",
						ReceiverChargesCurrency: "USD",
					},
					Currency: "GBP",
					DebtorParty: domain.DebtorPartyCommandData{
						AccountName:       "EJ Brown Black",
						AccountNumber:     "BG18RZBB91550123456789",
						AccountNumberCode: "IBAN",
						Address:           "10 Debtor Crescent Sourcetown NE1",
						BankID:            "203301",
						BankIDCode:        "GBDSC",
						Name:              "Emelia Jane Brown",
					},
					EndToEndReference: "Wil piano Jan",
					Fx: domain.FXCommandData{
						ContractReference: "FX123",
						ExchangeRate:      "2.00000",
						OriginalAmount:    "200.42",
						OriginalCurrency:  "USD",
					},
					NumericReference:     "1002001",
					PaymentID:            "123456789012345678",
					Purpose:              "Paying for goods/services",
					Scheme:               "FPS",
					PaymentType:          "Credit",
					ProcessingDate:       "2017-01-18",
					Reference:            "Payment for Em's piano lessons",
					SchemePaymentSubType: "InternetBanking",
					SchemePaymentType:    "ImmediatePayment",
					SponsorParty: domain.SponsorPartyCommandData{
						AccountNumber: "56781234",
						BankID:        "123123",
						BankIDCode:    "GBDSC",
					},
				},
			}

			command = validCommand
		})

		Context("Valid", func() {
			It("doesnt return an error", func() {
				err := command.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("missing ID", func() {
			JustBeforeEach(func() { command.ID = "" })

			It("returns an error", func() {
				Expect(command.Validate()).To(HaveOccurred())
			})
		})

		Context("invalid ID", func() {
			JustBeforeEach(func() { command.ID = "inv" })

			It("returns an error", func() {
				Expect(command.Validate()).To(HaveOccurred())
			})
		})

		Context("missing OrganisationID", func() {
			JustBeforeEach(func() { command.OrganisationID = "" })

			It("returns an error", func() {
				Expect(command.Validate()).To(HaveOccurred())
			})
		})

		Context("invalid OrganisationID", func() {
			JustBeforeEach(func() { command.OrganisationID = "inv" })

			It("returns an error", func() {
				Expect(command.Validate()).To(HaveOccurred())
			})
		})

		Context("missing Type", func() {
			JustBeforeEach(func() { command.Type = "" })

			It("returns an error", func() {
				Expect(command.Validate()).To(HaveOccurred())
			})
		})

		Context("invalid Type", func() {
			JustBeforeEach(func() { command.Type = "inv$" })

			It("returns an error", func() {
				Expect(command.Validate()).To(HaveOccurred())
			})
		})

		Context("invalid Version", func() {
			JustBeforeEach(func() { command.Version = -1 })

			It("returns an error", func() {
				Expect(command.Validate()).To(HaveOccurred())
			})
		})

		Context("missing PaymentID", func() {
			JustBeforeEach(func() { command.Data.PaymentID = "" })

			It("returns an error", func() {
				Expect(command.Validate()).To(HaveOccurred())
			})
		})

		Context("missing PaymentType", func() {
			JustBeforeEach(func() { command.Data.PaymentType = "" })

			It("returns an error", func() {
				Expect(command.Validate()).To(HaveOccurred())
			})
		})

		Context("missing Amount", func() {
			JustBeforeEach(func() { command.Data.Amount = "" })

			It("returns an error", func() {
				Expect(command.Validate()).To(HaveOccurred())
			})
		})

		Context("invalid Amount", func() {
			JustBeforeEach(func() { command.Data.Amount = "inv" })

			It("returns an error", func() {
				Expect(command.Validate()).To(HaveOccurred())
			})
		})

		Context("invalid SenderCharges amount", func() {
			JustBeforeEach(func() {
				command.Data.ChargesInformation.SenderCharges = []domain.SenderChargesCommandData{
					{Amount: "5.00", Currency: "GBP"},
					{Amount: "invalid", Currency: "USD"},
				}
			})

			It("returns an error", func() {
				Expect(command.Validate()).To(HaveOccurred())
			})

		})

	})

})
