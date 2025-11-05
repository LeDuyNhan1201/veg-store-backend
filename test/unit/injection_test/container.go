package injection_test

import (
	"veg-store-backend/injection"
)

func InjectMock() *injection.Container {
	return injection.Inject("test")
}
