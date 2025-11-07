package identity

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	JWTManagerModule,
)
