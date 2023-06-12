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

func tokenFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "token not found")
	}
	// log.Debugf("md: %+v", md)
	var tkn string
	for _, v := range md[authorizationHeader] {
		if strings.HasPrefix(v, bearerPrefix) {
			tkn = v[len(bearerPrefix):]
		}
	}
	if len(tkn) == 0 {
		return "", status.Error(codes.Unauthenticated, "token empty")
	}
	return tkn, nil
}

func ctxWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// InheritCtx 继承ctx
func InheritCtx(ctx context.Context) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	return metadata.NewOutgoingContext(ctx, md)
}

// UserIDFromCtx 从ctx中取出账号id
func UserIDFromCtx(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "userID not found")
	}
	return userID, nil
}
