package token

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// JwtTokenGenerator jwt token生成器
type JwtTokenGenerator struct {
	appid      string
	issuer     string
	alg        string
	privateKey interface{}
}

// NewJwtTokenGenerator 新建token生成器
// jwt测试：https://jwt.io/
// 公私钥生成：https://www.metools.info/code/c80.html
func NewJwtTokenGenerator() *JwtTokenGenerator {
	alg := viper.GetString("auth.jwt_alg")
	path := viper.GetString("auth.private_path")
	privKey, err := loadPrivateKey(path, alg)
	if err != nil {
		panic("loadPrivateKey err: " + err.Error())
	}
	return &JwtTokenGenerator{
		appid:      viper.GetString("auth.jwt_appid"),
		issuer:     viper.GetString("auth.jwt_issuer"),
		alg:        viper.GetString("auth.jwt_alg"),
		privateKey: privKey,
	}
}

// 加载私钥
func loadPrivateKey(keyPath, alg string) (interface{}, error) {
	pkFile, err := os.Open(keyPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open private key: %v", err)
	}
	pkBytes, err := io.ReadAll(pkFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read private key: %v", err)
	}
	var privKey interface{}
	switch alg {
	case "RS512":
		privKey, err = jwt.ParseRSAPrivateKeyFromPEM(pkBytes)
	case "ES256":
		privKey, err = jwt.ParseECPrivateKeyFromPEM(pkBytes)
	default:
		return nil, fmt.Errorf("invalid alg: %v", alg)
	}
	if err != nil {
		return nil, fmt.Errorf("cannot parse private key: %v", err)
	}
	return privKey, nil
}

func (t *JwtTokenGenerator) GenerateToken(userID string, expire int64) (string, error) {
	now := time.Now()
	var signMethod jwt.SigningMethod
	switch t.alg {
	case "RS512":
		signMethod = jwt.SigningMethodRS512
	case "ES256":
		signMethod = jwt.SigningMethodES256
	default:
		return "", fmt.Errorf("invalid alg: %v", t.alg)
	}
	token := jwt.NewWithClaims(signMethod, jwt.RegisteredClaims{
		Issuer:    t.issuer,
		IssuedAt:  &jwt.NumericDate{Time: now},
		ExpiresAt: &jwt.NumericDate{Time: now.Add(time.Second * time.Duration(expire))},
		Subject:   userID,
		Audience: []string{t.appid},
	})
	return token.SignedString(t.privateKey)
}
