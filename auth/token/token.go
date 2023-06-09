package token

import (
	"crypto/rsa"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/turing-era/turingera-shared/log"
)

// JWTTokenVerifier JWT验证器
type JWTTokenVerifier struct {
	PublicKey *rsa.PublicKey
}

// Verify JWT验证
func (v *JWTTokenVerifier) Verify(token string) (string, error) {
	t, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return v.PublicKey, nil
	})
	log.Debugf("t: %+v", t)
	if err != nil {
		return "", fmt.Errorf("cannot parse token: %v", err)
	}
	if !t.Valid {
		return "", fmt.Errorf("token not valid")
	}
	clm, ok := t.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", fmt.Errorf("token claim is not StandardClaims")
	}
	if err := clm.Valid(); err != nil {
		return "", fmt.Errorf("claim not valid: %v", err)
	}

	return clm.Subject, nil
}
