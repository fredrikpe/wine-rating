package internal

import (
	"fmt"
	"log"
	"net/http"

	vexcel "wine_rating/internal/vinmonopolet/excel"
	"wine_rating/internal/vivino"

	"github.com/xuri/excelize/v2"
)

func EnrichHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	excel, err := excelize.OpenReader(file)
	if err != nil {
		http.Error(w, "Failed to parse Excel", http.StatusInternalServerError)
		return
	}

	sheetName := excel.GetSheetName(0)
	rows, _ := excel.GetRows(sheetName)

	for i, row := range rows[1:] {
		columnIndexes, err := vexcel.FindColumnIndexes(excel)
		if err != nil {
			http.Error(w, fmt.Sprintf("error: %v", err), http.StatusBadRequest)
			return
		}
		name := row[columnIndexes.NameCol]
		producer := row[columnIndexes.ProducerCol]
		log.Print("Name", name)

		hit, err := vivino.FindVivinoMatch(name, producer)
		if err != nil {
			http.Error(w, fmt.Sprintf("error: %v", err), http.StatusInternalServerError)
			return
		}

		colCount := len(row)
		col1Name, _ := excelize.ColumnNumberToName(colCount + 1)
		col2Name, _ := excelize.ColumnNumberToName(colCount + 2)

		excel.SetCellValue(sheetName, fmt.Sprintf("%s%d", col1Name, i+2), hit.Statistics.RatingsAverage)
		excel.SetCellValue(sheetName, fmt.Sprintf("%s%d", col2Name, i+2), vivino.Url(hit.Id))
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=\"enriched.xlsx\"")
	excel.Write(w)
}
