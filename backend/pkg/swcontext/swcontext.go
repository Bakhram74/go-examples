package swcontext

import (
	"context"
	"errors"
	"fmt"
	"single-window/internal/entity"

	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"
)

const (
	tokenClaimsGinCtxKey = "token_claims"
)

type SwGinContext ContextClaimsAdapter[*gin.Context, entity.Claims]

type ContextClaimsAdapter[Ctx context.Context, Claims any] interface {
	GetContext(ctx context.Context) (Ctx, error)
	GetCookie(ctx Ctx, cookieName string) (*string, error)
	GetTokenClaims(ctx Ctx) (*Claims, error)
	SetTokenClaims(ctx Ctx, claims *Claims) error
}

type ginContextAdapter[T any] struct{}

func NewGinContextCookieAdapter[T any]() ContextClaimsAdapter[*gin.Context, T] {
	return &ginContextAdapter[T]{}
}

func (ca *ginContextAdapter[T]) GetContext(ctx context.Context) (*gin.Context, error) {
	return middleware.GetGinContext(ctx), nil
}

func (ca *ginContextAdapter[T]) GetCookie(ctx *gin.Context, cookieName string) (*string, error) {
	cookie, err := ctx.Cookie(cookieName)
	if err != nil {
		return nil, fmt.Errorf("error in get cookie: %w", err)
	}
	return &cookie, nil
}

func (ca *ginContextAdapter[T]) SetTokenClaims(ctx *gin.Context, claims *T) error {
	ctx.Set(tokenClaimsGinCtxKey, claims)
	return nil
}

func (ca *ginContextAdapter[T]) GetTokenClaims(ctx *gin.Context) (*T, error) {
	ctxClaims := ctx.Value(tokenClaimsGinCtxKey)

	claims, ok := ctxClaims.(*T)
	if !ok {
		return nil, errors.New("claims type is not correct")
	}

	return claims, nil
}
