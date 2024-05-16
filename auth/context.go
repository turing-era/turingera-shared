package auth

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/turing-era/turingera-shared/log"
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
	fmt.Printf("md: %+v\n", md)
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
	log.Debugf("user_id: %v", userID)
	return context.WithValue(ctx, userIDKey, userID)
}

// InheritCtx 继承ctx
func InheritCtx(ctx context.Context, server string) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(map[string]string{})
	}
	md.Set("from", "internal")
	md.Set("server", server)
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

// UidFromCtx 从ctx中取出uid
func UidFromCtx(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok {
		return 0, status.Error(codes.Unauthenticated, "uid not found")
	}
	uid, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return 0, status.Errorf(codes.Unauthenticated, "userID[%v] ParseInt err: %v", userID, err)
	}
	return uid, nil
}
