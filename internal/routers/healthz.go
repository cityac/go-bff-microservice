package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (r *BffRouter) MountHealthzRouter() {
	user := r.router.Group("/healthz")
	user.GET("", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
}
