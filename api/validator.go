package api

import (
	"github.com/go-playground/validator/v10"
	"simplebank/util"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		// Check currency is supported or not
		return util.IsSupportedCurrency(currency)
	}
	return false
}
