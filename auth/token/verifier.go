package token

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// JwtTokenVerifier JWT验证器
type JwtTokenVerifier struct {
	appid     string
	issuer    string
	alg       string
	publicKey interface{}
}

func NewJwtTokenVerifier(keyPath string) *JwtTokenVerifier {
	alg := viper.GetString("auth.jwt_alg")
	pubKey, err := loadPublicKey(keyPath, alg)
	if err != nil {
		panic("loadPublicKey err: " + err.Error())
	}
	return &JwtTokenVerifier{
		appid:     viper.GetString("auth.jwt_appid"),
		issuer:    viper.GetString("auth.jwt_issuer"),
		alg:       alg,
		publicKey: pubKey,
	}
}

// 加载公钥
func loadPublicKey(keyPath, alg string) (interface{}, error) {
	pkFile, err := os.Open(keyPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open public key file: %v", err)
	}
	pkBytes, err := io.ReadAll(pkFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read public key: %v", err)
	}
	var pubKey interface{}
	switch alg {
	case "RS512":
		pubKey, err = jwt.ParseRSAPublicKeyFromPEM(pkBytes)
	case "ES256":
		pubKey, err = jwt.ParseECPublicKeyFromPEM(pkBytes)
	default:
		return nil, fmt.Errorf("invalid alg: %v", alg)
	}
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
	if len(c.Audience) > 0 && c.Audience[0] != v.appid {
		return fmt.Errorf("aud claim must be your Privy AppID -- %+v", c)
	}
	if c.Issuer != c.Issuer {
		return fmt.Errorf("iss claim must be 'privy.io' -- %+v", c)
	}
	if c.ExpiresAt.Time.Before(time.Now()) {
		return fmt.Errorf("token is expired -- %+v", c)
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
	// fmt.Printf("accessToken: %s\n", accessToken)
	token, err := jwt.ParseWithClaims(accessToken, &jwt.RegisteredClaims{}, v.keyFunc)
	if err != nil {
		return "", fmt.Errorf("cannot parse token[%v]: %v", accessToken, err)
	}
	if !token.Valid {
		return "", fmt.Errorf("token not valid")
	}
	clm, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", fmt.Errorf("token claim is not RegisteredClaims")
	}
	// fmt.Printf("appid: %v, clm: %+v\n", v.appid, clm)
	if err = v.valid(clm); err != nil {
		return "", fmt.Errorf("claim not valid: %v", err)
	}
	return clm.Subject, nil
}
