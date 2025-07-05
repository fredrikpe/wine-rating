package internal

import (
	"fmt"
	"log"
	"net/http"

	"wine_rating/internal/db"
	vexcel "wine_rating/internal/vinmonopolet/excel"
	"wine_rating/internal/vivino"

	"github.com/xuri/excelize/v2"
)

func EnrichHandler(db *db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		firstEmptyColumn, err := vexcel.FirstEmptyColumn(excel)
		if err != nil {
			http.Error(w, fmt.Sprintf("fec error: %v", err), http.StatusBadRequest)
			return
		}
		col1Name, _ := excelize.ColumnNumberToName(firstEmptyColumn)
		col2Name, _ := excelize.ColumnNumberToName(firstEmptyColumn + 1)
		col3Name, _ := excelize.ColumnNumberToName(firstEmptyColumn + 2)

		for i, row := range rows[1:] {
			columnIndexes, err := vexcel.FindColumnIndexes(excel)
			if err != nil {
				http.Error(w, fmt.Sprintf("error: %v", err), http.StatusBadRequest)
				return
			}
			name := row[columnIndexes.NameCol]
			producer := row[columnIndexes.ProducerCol]
			vintage := vexcel.ResolveVintage(row, columnIndexes.VintageCol, name)
			log.Print("Name", name)

			match, err := vivino.FindMatch(db, name, producer, vintage)
			if err != nil {
				http.Error(w, fmt.Sprintf("error: %v", err), http.StatusInternalServerError)
				return
			}
			excel.SetCellValue(sheetName, fmt.Sprintf("%s%d", col1Name, i+2), match.RatingsAverage)
			excel.SetCellValue(sheetName, fmt.Sprintf("%s%d", col2Name, i+2), vivino.Url(match.Id))
			excel.SetCellValue(sheetName, fmt.Sprintf("%s%d", col3Name, i+2), match.Confidence)
		}

		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", "attachment; filename=\"enriched.xlsx\"")
		excel.Write(w)
	}
}
