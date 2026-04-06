package products

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Custom error messages for validation errors
var customErrorMessages = map[string]string{
	"Name.required":     "name is required",
	"Name.min":          "name must be at least 3 characters",
	"Description.min":   "description must be at least 10 characters",
	"Price.required":    "price is required",
	"Price.gte":         "price must be at least 0",
	"Category.required": "category is required",
	"Stock.required":    "stock is required",
	"Stock.gte":         "stock must be at least 0",
}

func (p *Product) Validate() error {
	v := validator.New()
	if err := v.Struct(p); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			var errs []string
			for _, e := range ve {
				field := strings.Split(e.Namespace(), ".")[len(strings.Split(e.Namespace(), "."))-1]
				tag := e.Tag()
				key := field + "." + tag
				fieldKey := strings.Split(key, ".")[1]
				if msg, found := customErrorMessages[key]; found {
					errs = append(errs, msg)
				} else {
					errs = append(errs, fmt.Sprintf("%s is invalid", fieldKey))
				}
			}
			return errors.New(strings.Join(errs, ", "))
		}
		return err
	}
	return nil
}
