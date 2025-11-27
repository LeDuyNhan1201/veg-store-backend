package service

import (
	"veg-store-backend/internal/infrastructure/core"
)

type Service[TDatabase any, TRepository any] struct {
	DB         TDatabase
	Repository TRepository
	*core.Core
}
