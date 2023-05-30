package auth

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/funstartech/funstar-shared/auth/token"
)

// Interceptor 拦截器方法
func Interceptor(publicKeyFile string) (grpc.UnaryServerInterceptor, error) {
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
	i := &interceptor{
		verifier: &token.JWTTokenVerifier{PublicKey: pubKey},
	}
	return i.handleReq, nil

}

type interceptor struct {
	verifier *token.JWTTokenVerifier
}

func (i *interceptor) handleReq(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var userID string
	if info.FullMethod != "/turingera.server.userinfo.Userinfo/InitUser" {
		tkn, err := tokenFromCtx(ctx)
		if err != nil {
			return nil, err
		}
		userID, err = i.verifier.Verify(tkn)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "token not valid: %v", err)
		}
	}
	return handler(ctxWithUserID(ctx, userID), req)
}
