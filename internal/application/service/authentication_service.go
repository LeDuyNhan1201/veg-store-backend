package service

import (
	"veg-store-backend/internal/application/dto"
	"veg-store-backend/internal/application/infra_interface"
	"veg-store-backend/internal/infrastructure/core"
)

type AuthenticationService interface {
	Tokens(request dto.SignInRequest) (*dto.Tokens, error)
	Me(id string) (string, error)
}

type authenticationService struct {
	*core.Core
	userService UserService
	jwtManager  infra_interface.JWTManager
}

func NewAuthenticationService(core *core.Core, userService UserService, jwtManager infra_interface.JWTManager) AuthenticationService {
	return &authenticationService{
		Core:        core,
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

	accessToken, err := service.jwtManager.Sign(false, user.Id.String())
	if err != nil {
		return nil, service.Error.Auth.Unauthenticated
	}
	refreshToken, err := service.jwtManager.Sign(true, user.Id.String())
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
