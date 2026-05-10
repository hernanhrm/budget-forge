package validation

import (
	"context"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

var (
	By            = validation.By
	Required      = validation.Required
	NotNil        = validation.NotNil
	NilOrNotEmpty = validation.NilOrNotEmpty
	Nil           = validation.Nil
	Empty         = validation.Empty
	Skip          = validation.Skip
	In            = validation.In
	NotIn         = validation.NotIn
	Length        = validation.Length
	RuneLength    = validation.RuneLength
	Min           = validation.Min
	Max           = validation.Max
	Match         = validation.Match
	Date          = validation.Date
	Each          = validation.Each
	When          = validation.When
)

var (
	IsEmail        = is.Email
	IsEmailFormat  = is.EmailFormat
	IsURL          = is.URL
	IsRequestURL   = is.RequestURL
	IsRequestURI   = is.RequestURI
	IsAlpha        = is.Alpha
	IsDigit        = is.Digit
	IsAlphanumeric = is.Alphanumeric
	IsUUID         = is.UUID
	IsUUIDv4       = is.UUIDv4
	IsInt          = is.Int
	IsFloat        = is.Float
	IsLowerCase    = is.LowerCase
	IsUpperCase    = is.UpperCase
)

type Validatable interface {
	Validate(ctx context.Context) error
}

func Validate(ctx context.Context, v Validatable) error {
	if err := v.Validate(ctx); err != nil {
		return fmt.Errorf("validation: %w", err)
	}
	return nil
}

func ValidateStruct(ctx context.Context, structPtr any, fields ...*validation.FieldRules) error {
	if err := validation.ValidateStruct(structPtr, fields...); err != nil {
		return fmt.Errorf("validation: %w", err)
	}
	return nil
}

func Field(fieldPtr any, rules ...validation.Rule) *validation.FieldRules {
	return validation.Field(fieldPtr, rules...)
}

type Rules []validation.Rule
