package persistence

import (
	"github.com/globalsign/mgo/bson"
)

// QueryBuilder builds DB queries for specific Payment fields
type QueryBuilder interface {
	Currency(currency *string)
	Build() bson.M
}

// MongoQueryBuilder builds a mongo query for specific Payment fields
type MongoQueryBuilder struct {
	currency *string
}

// Currency sets a currency filter
func (b *MongoQueryBuilder) Currency(currency *string) {
	b.currency = currency
}

// Build builds the mongo query
func (b *MongoQueryBuilder) Build() bson.M {
	query := bson.M{}

	if b.currency != nil {
		query["attributes.currency"] = *b.currency
	}

	return query
}
