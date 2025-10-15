package routers

import (
	"bff/internal/controllers"
)

func (r *BffRouter) MountPriceListRouter() {
	priceLists := r.router.Group("/api/price-lists")
	priceLists.GET("/:productId", controllers.Controller.PriceListByProductId)
	priceLists.GET("/download/:productId", controllers.Controller.DownloadPriceList)

}
