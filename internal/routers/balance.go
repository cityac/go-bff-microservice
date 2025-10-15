package routers

import (
	"bff/internal/controllers"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/auth"
)

func (r *BffRouter) MountBalanceRouter() {
	balance := r.router.Group("/api/balance")

	All := auth.MiddlewareFactory(r.fga, auth.ModeAll)
	Any := auth.MiddlewareFactory(r.fga, auth.ModeOneOf)

	balance.POST("", All(auth.AccessPartnersCreator), controllers.Controller.CreatePartnerBalance)
	balance.GET("/:partnerId", Any(auth.AccessPartnersViewer, auth.AccessPartnersDraftViewer), controllers.Controller.GetPartnerBalance)
	balance.PUT("", All(auth.AccessPartnersEditor), controllers.Controller.UpdatePartnerBalance)

	balance.POST("/test", controllers.Controller.SendTestBalanceUpdate)
}
