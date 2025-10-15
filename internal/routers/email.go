package routers

import controllers "bff/internal/controllers"

func (r *BffRouter) MountEmailRouter() {
	email := r.router.Group("/api/email")
	email.POST("/send", controllers.Controller.SendEmail)
}
