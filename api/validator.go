package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/santiagot714/SimpleBank/util"
)

// validCurrency is a custom validator for currency.
// It checks if the currency is supported.
var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	currency, ok := fieldLevel.Field().Interface().(string)
	if !ok {
		return false
	}
	return util.IsSupportedCurrency(currency)
}
