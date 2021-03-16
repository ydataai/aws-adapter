package server

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s Server) tracing() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := uuid.New().String()
		if ctx.Request.Header.Get("X-Request-Id") != "" {
			requestID = ctx.Request.Header.Get("X-Request-Id")
		}

		s.log.Infof("Path: %v", ctx.Request.URL.Path)
		s.log.Infof("X-Request-Id: %v", requestID)
		ctx.Set("X-Request-Id", requestID)
		ctx.Next()
	}
}
