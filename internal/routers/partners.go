package routers

import (
	"bff/internal/controllers"

	"gitlab.HIDDEN.com/HIDDEN/sms-core/auth"
)

func (r *BffRouter) MountPartnerRouter() {
	// Register bff router to query db broker for rate plan by productId
	// r.router.GET("/api/rate-plans", RatePlanByProductIdHandler)
	All := auth.MiddlewareFactory(r.fga, auth.ModeAll)
	Any := auth.MiddlewareFactory(r.fga, auth.ModeOneOf)

	partners := r.router.Group("/api/partners")
	//partnerProducts := r.router.Group("/api/partners")
	upload := r.router.Group("/api/upload")
	products := r.router.Group("/api/products")
	channels := r.router.Group("/api/channels")

	products.GET("/:productId", controllers.Controller.GetProductById)
	products.GET("", controllers.Controller.GetPartnersProducts)
	products.GET("/vendors", controllers.Controller.GetVendorProductsPrices)

	partners.GET("", controllers.Controller.Partners)
	partners.POST("", All(auth.AccessPartnersCreator), controllers.Controller.CreatePartner)
	partners.PUT("/:partnerId", All(auth.AccessPartnersEditor), controllers.Controller.UpdatePartner)
	partners.POST("/:partnerId/deactivate", All(auth.AccessPartnersToggler), controllers.Controller.DeactivatePartner)
	partners.PATCH("/:partnerId", All(auth.AccessPartnersToggler), controllers.Controller.TogglePartner)

	//partners.GET("/:partnerId/products", controllers.Controller.Products)
	partners.GET("/:partnerId", Any(auth.AccessPartnersViewer, auth.AccessPartnersDraftViewer), controllers.Controller.GetPartnerById)
	channels.GET("/:productId", Any(auth.AccessPartnersViewer, auth.AccessPartnersDraftViewer), controllers.Controller.GetChannelsActiveCount)

	partners.GET("/:partnerId/contacts", Any(auth.AccessPartnersViewer, auth.AccessPartnersDraftViewer), controllers.Controller.GetContactsByPartnerId)
	partners.POST("/:partnerId/contacts", All(auth.AccessPartnersCreator), controllers.Controller.CreateContact)
	partners.PUT("/:partnerId/contacts/:contactId", All(auth.AccessPartnersEditor), controllers.Controller.UpdateContact)
	partners.DELETE("/:partnerId/contacts/:contactId", All(auth.AccessPartnersEditor), controllers.Controller.DeleteContact)

	partners.GET("/:partnerId/products" /*Any(auth.AccessPartnersViewer, auth.AccessPartnersDraftViewer),*/, controllers.Controller.GetPartnerProductsById)
	partners.POST("/:partnerId/products", All(auth.AccessPartnersCreator), controllers.Controller.CreatePartnerProduct)
	partners.PUT("/:partnerId/products/:productId", All(auth.AccessPartnersEditor), controllers.Controller.UpdatePartnerProduct)
	partners.DELETE("/:partnerId/products/:productId", All(auth.AccessPartnersEditor), controllers.Controller.DeletePartnerProduct)

	partners.GET("/:partnerId/channels", Any(auth.AccessPartnersViewer, auth.AccessPartnersDraftViewer), controllers.Controller.GetChannelsByPartnerId)
	partners.GET("/:partnerId/products/:productId/channels", Any(auth.AccessPartnersViewer, auth.AccessPartnersDraftViewer), controllers.Controller.GetChannelsByProductId)
	partners.POST("/:partnerId/products/:productId/channels", All(auth.AccessPartnersCreator), controllers.Controller.CreateProductChannel)
	partners.PUT("/:partnerId/products/:productId/channels/:channelId", All(auth.AccessPartnersEditor), controllers.Controller.UpdateProductChannel)
	partners.DELETE("/:partnerId/products/:productId/channels/:channelId", All(auth.AccessPartnersEditor), controllers.Controller.DeleteProductChannel)

	upload.POST("", controllers.AllPartneHIDDENRequestHandler)
	upload.DELETE("", controllers.AllPartneHIDDENRequestHandler)
}
