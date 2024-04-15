package v1

import (
	"net/http"
	"single-window/internal/entity"
	"single-window/internal/usecase"
	"single-window/pkg/httpserver"

	"github.com/gin-gonic/gin"
)

type authRoutes struct {
	auc              usecase.IAuthUseCase
	authCookieName   string
	authCookiePath   string
	authCookieDomain string
}

func NewAuthRoutes(auc usecase.IAuthUseCase, name, path, domain string) *authRoutes {
	return &authRoutes{
		auc:              auc,
		authCookieName:   name,
		authCookiePath:   path,
		authCookieDomain: domain,
	}
}

func (a *authRoutes) AuthUser(ctx *gin.Context, h AuthUserParams) {
	body := entity.AuthUserJSONBody{}

	if err := ctx.Bind(&body); err != nil {
		httpserver.ErrorResponse(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}

	token, err := a.auc.AuthHandle(body, entity.AuthUserParams(h))
	if err != nil {
		httpserver.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie(
		a.authCookieName,
		token.Token,
		3600*24, // TODO: получать из конфига (сейчас неправильный парсинг)
		a.authCookiePath,
		a.authCookieDomain,
		true,
		true,
	)
	ctx.JSON(http.StatusOK, httpserver.ServerResponse[entity.Claims]{Data: token.Claims})
}
