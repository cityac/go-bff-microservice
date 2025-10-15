package routers

import (
	"bff/internal/controllers"
)

func (r *BffRouter) MountUploadRoutes() {
	r.router.POST("/api/vendors/price-updates/upload", controllers.AllVendorsHIDDENRequestsHandler)
	r.router.GET("/api/vendors/price-updates/upload/:jobID", controllers.AllVendorsHIDDENRequestsHandler)
	r.router.PUT("/api/vendors/price-updates/upload/:jobID", controllers.AllVendorsHIDDENRequestsHandler)
	r.router.GET("/api/vendors/price-updates/upload/options", controllers.AllVendorsHIDDENRequestsHandler)
	r.router.GET("/api/vendors/price-updates", controllers.AllVendorsHIDDENRequestsHandler)
	r.router.POST("/api/vendors/price-updates/confirm-upload", controllers.AllVendorsHIDDENRequestsHandler)
	r.router.GET("/api/vendors/price-updates/validation-report", controllers.AllVendorsHIDDENRequestsHandler)
	r.router.POST("/api/vendors/price-updates/structure", controllers.AllVendorsHIDDENRequestsHandler)
	r.router.POST("/api/vendors/price-updates/effective-from", controllers.AllVendorsHIDDENRequestsHandler)
	r.router.POST("/api/vendors/price-updates/confirm", controllers.AllVendorsHIDDENRequestsHandler)

	r.router.POST("/api/clients/price-updates/upload", controllers.AllClientsHIDDENRequestsHandler)
	r.router.GET("/api/clients/price-updates/upload/:jobID", controllers.AllClientsHIDDENRequestsHandler)
	r.router.PUT("/api/clients/price-updates/upload/:jobID", controllers.AllClientsHIDDENRequestsHandler)
	r.router.GET("/api/clients/price-updates/upload/options", controllers.AllClientsHIDDENRequestsHandler)
	r.router.GET("/api/clients/price-updates", controllers.AllClientsHIDDENRequestsHandler)
	r.router.POST("/api/clients/price-updates/confirm-upload", controllers.AllClientsHIDDENRequestsHandler)
	r.router.GET("/api/clients/price-updates/validation-report", controllers.AllClientsHIDDENRequestsHandler)
	r.router.POST("/api/clients/price-updates/structure", controllers.AllClientsHIDDENRequestsHandler)
	r.router.POST("/api/clients/price-updates/effective-from", controllers.AllClientsHIDDENRequestsHandler)
	r.router.POST("/api/clients/price-updates/confirm", controllers.AllClientsHIDDENRequestsHandler)
}
