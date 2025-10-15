package routers

import (
	"bff/internal/controllers"
)

func (r *BffRouter) DictionariesRouter() {
	dictionaries := r.router.Group("/api/dictionaries")

	dictionaries.GET("/contract-types", controllers.Controller.ContractTypes)
	dictionaries.GET("/account-managers", controllers.Controller.AccountManagers)
	dictionaries.GET("/business-units", controllers.Controller.BusinessUnits)
	dictionaries.GET("/sub-units", controllers.Controller.SubUnits)
	dictionaries.GET("/units-map", controllers.Controller.UnitsMap)
	dictionaries.GET("/countries", controllers.Controller.CountriesList)
	dictionaries.GET("/contract-companies", controllers.Controller.ContractCompanies)
	dictionaries.GET("/timezones", controllers.Controller.Timezones)
	dictionaries.GET("/HIDDEN-models", controllers.Controller.HIDDENModels)
	dictionaries.GET("/HIDDEN-periods", controllers.Controller.HIDDENPeriods)
	dictionaries.GET("/contact-types", controllers.Controller.ContactTypes)
	dictionaries.GET("/tags", controllers.Controller.Tags)
	dictionaries.GET("/partners", controllers.Controller.PartnerOptions)
	dictionaries.GET("/partners/:partnerId/products", controllers.Controller.PartnerProductOptions)
	dictionaries.GET("/presets", controllers.AllUserRequestsHandler)
	dictionaries.GET("/currencies", controllers.Controller.Currencies)
}
