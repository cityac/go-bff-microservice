package routers

import (
	"bff/internal/controllers"
)

func (r *BffRouter) MountRulesRoutes() {
	r.router.GET("/api/rules", controllers.AllVendorsHIDDENRequestsHandler)
	r.router.POST("/api/rules", controllers.AllVendorsHIDDENRequestsHandler)
	r.router.PUT("/api/rules", controllers.AllVendorsHIDDENRequestsHandler)
	r.router.GET("/api/cards", controllers.AllVendorsHIDDENRequestsHandler)
	r.router.DELETE("/api/rules/:ruleID", controllers.AllVendorsHIDDENRequestsHandler)

}
