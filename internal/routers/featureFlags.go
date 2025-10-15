package routers

import "bff/internal/controllers"

func (r *BffRouter) MountFeatureRouter() {
	features := r.router.Group("/api/features")

	features.GET("", controllers.Controller.GetFeatureFlags)
	features.POST("/replace/:orgId", controllers.Controller.ReplaceAllFeatureFlagUsers)
	features.POST("/:orgId/:userId", controllers.Controller.UpdateFeatureFlag)
}
