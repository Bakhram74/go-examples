package middlewares

import (
	"context"
	"fmt"
	"single-window/internal/entity"
	"single-window/pkg/jwt"
	"single-window/pkg/swcontext"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
)

type Authenticator struct {
	processor jwt.IJWTParser
	adapter   swcontext.ContextClaimsAdapter[*gin.Context, entity.Claims]
}

func NewAuthenticator(
	processor jwt.IJWTParser,
	adapter swcontext.ContextClaimsAdapter[*gin.Context, entity.Claims],
) *Authenticator {
	return &Authenticator{
		processor: processor,
		adapter:   adapter,
	}
}

func NewCookieAuthContextAdapter() swcontext.ContextClaimsAdapter[*gin.Context, entity.Claims] {
	return swcontext.NewGinContextCookieAdapter[entity.Claims]()
}

func (a *Authenticator) OAPIAuthenticate(ctx context.Context, ai *openapi3filter.AuthenticationInput) error {
	isSecuritySchemeSwCookieAuth := ai.SecuritySchemeName == "SWCookieAuth"
	isApiKey := ai.SecurityScheme.Type == "apiKey"
	isCookie := ai.SecurityScheme.In == "cookie"

	if isApiKey && isCookie && isSecuritySchemeSwCookieAuth {
		if err := a.getTokenClaimsFromGin(ctx, ai.SecurityScheme.Name); err != nil {
			return fmt.Errorf("cannot get token: %w", err)
		}
	}
	return nil
}

func (a *Authenticator) getTokenClaimsFromGin(ctx context.Context, cookieName string) error {
	c, err := a.adapter.GetContext(ctx)
	if err != nil {
		return fmt.Errorf("cannot get context: %w", err)
	}

	cookie, err := a.adapter.GetCookie(c, cookieName)
	if err != nil {
		return fmt.Errorf("error in get cookie: %w", err)
	}

	claims := entity.Claims{}
	if err := a.processor.ParseToken(*cookie, &claims); err != nil {
		return fmt.Errorf("error in token parser: %w", err)
	}
	if err := a.adapter.SetTokenClaims(c, &claims); err != nil {
		return fmt.Errorf("error in set claims: %w", err)
	}
	return nil
}
