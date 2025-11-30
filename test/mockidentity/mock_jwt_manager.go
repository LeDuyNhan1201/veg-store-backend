package mockidentity

import (
	"veg-store-backend/internal/application/iface"

	"github.com/stretchr/testify/mock"
)

type MockJWTManager struct {
	mock.Mock
}

func (mck *MockJWTManager) Sign(isRefresh bool, userID string, roles ...string) (string, error) {
	args := mck.Called(isRefresh, userID, roles)
	return args.String(0), args.Error(1)
}

func (mck *MockJWTManager) Verify(rawToken string) (*iface.JWTClaims, error) {
	args := mck.Called(rawToken)
	claims, _ := args.Get(0).(*iface.JWTClaims)
	return claims, args.Error(1)
}
