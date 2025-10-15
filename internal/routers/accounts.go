package routers

import "bff/internal/controllers"

func (r *BffRouter) MountAccountsRouter() {
	ratePlans := r.router.Group("/api/accounts")
	ratePlans.GET("/:tenantName", controllers.Controller.Accounts)
}
