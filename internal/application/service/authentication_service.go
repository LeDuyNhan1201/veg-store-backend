package service

import (
	"veg-store-backend/internal/application/dto"
	"veg-store-backend/internal/application/infra_interface"
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/internal/infrastructure/data"
)

type AuthenticationService interface {
	Tokens(request dto.SignInRequest) (*dto.Tokens, error)
	Me(id string) (string, error)
}

type authenticationService struct {
	Service[*data.PostgresDB, any]
	userService UserService
	jwtManager  infra_interface.JWTManager
}

func NewAuthenticationService(
	core *core.Core,
	db *data.PostgresDB,
	userService UserService,
	jwtManager infra_interface.JWTManager,
) AuthenticationService {
	return &authenticationService{
		Service: Service[*data.PostgresDB, any]{
			Core: core,
			DB:   db,
		},
		userService: userService,
		jwtManager:  jwtManager,
	}
}

func (service *authenticationService) Tokens(request dto.SignInRequest) (*dto.Tokens, error) {
	var err error
	user, err := service.userService.FindByUsername(request.Username)
	if err != nil {
		return nil, service.Error.Invalid.Username
	}

	accessToken, err := service.jwtManager.Sign(false, user.ID.String())
	if err != nil {
		return nil, service.Error.Auth.Unauthenticated
	}
	refreshToken, err := service.jwtManager.Sign(true, user.ID.String())
	if err != nil {
		return nil, service.Error.Auth.Unauthenticated
	}

	return &dto.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (service *authenticationService) Me(id string) (string, error) {
	user, err := service.userService.FindById(id)
	if err != nil {
		return "", err
	}
	return user.Name, nil
}
