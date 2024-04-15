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

type disputeRoutes struct {
	duc          usecase.IDisputeUseCase
	l            logger.ILogger
	swCtxAdapter swcontext.SwGinContext
}

func NewDisputeRoutes(duc usecase.IDisputeUseCase, l logger.ILogger, swCtxAdapter swcontext.SwGinContext) *disputeRoutes {
	return &disputeRoutes{
		duc:          duc,
		l:            l,
		swCtxAdapter: swCtxAdapter,
	}
}

func (r *disputeRoutes) GetDisputesV1(c *gin.Context, params GetDisputesV1Params) {
	claims, err := r.swCtxAdapter.GetTokenClaims(c)
	if err != nil {
		r.l.Error(err, "http - v1 - dispute - claims")
		httpserver.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	disputes, hasMoreData, err := r.duc.GetDisputesV1(*claims, *params.Limit, *params.Offset, entity.EntityStatus(*params.Status))
	if err != nil {
		r.l.Error(err, "http - v1 - disputes")
		httpserver.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, httpserver.ServerResponse[[]entity.DisputeList]{Data: disputes, HasMoreData: hasMoreData})
}

func (r *disputeRoutes) GetDisputeV1(c *gin.Context, params GetDisputeV1Params) {
	dispute, err := r.duc.GetDisputeV1(params.DisputeId, params.GoodsId)
	if err != nil {
		r.l.Error(err, "http - v1 - dispute - get")
		if dispute == nil {
			httpserver.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, httpserver.ServerResponse[*entity.Dispute]{Data: dispute})
}

func (r *disputeRoutes) CreateDisputeV1(c *gin.Context) {
	claims, err := r.swCtxAdapter.GetTokenClaims(c)
	if err != nil {
		r.l.Error(err, "http - v1 - dispute - create")
		httpserver.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		r.l.Error(err, "http - v1 - dispute - create")
		httpserver.ErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	var createDisputeReqBody entity.CreateDisputeV1MultipartBody
	createDisputeReqBody.ShortageId = form.Value["shortage_id"][0]
	createDisputeReqBody.OrganizationId = form.Value["organization_id"][0]
	createDisputeReqBody.MessageBody = form.Value["message_body"][0]
	if len(form.File["file"]) == 1 {
		file := &openapi_types.File{}
		file.InitFromMultipart(form.File["file"][0])
		createDisputeReqBody.File = file
	} else {
		createDisputeReqBody.File = nil
	}

	dispute, err := r.duc.CreateDisputeV1(claims, &createDisputeReqBody)
	if err != nil {
		r.l.Error(err, "http - v1 - dispute - create")
		httpserver.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, httpserver.ServerResponse[*entity.CreateDispute]{Data: dispute})
}

func (r *disputeRoutes) CloseDisputeV1(c *gin.Context) {
	claims, err := r.swCtxAdapter.GetTokenClaims(c)
	if err != nil {
		r.l.Error(err, "http - v1 - dispute - create")
		httpserver.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var closeDisputeReqBody CloseDisputeV1JSONBody
	if err := c.Bind(&closeDisputeReqBody); err != nil {
		r.l.Error(err, "http - v1 - dispute - close")
		httpserver.ErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	dispute, err := r.duc.CloseDisputeV1(claims, closeDisputeReqBody.DisputeId, *closeDisputeReqBody.GuiltyWorkerIds)
	if err != nil {
		r.l.Error(err, "http - v1 - dispute - close")
		if dispute == nil {
			httpserver.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, httpserver.ServerResponse[*entity.CloseDispute]{Data: dispute})
}

func (r *disputeRoutes) GetShortagesV1(c *gin.Context, param GetShortagesV1Params) {
	claims, err := r.swCtxAdapter.GetTokenClaims(c)
	if err != nil {
		r.l.Error(err, "http - v1 - shortages")
		httpserver.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	shortages, hasMoreData, err := r.duc.GetShortagesV1(claims.UserId, *param.Limit, *param.Offset)
	if err != nil {
		r.l.Error(err, "http - v1 - shortages")
		if shortages == nil {
			httpserver.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, httpserver.ServerResponse[[]entity.Shortage]{Data: shortages, HasMoreData: hasMoreData})
}
