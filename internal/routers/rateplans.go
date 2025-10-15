package routers

import (
	"bff/internal/controllers"
)

func (r *BffRouter) MountRatePlanRouter() {
	ratePlans := r.router.Group("/api/rate-plans")
	ratePlans.GET("", controllers.Controller.RatePlans)
	ratePlans.POST("/link", controllers.Controller.LinkRatePlanToProduct)
	ratePlans.GET("/plan", controllers.Controller.RatePlan)
	ratePlans.GET("/download", controllers.Controller.DownloadRatePlan)
}
