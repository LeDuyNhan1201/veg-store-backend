package validation

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var RangeValidator validator.Func = func(fl validator.FieldLevel) bool {
	// Pass nil values
	value := fl.Field().String()
	if value == "" {
		return true
	}

	param := fl.Param() // ví dụ: "1-100"
	if param == "" {
		return true
	}

	// parse minParam và maxParam từ tag
	minParam, maxParam := ParseRangeParam(param)

	field := fl.Field()
	switch field.Kind() {
	case reflect.String:
		length := float64(len(field.String()))
		return length >= minParam && length <= maxParam

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := float64(field.Int())
		return val >= minParam && val <= maxParam

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val := float64(field.Uint())
		return val >= minParam && val <= maxParam

	case reflect.Float32, reflect.Float64:
		val := field.Float()
		return val >= minParam && val <= maxParam

	default:
		return false
	}
}

func ParseRangeParam(param string) (float64, float64) {
	// parse minParam và maxParam từ tag
	params := strings.Split(param, "-")
	if len(params) != 2 {
		zap.L().Fatal("Cannot parse range-param")
	}

	minParam, err1 := strconv.ParseFloat(params[0], 64)
	maxParam, err2 := strconv.ParseFloat(params[1], 64)
	if err1 != nil || err2 != nil {
		zap.L().Fatal("Cannot parse range-param", zap.Error(err1), zap.Error(err2))
	}

	return minParam, maxParam
}
