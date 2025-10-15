package routers

import (
	"bff/internal/controllers"
)

func (r *BffRouter) MountPricesRoutes() {
	prices := r.router.Group("/api/prices")

	prices.GET("timeline", controllers.Controller.PriceTimeline)
	prices.GET("current", controllers.Controller.CurrentPrices)
	prices.GET("last", controllers.Controller.LastPrices)
	prices.GET("/download/lastprices/:productId", controllers.Controller.DownloadLastPrices)

}
