package healthcheck

import (
	"kredit-plus/app/constants"
	"kredit-plus/app/controller"
	"kredit-plus/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type IHealthCheckController interface {
	HealthCheck(c *gin.Context)
}

type HealthCheckController struct {
}

func NewHealthCheckController() IHealthCheckController {
	return &HealthCheckController{}
}

func (h *HealthCheckController) HealthCheck(c *gin.Context) {
	message := "Congratulations! Your application is as fit as a fiddle, and it's ready to dance on the cloud. 🕺☁️"
	data := gin.H{
		"version":   constants.Config.ProjectVersion,
		"uptime":    time.Since(config.StartTime).String(),
		"timestamp": time.Now().Format(time.RFC3339),
	}
	controller.RespondWithSuccess(c, http.StatusOK, message, data, nil)
}
