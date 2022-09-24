package controller

import (
	"context"
	"net/http"

	"github.com/ydataai/aws-adapter/pkg/service"
	"github.com/ydataai/go-core/pkg/common/config"
	"github.com/ydataai/go-core/pkg/common/logging"
	"github.com/ydataai/go-core/pkg/common/server"

	"github.com/gin-gonic/gin"
)

// RESTController defines rest controller
type RESTController struct {
	log           logging.Logger
	restService   service.RESTServiceInterface
	configuration config.RESTControllerConfiguration
}

// NewRESTController initializes rest controller
func NewRESTController(
	log logging.Logger,
	restService service.RESTServiceInterface,
	configuration config.RESTControllerConfiguration,
) RESTController {
	return RESTController{
		restService:   restService,
		log:           log,
		configuration: configuration,
	}
}

// Boot ...
func (r RESTController) Boot(s *server.Server) {
	s.Router.GET("/healthz", r.healthCheck())
	s.Router.GET("/available/gpu", r.getAvailableGPU())
}

func (r RESTController) healthCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	}
}

func (r RESTController) getAvailableGPU() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tCtx, cancel := context.WithTimeout(ctx, r.configuration.HTTPRequestTimeout)
		defer cancel()

		gpu, err := r.restService.AvailableGPU(tCtx)
		if err != nil {
			r.log.Errorf("while fetching available resources. Error: %s", err.Error())
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gpu)
	}
}
