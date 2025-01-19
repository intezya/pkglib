package pkglib

import validatorlib "github.com/go-playground/validator/v10"

type validator struct {
	v *validatorlib.Validate
}

func (v validator) Validate(s struct{}) error {
	return v.v.Struct(s)
}

var Validator = validator{
	v: validatorlib.New(),
}
