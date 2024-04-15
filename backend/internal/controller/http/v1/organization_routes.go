package v1

import (
	"net/http"
	"single-window/internal/entity"
	"single-window/internal/usecase"
	"single-window/pkg/httpserver"
	"single-window/pkg/logger"
	"single-window/pkg/swcontext"

	"github.com/gin-gonic/gin"
)

type organizationRoutes struct {
	ouc          usecase.IOrganizationUseCase
	l            logger.ILogger
	swCtxAdapter swcontext.SwGinContext
}

func newOrganizationRoutes(ouc usecase.IOrganizationUseCase, l logger.ILogger, swCtxAdapter swcontext.SwGinContext) *organizationRoutes {
	return &organizationRoutes{
		ouc:          ouc,
		l:            l,
		swCtxAdapter: swCtxAdapter,
	}
}

func (r *organizationRoutes) GetOrganizationsV1(c *gin.Context, params GetOrganizationsV1Params) {
	claims, err := r.swCtxAdapter.GetTokenClaims(c)
	if err != nil {
		r.l.Error(err, "http - v1 - organizations")
		httpserver.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	organizations, err := r.ouc.GetOrganizationsV1(
		entity.GetOrganizationsData{
			Claims:           *claims,
			OrganizationCode: params.OrganizationCode,
			SearchToken:      params.SearchToken,
		},
		entity.Pagination{
			Limit: *params.Limit, Offset: *params.Offset,
		},
	)
	if err != nil {
		r.l.Error(err, "http - v1 - organizations")
		httpserver.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, httpserver.ServerResponse[[]entity.Organization]{
		Data:        organizations.Data,
		HasMoreData: &organizations.HasMoreData,
	})
}
