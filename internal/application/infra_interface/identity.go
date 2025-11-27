package infra_interface

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	UserId string   `json:"user_id"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

type JWTManager interface {
	Sign(isRefresh bool, userID string, roles ...string) (string, error)
	Verify(rawToken string) (*JWTClaims, error)
}
