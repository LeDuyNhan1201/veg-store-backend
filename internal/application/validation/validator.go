package validation

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func Init() {
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := validate.RegisterValidation("range", RangeValidator)
		if err != nil {
			zap.L().Fatal("validation validator register error", zap.Error(err))
		}
		err = validate.RegisterValidation("fields", FieldsValidator)
		if err != nil {
			zap.L().Fatal("fields validator register error", zap.Error(err))
		}
	}
}
