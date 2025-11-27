package identity_mock

import (
	"veg-store-backend/internal/application/infra_interface"

	"github.com/stretchr/testify/mock"
)

type MockJWTManager struct {
	mock.Mock
}

func (mck *MockJWTManager) Sign(isRefresh bool, userID string, roles ...string) (string, error) {
	args := mck.Called(isRefresh, userID, roles)
	return args.String(0), args.Error(1)
}

func (mck *MockJWTManager) Verify(rawToken string) (*infra_interface.JWTClaims, error) {
	args := mck.Called(rawToken)
	claims, _ := args.Get(0).(*infra_interface.JWTClaims)
	return claims, args.Error(1)
}
