package app

import (
	"github.com/ymgyt/happy-developing/hpdev/errors"
	validator "gopkg.in/go-playground/validator.v9"
)

// Validator -
type Validator interface {
	Validate(interface{}) error
}

// MustValidator -
func MustValidator() Validator {
	return &validatorImpl{v: validator.New()}
}

type validatorImpl struct {
	v *validator.Validate
}

// Validate -
func (v *validatorImpl) Validate(s interface{}) error {
	err := v.v.Struct(s)
	if err != nil {
		return errors.InvalidInput("validation error", err)
	}
	return nil
}
