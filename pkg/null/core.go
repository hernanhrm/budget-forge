package null

import "github.com/guregu/null/v6"

type String = null.String
type Bool = null.Bool
type Int = null.Int
type Float = null.Float
type Time = null.Time
type Byte = null.Byte

func NewString(value string, valid bool) String {
	if !valid {
		return String{}
	}
	return null.StringFrom(value)
}
