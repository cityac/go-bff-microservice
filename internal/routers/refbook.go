package routers

import (
	"bff/internal/controllers"
)

func (r *BffRouter) MountRefBookRoutes() {
	refbook := r.router.Group("/api/refbook")
	refbook.GET("", controllers.Controller.Refbook)
	refbook.GET("/countries", controllers.Controller.Countries)
	refbook.POST("/mccmnc", controllers.Controller.MccMnc)
}
