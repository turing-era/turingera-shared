package token

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"

	"github.com/turing-era/turingera-shared/log"
)

// JwtTokenVerifier JWT验证器
type JwtTokenVerifier struct {
	appid      string
	issuer     string
	alg        string
	publicKey *rsa.PublicKey
}

func NewJwtTokenVerifier(publicKeyFile string) *JwtTokenVerifier {
	pubKey, err := loadPublicKey(publicKeyFile)
	if err != nil {
		panic("loadPublicKey err: " + err.Error())
	}
	return &JwtTokenVerifier{
		appid:      viper.GetString("auth.jwt_appid"),
		issuer:     viper.GetString("auth.jwt_issuer"),
		alg:        viper.GetString("auth.jwt_alg"),
		publicKey: pubKey,
	}
}

// 加载公钥
func loadPublicKey(publicKeyFile string) (*rsa.PublicKey, error) {
	f, err := os.Open(publicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("cannot open public key file: %v", err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("cannot read public key: %v", err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(b)
	if err != nil {
		return nil, fmt.Errorf("cannot parse public key: %v", err)
	}
	return pubKey, nil
}

// PrivyClaims Defining a Go type for Privy JWTs
type PrivyClaims struct {
	AppId      string `json:"aud,omitempty"`
	Expiration uint64 `json:"exp,omitempty"`
	Issuer     string `json:"iss,omitempty"`
	UserId     string `json:"sub,omitempty"`
}

// This method will be used to check the token's claims later
func (v *JwtTokenVerifier) valid(c *jwt.RegisteredClaims) error {
	if len(c.Audience) == 0 || c.Audience[0] != v.appid {
		return errors.New("aud claim must be your Privy App ID")
	}
	if c.Issuer != c.Issuer {
		return errors.New("iss claim must be 'privy.io'")
	}
	if c.ExpiresAt.Time.Before(time.Now()) {
		return errors.New("token is expired")
	}
	return nil
}

// This method will be used to load the verification key in the required format later
func (v *JwtTokenVerifier) keyFunc(token *jwt.Token) (interface{}, error) {
	if token.Method.Alg() != viper.GetString("auth.jwt_alg") {
		return nil, fmt.Errorf("unexpected JWT signing method=%v", token.Header["alg"])
	}
	return v.publicKey, nil
}

// Verify JWT验证
func (v *JwtTokenVerifier) Verify(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &jwt.RegisteredClaims{}, v.keyFunc)
	if err != nil {
		return "", fmt.Errorf("cannot parse token: %v", err)
	}
	if !token.Valid {
		return "", fmt.Errorf("token not valid")
	}
	clm, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", fmt.Errorf("token claim is not RegisteredClaims")
	}
	log.Debugf("appid: %v, clm: %+v", v.appid, clm)
	if err = v.valid(clm); err != nil {
		return "", fmt.Errorf("claim not valid: %v", err)
	}
	return clm.Subject, nil
}
