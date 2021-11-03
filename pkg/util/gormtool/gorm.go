// Package gormtool is a util to convert offset and limit to default values.
package gormtool

// DefaultLimit define the default number of records to be retrieved.
const DefaultLimit = 1000

// LimitAndOffset contains offset and limit fields.
type LimitAndOffset struct {
	Offset int
	Limit  int
}

// Unpointer fill LimitAndOffset with default values if offset/limit is nil
// or it will be filled with the passed value.
func Unpointer(offset *int64, limit *int64) *LimitAndOffset {
	var o, l int = 0, DefaultLimit

	if offset != nil {
		o = int(*offset)
	}

	if limit != nil {
		l = int(*limit)
	}

	return &LimitAndOffset{
		Offset: o,
		Limit:  l,
	}
}
