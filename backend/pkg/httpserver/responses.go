package httpserver

import "github.com/gin-gonic/gin"

// TODO: закинуть методы по работе с респонсами
type ServerResponse[T any] struct {
	Data        T       `json:"data"`
	Error       *string `json:"error,omitempty"`
	HasMoreData *bool   `json:"has_more_data,omitempty"`
}

func ErrorResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, ServerResponse[any]{
		Error: &msg,
	})
}
