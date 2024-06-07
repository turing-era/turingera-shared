package token

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
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
	privateKey *rsa.PrivateKey
}

// NewJwtTokenGenerator 新建token生成器
// jwt测试：https://jwt.io/
// 公私钥生成：https://www.metools.info/code/c80.html
func NewJwtTokenGenerator() *JwtTokenGenerator {
	path := viper.GetString("auth.private_path")
	privKey, err := loadPrivateKey(path)
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
func loadPrivateKey(keyPath string) (*rsa.PrivateKey, error) {
	pkFile, err := os.Open(keyPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open private key: %v", err)
	}
	pkBytes, err := ioutil.ReadAll(pkFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read private key: %v", err)
	}
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(pkBytes)
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
