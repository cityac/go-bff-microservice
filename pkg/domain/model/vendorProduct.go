package model

type VendorProduct struct {
	Id           string               `json:"id"`
	Name         string               `json:"name"`
	MccMnc       string               `json:"mccmnc"`
	Network      []string             `json:"network"`
	Rate         string               `json:"rate"`
	VendorPrices []VendorProductPrice `json:"vendorPrices"`
}

type VendorProductPrice struct {
	ID      string `json:"id"`
	Mccmnc  string `json:"mccmnc"`
	Network string `json:"network"`
	Rate    string `json:"rate"`
}
