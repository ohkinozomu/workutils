package workutils

import "github.com/go-playground/validator/v10"

type Resource struct {
	Group     string
	Version   string `validate:"required"`
	Kind      string `validate:"required"`
	Name      string `validate:"required"`
	Namespace string
}

func validate(resource Resource) error {
	validate := validator.New()
	return validate.Struct(resource)
}
