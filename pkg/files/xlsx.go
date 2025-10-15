package files

import (
	"bff/pkg/domain/model"
	"bytes"
	"fmt"
	"github.com/xuri/excelize/v2"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
	"time"
)

func PriceListToXLSBuffer(prices model.FormattedPriceList, partnerName string, productName string) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	now := time.Now().Format("2006-01-02_15-04-05")
	fileName := fmt.Sprintf("%s - %s_%s.xlsx", partnerName, productName, now)

	f := excelize.NewFile()

	// Add headers
	headers := []string{"MCC", "MNC", "Country", "Network", "Effective date", "Rate", "Change type"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue("Sheet1", cell, header)
	}

	// Style for headers
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#D3D3D3"}, Pattern: 1},
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create header style")
		return nil, err
	}

	// Apply header style
	for i := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellStyle("Sheet1", cell, cell, headerStyle)
	}

	// Add data
	for i, row := range prices.Rows {
		rowNum := i + 2
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", rowNum), row.MCC)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", rowNum), row.MNC)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", rowNum), row.Country)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", rowNum), row.Network)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", rowNum), row.EffectiveDate.Format("2006-01-02 15:04"))
		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", rowNum), row.Rate)
		f.SetCellValue("Sheet1", fmt.Sprintf("G%d", rowNum), row.ChangeType)
	}

	// Auto-fit columns
	for i := range headers {
		col := string('A' + i)
		width, _ := f.GetColWidth("Sheet1", col)
		f.SetColWidth("Sheet1", col, col, width+2)
	}

	// Set file properties
	f.SetDocProps(&excelize.DocProperties{
		Title:    fileName,
		Created:  now,
		Modified: now,
	})

	err = f.Write(&buf)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
