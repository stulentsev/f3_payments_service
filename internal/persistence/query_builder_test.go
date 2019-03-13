package persistence_test

import (
	"github.com/globalsign/mgo/bson"
	"github.com/nvloff/f3_payments_service/internal/persistence"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryBuilder", func() {
	var (
		builder  persistence.QueryBuilder
		query    bson.M
		currency *string
	)

	BeforeEach(func() {
		builder = &persistence.MongoQueryBuilder{}
	})

	JustBeforeEach(func() {
		builder.Currency(currency)
		query = builder.Build()
	})

	Context("Empty", func() {
		It("returns an empty query", func() {
			Expect(query).To(Equal(bson.M{}))
		})
	})

	Context("With currency", func() {
		BeforeEach(func() {
			c := "whatever"
			currency = &c
		})

		It("returns a currency query", func() {
			Expect(query).To(Equal(bson.M{"attributes.currency": "whatever"}))
		})

	})

})
