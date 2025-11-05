package core

import (
	"veg-store-backend/internal/application/exception"

	"go.uber.org/zap"
)

var Error *exception.AppError
var Configs = &Config{}
var Logger *zap.Logger
var Translator *Localizer
