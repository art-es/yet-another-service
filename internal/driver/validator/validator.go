package validator

import validatorV10 "github.com/go-playground/validator/v10"

func New() *validatorV10.Validate {
	return validatorV10.New()
}
