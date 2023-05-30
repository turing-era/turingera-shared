package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authorizationHeader = "authorization"
	bearerPrefix        = "Bearer "
	userIDKey           = "user_id"
)

var unauthenticated = status.Error(codes.Unauthenticated, "unauthenticated")

func tokenFromCtx(c context.Context) (string, error) {
	m, ok := metadata.FromIncomingContext(c)
	if !ok {
		return "", unauthenticated
	}
	tkn := ""
	for _, v := range m[authorizationHeader] {
		if strings.HasPrefix(v, bearerPrefix) {
			tkn = v[len(bearerPrefix):]
		}
	}
	if tkn == "" {
		return "", unauthenticated
	}
	return tkn, nil
}

func ctxWithUserID(c context.Context, aid string) context.Context {
	return context.WithValue(c, userIDKey, aid)
}

// UserIDFromCtx 从ctx中取出账号id
func UserIDFromCtx(c context.Context) (string, error) {
	v := c.Value(userIDKey)
	aid, ok := v.(string)
	if !ok {
		return "", unauthenticated
	}
	return aid, nil
}
