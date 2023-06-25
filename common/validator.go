package common

import playground "github.com/go-playground/validator/v10"

func ValidateRequest(s interface{}) error {
	validate := playground.New()

	err := validate.Struct(s)
	if err != nil {
		return TransformValidator(err)
	}

	return nil
}