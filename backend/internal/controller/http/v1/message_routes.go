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

type messageRoutes struct {
	muc          usecase.IMessageUseCase
	l            logger.ILogger
	swCtxAdapter swcontext.SwGinContext
}

func NewMessageRoutes(muc usecase.IMessageUseCase, l logger.ILogger, swCtxAdapter swcontext.SwGinContext) *messageRoutes {
	return &messageRoutes{
		muc:          muc,
		l:            l,
		swCtxAdapter: swCtxAdapter,
	}
}

func (m *messageRoutes) GetMessagesV1(c *gin.Context, params GetMessagesV1Params) {
	claims, err := m.swCtxAdapter.GetTokenClaims(c)
	if err != nil {
		m.l.Error(err, "http - v1 - messages - get")
		httpserver.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	messages, hasMoreData, err := m.muc.GetMessagesV1(claims, params.DisputeId, *params.Limit, *params.Offset)
	if err != nil {
		m.l.Error(err, "http - v1 - messages - get")
		httpserver.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, httpserver.ServerResponse[[]entity.Message]{Data: messages, HasMoreData: hasMoreData})
}

func (m *messageRoutes) CreateMessageV1(c *gin.Context) {
	claims, err := m.swCtxAdapter.GetTokenClaims(c)
	if err != nil {
		m.l.Error(err, "http - v1 - message - create")
		httpserver.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		m.l.Error(err, "http - v1 - message - create")
		httpserver.ErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	var createMsgReqBody entity.CreateMessageV1MultipartBody
	createMsgReqBody.DisputeId = form.Value["dispute_id"][0]
	createMsgReqBody.MessageBody = &form.Value["message_body"][0]
	if len(form.File["file"]) == 1 {
		file := &openapi_types.File{}
		file.InitFromMultipart(form.File["file"][0])
		createMsgReqBody.File = file
	} else {
		createMsgReqBody.File = nil
	}

	msgId, err := m.muc.CreateMessageV1(claims, createMsgReqBody)
	if err != nil {
		m.l.Error(err, "http - v1 - message - create")
		httpserver.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, httpserver.ServerResponse[entity.CreateMsgResponse]{Data: *msgId})
}
