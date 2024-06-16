package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/turing-era/turingera-shared/auth/token"
	"github.com/turing-era/turingera-shared/cutils"
)

type interceptor struct {
	verifier *token.JwtTokenVerifier
}

// PrivyPrefix privy SDK userID前缀
const PrivyPrefix = "did:privy:"

var DefaultInterceptor *interceptor

// LoadInterceptor 注入方式拦截器方法
func LoadInterceptor(keyPath string) (grpc.UnaryServerInterceptor, error) {
	intercept := &interceptor{
		verifier: token.NewJwtTokenVerifier(keyPath),
	}
	return intercept.handleReq, nil

}

// InitInterceptor 非注入方式初始化拦截器
func InitInterceptor(keyPath string) error {
	intercept := &interceptor{
		verifier: token.NewJwtTokenVerifier(keyPath),
	}
	DefaultInterceptor = intercept
	return nil
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
	userID = strings.ReplaceAll(userID, PrivyPrefix, "")
	return userID, true
}

func (i *interceptor) handleReq(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var userID string
	// 内部调用
	from := metadata.ValueFromIncomingContext(ctx, "from")
	internal := len(from) > 0 && from[0] == "internal"
	// 开放方法
	openMethod := cutils.InStringList(viper.GetStringSlice("auth.open_method"), info.FullMethod)
	if !internal && !openMethod {
		tkn, err := tokenFromCtx(ctx)
		if err != nil {
			return nil, err
		}
		userID, err = i.verifier.Verify(tkn)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "token not valid: %v", err)
		}
		userID = strings.ReplaceAll(userID, PrivyPrefix, "")
		ctx = ctxWithUserID(ctx, userID)
	}
	return handler(ctx, req)
}
