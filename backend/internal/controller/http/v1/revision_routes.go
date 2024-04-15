package v1

import (
	"net/http"
	"single-window/internal/entity"
	"single-window/internal/usecase"
	"single-window/pkg/httpserver"
	"single-window/pkg/logger"
	"single-window/pkg/swcontext"

	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type revisionRoutes struct {
	ruc          usecase.IRevisionUseCase
	l            logger.ILogger
	swCtxAdapter swcontext.SwGinContext
}

func newRevisionRoutes(ruc usecase.IRevisionUseCase, l logger.ILogger, swCtxAdapter swcontext.SwGinContext) *revisionRoutes {
	return &revisionRoutes{
		ruc:          ruc,
		l:            l,
		swCtxAdapter: swCtxAdapter,
	}
}

func (r *revisionRoutes) GetRevisionsV1(c *gin.Context, params GetRevisionsV1Params) {
	claims, err := r.swCtxAdapter.GetTokenClaims(c)
	if err != nil {
		r.l.Error(err, "http - v1 - revisions")
		httpserver.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	revisions, err := r.ruc.GetRevisionsV1(
		entity.GetRevisionsData{
			Claims:    *claims,
			Status:    entity.EntityStatus(*params.Status),
			DisputeID: params.DisputeId,
		},
		entity.Pagination{
			Limit: *params.Limit, Offset: *params.Offset,
		},
	)
	if err != nil {
		r.l.Error(err, "http - v1 - revisions")
		httpserver.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, httpserver.ServerResponse[[]entity.Revision]{
		Data:        revisions.Data,
		HasMoreData: &revisions.HasMoreData,
	})
}

func (r *revisionRoutes) CreateRevisionV1(c *gin.Context) {
	claims, err := r.swCtxAdapter.GetTokenClaims(c)
	if err != nil {
		r.l.Error(err, "http - v1 - revision - create")
		httpserver.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		r.l.Error(err, "http - v1 - revision - create")
		httpserver.ErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	var body entity.CreateRevisionV1MultipartBody
	body.DisputeId = form.Value["dispute_id"][0]
	body.OrganizationId = form.Value["organization_id"][0]
	body.MessageBody = form.Value["message_body"][0]
	if len(form.File["file"]) == 1 {
		file := &openapi_types.File{}
		file.InitFromMultipart(form.File["file"][0])
		body.File = file
	} else {
		body.File = nil
	}

	revision, err := r.ruc.CreateRevisionV1(claims, &body)
	if err != nil {
		r.l.Error(err, "http - v1 - revision - create")
		httpserver.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, httpserver.ServerResponse[*entity.Revision]{Data: revision})
}

func (r *revisionRoutes) GetCorrespondencesV1(c *gin.Context, param GetCorrespondencesV1Params) {
	claims, err := r.swCtxAdapter.GetTokenClaims(c)
	if err != nil {
		r.l.Error(err, "http - v1 - correspondences - get")
		httpserver.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	correspondences, hasMoreData, err := r.ruc.GetCorrespondencesV1(
		entity.GetCorrespondencesData{
			Claims:     *claims,
			RevisionID: param.RevisionId,
		},
		entity.Pagination{Limit: *param.Limit,
			Offset: *param.Offset,
		},
	)
	if err != nil {
		r.l.Error(err, "http - v1 - correspondences - get")
		httpserver.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK,
		httpserver.ServerResponse[[]entity.Correspondence]{
			Data:        correspondences,
			HasMoreData: hasMoreData,
		},
	)
}

func (r *revisionRoutes) GetRevisionV1(c *gin.Context, params GetRevisionV1Params) {
	claims, err := r.swCtxAdapter.GetTokenClaims(c)
	if err != nil {
		r.l.Error(err, "http - v1 - revision - get")
		httpserver.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	revision, err := r.ruc.GetRevisionV1(claims, params.RevisionId)
	if err != nil {
		r.l.Error(err, "http - v1 - revision - get")
		httpserver.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, httpserver.ServerResponse[*entity.Revision]{Data: revision})
}
