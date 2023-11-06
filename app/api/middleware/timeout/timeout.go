package timeout

import (
	"kredit-plus/app/constants"
	"kredit-plus/app/controller"
	"net/http"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

// TimeoutMiddleware sets a timeout for incoming requests and responds with an error if the request takes too long.
func TimeoutMiddleware() gin.HandlerFunc {
	return timeout.New(timeout.WithTimeout(15*time.Second), timeout.WithResponse(timeoutHandler))
}

func timeoutHandler(c *gin.Context) {
	controller.RespondWithError(c, http.StatusRequestTimeout, constants.TIMEOUT_ERROR, nil)
}
