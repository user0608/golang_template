package security

import (
	"crypto/rsa"
	"errors"
	"goqrs/envs"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ksaucedo002/answer/errores"
)

var (
	once       sync.Once
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)
var (
	jwtTokenLife time.Duration = 730 * time.Hour
)

func LoadRSAKeys() (err error) {
	jwtTokenLife, err = time.ParseDuration(envs.FindEnv("TOKEN_LIFE", "24h"))
	if err != nil {
		return err
	}
	keyPath := envs.FindEnv("GOQRS_RSA_PRIVATE", "certificates/id_rsa")
	pubPath := envs.FindEnv("GOQRS_RSA_PUBLIC", "certificates/id_rsa.pub")
	once.Do(func() {
		var private []byte
		var public []byte
		private, err = os.ReadFile(keyPath)
		if err != nil {
			return
		}
		public, err = os.ReadFile(pubPath)
		if err != nil {
			return
		}
		privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(private)
		if err != nil {
			return
		}
		publicKey, err = jwt.ParseRSAPublicKeyFromPEM(public)
		if err != nil {
			return
		}
	})
	return err
}

type JWTValues struct {
	Username string `json:"username"`
}
type jwtCustomClaims struct {
	JWTValues
	jwt.RegisteredClaims
}

func GenToken(values JWTValues) (string, error) {
	now := time.Now()
	customClaims := jwtCustomClaims{
		JWTValues: values,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ksaucedo",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(jwtTokenLife)),
			Audience:  jwt.ClaimStrings{"qrsystems"},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, customClaims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", nil
	}
	return tokenString, nil

}
func ValidateToken(tokenString string) (values JWTValues, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtCustomClaims{}, verifyWithKey)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return JWTValues{}, errores.NewUnauthorizedf(nil, "[close] la sesión a expirado")
		}
		return JWTValues{}, errores.NewUnauthorizedf(nil, "[close] no se reconoció la sesión")
	}
	if !token.Valid {
		return JWTValues{}, errores.NewUnauthorizedf(nil, "[close] token invalido")
	}
	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok {
		return JWTValues{}, errores.NewInternalf(nil, "[close] claims invalidos")
	}
	return claims.JWTValues, nil
}
func verifyWithKey(token *jwt.Token) (any, error) {
	return publicKey, nil
}
