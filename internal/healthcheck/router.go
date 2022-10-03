package healthcheck

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetHealthCheckRoute(r *gin.Engine) {
	r.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "OK")
	})
}
