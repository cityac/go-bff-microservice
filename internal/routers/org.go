package routers

import (
	"bff/internal/controllers"
)

func (r *BffRouter) MountOrgRouter() {
	org := r.router.Group("/api/org")
	org.POST("", controllers.AllUserRequestsHandler)
	org.GET("", controllers.AllUserRequestsHandler)
}
