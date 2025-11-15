package identity

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"
	"veg-store-backend/internal/application/infra_interface"
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/util"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type jwtManager struct {
	*core.Core
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewJWTManager(core *core.Core) (infra_interface.JWTManager, error) {
	// Set config path to .../.../keypair
	keypairPath := util.GetConfigPathFromGoMod("secrets/keypair")
	privateKeyPath := fmt.Sprintf("%s/%s", keypairPath, core.Config.JWT.PrivateKey)
	publicKeyPath := fmt.Sprintf("%s/%s", keypairPath, core.Config.JWT.PublicKey)
	core.Logger.Info(fmt.Sprintf("Private key path: %s", privateKeyPath))
	core.Logger.Info(fmt.Sprintf("Public key path: %s", publicKeyPath))

	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		core.Logger.Fatal("Error to read private key", zap.Error(err))
	}

	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		core.Logger.Fatal("Error to read public key", zap.Error(err))
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		core.Logger.Fatal("Error to parse private key", zap.Error(err))
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		core.Logger.Fatal("Error to parse public key", zap.Error(err))
	}

	return &jwtManager{
		Core:       core,
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func (manager *jwtManager) Sign(isRefresh bool, userId string, roles ...string) (string, error) {
	var err error
	Expiration, err := util.ParseDuration(manager.Config.JWT.AccessDuration)
	if err != nil {
		manager.Logger.Error("Error to parse string to duration", zap.Error(err))
		return "", manager.Error.Auth.Unauthenticated
	}

	if isRefresh {
		Expiration, err = util.ParseDuration(manager.Config.JWT.RefreshDuration)
		if err != nil {
			manager.Logger.Error("Error to parse string to duration", zap.Error(err))
			return "", manager.Error.Auth.Unauthenticated
		}
	}

	if roles == nil || len(roles) == 0 {
		roles = []string{}
	}

	claims := &infra_interface.JWTClaims{
		UserId: userId,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    manager.Config.JWT.ExpectedIssuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(Expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(manager.privateKey)
}

func (manager *jwtManager) Verify(tokenStr string) (*infra_interface.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &infra_interface.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return manager.publicKey, nil
	})
	if err != nil {
		manager.Logger.Error("Invalid token", zap.Error(err))
		return nil, manager.Error.Invalid.Token
	}

	claims, ok := token.Claims.(*infra_interface.JWTClaims)
	if !ok || !token.Valid {
		manager.Logger.Error("Invalid claims", zap.Error(err))
		return nil, manager.Error.Auth.Unauthenticated
	}
	return claims, nil
}

func (manager *jwtManager) toJWTSigningMethod(jwtAlgorithm string) *jwt.SigningMethodRSA {
	switch jwtAlgorithm {
	case "RS256":
		return jwt.SigningMethodRS256
	case "RS384":
		return jwt.SigningMethodRS384
	case "RS512":
		return jwt.SigningMethodRS512
	default:
		manager.Logger.Fatal("Unsupported signing method", zap.String("algorithm", jwtAlgorithm))
	}
	return nil
}
