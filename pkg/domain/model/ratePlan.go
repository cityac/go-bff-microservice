package model

import "time"

type RatePlan struct {
	Id             string    `json:"id"`
	Name           string    `json:"name"`
	OrganizationId string    `json:"organizationId"`
	PlanType       string    `json:"planType"`
	Prices         []Price   `json:"prices"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type CustomRatePlan struct {
	RatePlan
	PartnerId            string `json:"partnerId"`
	ProductId            string `json:"productId"`
	SystemCountriesCount int    `json:"systemCountriesCount"`
}
type GlobalRatePlan struct {
	RatePlan
	Currency            string `json:"currency"`
	LinkedProductsCount int    `json:"linkedProductsCount"`
	PlanBackwardId      string `json:"planBackwardId"`
}
