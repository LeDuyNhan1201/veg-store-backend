package identity_test

import (
	"fmt"
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/application/infra_interface"

	"github.com/stretchr/testify/mock"
)

type MockJWTManager struct {
	mock.Mock
}

func (mockManager *MockJWTManager) Sign(isRefresh bool, userID string, roles ...string) (string, error) {
	args := mockManager.Called(isRefresh, userID, roles)
	return args.String(0), args.Error(1)
}

func (mockManager *MockJWTManager) Verify(rawToken string) (*infra_interface.JWTClaims, error) {
	args := mockManager.Called(rawToken)
	claims, _ := args.Get(0).(*infra_interface.JWTClaims)
	return claims, args.Error(1)
}

/*----------------------------------INJECTION--------------------------------------*/

func (mockManager *MockJWTManager) Name() string { return "MockJWTManager" }
func (mockManager *MockJWTManager) Start() error {
	core.Logger.Debug(fmt.Sprintf("%s initialized", mockManager.Name()))
	return nil
}
func (mockManager *MockJWTManager) Stop() error {
	core.Logger.Debug(fmt.Sprintf("%s destroyed", mockManager.Name()))
	return nil
}
