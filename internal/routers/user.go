package routers

import (
	"bff/internal/controllers"
)

func (r *BffRouter) MountUserRouter() {
	user := r.router.Group("/api/user")
	users := r.router.Group("/api/users")
	permissions := r.router.Group("/api/permissions")
	presets := permissions.Group("/presets")

	user.GET("", controllers.Controller.GetUser)
	user.GET("/:id", controllers.AllUserRequestsHandler)
	user.POST("", controllers.AllUserRequestsHandler)
	user.POST("/toggle", controllers.AllUserRequestsHandler)
	user.PUT("/:id", controllers.AllUserRequestsHandler)
	user.GET("/HIDDEN-admin", controllers.AllUserRequestsHandler)
	user.POST("/HIDDEN-admin", controllers.AllUserRequestsHandler)

	users.GET("", controllers.AllUserRequestsHandler)
	users.GET("account-managers", controllers.AllUserRequestsHandler)

	presets.GET(":id", controllers.AllUserRequestsHandler)
	presets.POST("", controllers.AllUserRequestsHandler)
	presets.PUT(":id", controllers.AllUserRequestsHandler)
	presets.GET("", controllers.AllUserRequestsHandler)
}
