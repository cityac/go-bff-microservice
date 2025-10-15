package routers

import (
	"time"

	"bff/internal/controllers"
	"bff/internal/middlewares/debug"
	"bff/internal/utils"

	"github.com/gin-gonic/gin"
)

type HIDDENRouter struct {
	*gin.Engine
}

func (r *BffRouter) MountHIDDENRoutes() {
	HIDDEN := r.router.Group("/api/HIDDEN")
	globalHIDDEN := r.router.Group("/api/global-HIDDEN")
	canvas := r.router.Group("/api/canvas")
	da := r.router.Group("/api/da")
	sa := r.router.Group("/api/sa")
	vg := r.router.Group("/api/vendorGroups")
	clientGroups := r.router.Group("/api/client-groups")
	customerProducts := r.router.Group("/api/customerProducts")
	smsNetworksAndCodes := r.router.Group("/api/smsNetworksAndCodes")
	vendorProducts := r.router.Group("/api/vendorProducts")
	hierarchy := r.router.Group("/api/hierarchy")
	countries := r.router.Group("/api/countries")
	margins := r.router.Group("/api/margins")
	tech := r.router.Group("/api/tech")
	migration := r.router.Group("/api/migration")

	HIDDEN.GET("", controllers.AllHIDDENRequestsHandler)
	HIDDEN.POST("", controllers.AllHIDDENRequestsHandler)

	HIDDEN.PUT(
		"/:id",
		r.lockRule(lockRuleEntity, 3*time.Minute),
		controllers.AllHIDDENRequestsHandler,
		debug.LockTimeout,
		r.unlockFromCtx(),
	)

	HIDDEN.POST("/checkMsgCount/:id", controllers.AllHIDDENRequestsHandler)
	HIDDEN.GET("/checkMmcStatus", controllers.AllHIDDENRequestsHandler)
	HIDDEN.DELETE("/:id", controllers.AllHIDDENRequestsHandler)

	globalHIDDEN.GET("", controllers.AllHIDDENRequestsHandler)
	globalHIDDEN.POST("", controllers.AllHIDDENRequestsHandler)
	globalHIDDEN.PUT(
		"/:id",
		r.lockRule(lockGlobalRuleEntity, 5*time.Minute),
		controllers.AllHIDDENRequestsHandler,
		debug.LockTimeout,
		r.unlockFromCtx(),
	)
	globalHIDDEN.GET("/checkMmcStatus", controllers.AllHIDDENRequestsHandler)
	globalHIDDEN.PUT(
		"/customerProducts",
		lockByPayload[masterRuleIds](r.redis.GetClient(), lockGlobalRuleEntity, 5*time.Minute),
		proxyLockTokenHeader,
		controllers.AllHIDDENRequestsHandler,
		debug.LockTimeout,
		r.unlockFromCtx(),
	)
	globalHIDDEN.PUT(
		"/multi-edit",
		lockByPayload[ids](r.redis.GetClient(), lockGlobalRuleEntity, 5*time.Minute),
		proxyLockTokenHeader,
		controllers.AllHIDDENRequestsHandler,
		debug.LockTimeout,
		r.unlockFromCtx(),
	)
	globalHIDDEN.GET("/countriesAndCpsMappings", controllers.AllHIDDENRequestsHandler)
	globalHIDDEN.DELETE("/:id", controllers.AllHIDDENRequestsHandler)

	canvas.GET("", controllers.AllHIDDENRequestsHandler)

	da.GET("", controllers.AllHIDDENRequestsHandler)
	da.GET("/:fName", controllers.AllHIDDENRequestsHandler)
	da.GET("/file/:bucketId", controllers.AllHIDDENRequestsHandler)
	da.POST("/upload", controllers.AllHIDDENRequestsHandler)
	da.PUT(
		"/:id",
		r.lockById(lockDaEntity, 3*time.Minute),
		controllers.AllHIDDENRequestsHandler,
		debug.LockTimeout,
		r.unlockFromCtx(),
	)
	da.GET("/check/:daId", controllers.AllHIDDENRequestsHandler)
	da.POST("/check/:daId", controllers.AllHIDDENRequestsHandler)
	da.GET("/connections/:daId", controllers.AllHIDDENRequestsHandler)
	da.DELETE("/:daId", controllers.AllHIDDENRequestsHandler)

	sa.GET("", controllers.AllHIDDENRequestsHandler)
	sa.GET("/file/:bucketId", controllers.AllHIDDENRequestsHandler)
	sa.GET("/:fName", controllers.AllHIDDENRequestsHandler)
	sa.POST("/upload", controllers.AllHIDDENRequestsHandler)
	sa.GET("/connections-with-details/:saId", controllers.AllHIDDENRequestsHandler)
	sa.DELETE("/:saId", controllers.AllHIDDENRequestsHandler)

	clientGroups.GET("", controllers.AllHIDDENRequestsHandler)
	clientGroups.POST("", controllers.AllHIDDENRequestsHandler)
	clientGroups.GET("/:id", controllers.AllHIDDENRequestsHandler)
	clientGroups.PUT("/:id", controllers.AllHIDDENRequestsHandler)
	clientGroups.DELETE("/:id", controllers.AllHIDDENRequestsHandler)
	clientGroups.GET("/:id/products", controllers.AllHIDDENRequestsHandler)
	clientGroups.POST("/preview-products", controllers.AllHIDDENRequestsHandler)

	vg.GET("", controllers.AllHIDDENRequestsHandler)
	vg.GET("/:vgName", controllers.AllHIDDENRequestsHandler)
	vg.POST("", controllers.AllHIDDENRequestsHandler)
	vg.PUT(
		"/:id",
		r.lockById(lockVGEntity, 3*time.Minute),
		proxyLockTokenHeader,
		controllers.AllHIDDENRequestsHandler,
		debug.LockTimeout,
		r.unlockFromCtx(),
	)
	vg.DELETE(
		"/:id",
		r.lockById(lockVGEntity, 3*time.Minute),
		controllers.AllHIDDENRequestsHandler,
		debug.LockTimeout,
		r.unlockFromCtx(),
	)
	vg.GET("/connections/:id", controllers.AllHIDDENRequestsHandler)

	customerProducts.GET("", controllers.AllHIDDENRequestsHandler)
	customerProducts.GET("/:id", controllers.AllHIDDENRequestsHandler)
	customerProducts.GET("/unavailable", controllers.AllHIDDENRequestsHandler)

	smsNetworksAndCodes.GET("", controllers.AllHIDDENRequestsHandler)
	smsNetworksAndCodes.GET("/global", controllers.AllHIDDENRequestsHandler)
	smsNetworksAndCodes.POST("/canvas", controllers.AllHIDDENRequestsHandler)

	vendorProducts.GET("", controllers.AllHIDDENRequestsHandler)
	vendorProducts.GET("/:id", controllers.AllHIDDENRequestsHandler)

	hierarchy.GET("", controllers.AllHIDDENRequestsHandler)

	countries.GET("", controllers.AllHIDDENRequestsHandler)
	countries.GET("/unavailable", controllers.AllHIDDENRequestsHandler)
	//countries.GET("/raw", controllers.AllHIDDENRequestsHandler)

	//
	tech.POST("/check", controllers.AllHIDDENRequestsHandler)
	tech.POST("/seed/rules", controllers.AllHIDDENRequestsHandler)
	tech.POST("/mr_update", controllers.AllHIDDENRequestsHandler)
	tech.POST("/sr_update", controllers.AllHIDDENRequestsHandler)
	tech.POST("/customerProducts/bulk/:amount", controllers.AllHIDDENRequestsHandler)

	margins.POST("/calculate", utils.FeatureFlag("HIDDEN_MARGIN_LEGACY_MODE", controllers.AllHIDDENRequestsHandler, controllers.Controller.CalculateMargins))

	migration.POST("/mrsCount", controllers.Controller.MigrateMRsCount)
	migration.POST("/mrsCountByIds", controllers.Controller.MigrateMRsCountByIds)

	migration.POST("/ratePlans/cgrates-to-price-catalog", controllers.Controller.MigrateIzonetRatePlansFromCGRates)
	migration.POST("/ratePlans/cgrates-to-price-catalog/linkage", controllers.Controller.MigrateGlobalProductsLinkageFromCGRates)
	migration.POST("/partners/migrate", controllers.Controller.MigrateIzonetPartners)
}
