package auth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"

	"github.com/turing-era/turingera-shared/cutils"
	"github.com/turing-era/turingera-shared/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/funstartech/funstar-shared/auth/token"
)

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

// Interceptor 拦截器方法
func Interceptor(publicKeyFile string) (grpc.UnaryServerInterceptor, error) {
	pubKey, err := loadPublicKey(publicKeyFile)
	if err != nil {
		return nil, err
	}
	log.Debugf("load publicKeyFile: %v, pubKey: %v", publicKeyFile, pubKey)
	i := &interceptor{
		verifier: &token.JWTTokenVerifier{PublicKey: pubKey},
	}
	return i.handleReq, nil

}

var DefaultInterceptor *interceptor

// InitInterceptor 原生拦截器
func InitInterceptor(publicKeyFile string) error {
	pubKey, err := loadPublicKey(publicKeyFile)
	if err != nil {
		return err
	}
	DefaultInterceptor = &interceptor{
		verifier: &token.JWTTokenVerifier{PublicKey: pubKey},
	}
	return nil
}

type interceptor struct {
	verifier *token.JWTTokenVerifier
}

// GetAuthUserID 获取鉴权userID
func (i *interceptor) GetAuthUserID(w http.ResponseWriter, r *http.Request) (string, bool) {
	var tkn string
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, bearerPrefix) {
		tkn = auth[len(bearerPrefix):]
	}
	if len(tkn) == 0 {
		http.Error(w, fmt.Sprintf("token empty"), http.StatusUnauthorized)
		return "", false
	}
	userID, err := i.verifier.Verify(tkn)
	if err != nil {
		http.Error(w, fmt.Sprintf("token not valid: %v", err), http.StatusUnauthorized)
		return "", false
	}
	return userID, true
}

func (i *interceptor) handleReq(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var userID string
	// 内部调用
	from := metadata.ValueFromIncomingContext(ctx, "from")
	internal := len(from) > 0 && from[0] == "internal"
	// 开放方法
	openMethods := viper.GetStringSlice("auth.open_method")
	openMethod := cutils.InStringList(openMethods, info.FullMethod)
	if !internal && !openMethod {
		tkn, err := tokenFromCtx(ctx)
		if err != nil {
			return nil, err
		}
		userID, err = i.verifier.Verify(tkn)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "token not valid: %v", err)
		}
		ctx = ctxWithUserID(ctx, userID)
	}
	return handler(ctx, req)
}
