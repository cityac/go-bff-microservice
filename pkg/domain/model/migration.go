package model

type MigrationResult struct {
	Error           error
	Message         string
	ProductId       string
	ProductCurrency string
}
