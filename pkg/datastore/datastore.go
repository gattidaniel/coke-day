package datastore

import "github.com/aws/aws-sdk-go/service/dynamodb/expression"

// Datastore -
type Datastore interface {
	Scan(filt expression.ConditionBuilder, castTo interface{}) error
	Get(id, kind string, castTo interface{}) error
	Store(item interface{}) error
	Delete(id, kind string) error
}
