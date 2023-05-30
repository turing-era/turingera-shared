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

func tokenFromCtx(c context.Context) (string, error) {
	values := metadata.ValueFromIncomingContext(c, authorizationHeader)
	var tkn string
	for _, v := range values {
		if strings.HasPrefix(v, bearerPrefix) {
			tkn = v[len(bearerPrefix):]
		}
	}
	if len(tkn) == 0 {
		return "", status.Error(codes.Unauthenticated, "token not found")
	}
	return tkn, nil
}

func ctxWithUserID(c context.Context, userID string) context.Context {
	return context.WithValue(c, userIDKey, userID)
}

// UserIDFromCtx 从ctx中取出账号id
func UserIDFromCtx(c context.Context) (string, error) {
	userID, ok := c.Value(userIDKey).(string)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "userID not found")
	}
	return userID, nil
}
