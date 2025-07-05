package web

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"wine_rating/internal/db"
	"wine_rating/internal/similarity"
	"wine_rating/internal/vivino"

	"github.com/xuri/excelize/v2"
)

type OutputColumns struct {
	RatingCol string
	URLCol    string
	ConfCol   string
}

func readUploadedExcel(r *http.Request) (*excelize.File, error) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return nil, fmt.Errorf("failed to parse form: %w", err)
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("invalid file: %w", err)
	}
	defer file.Close()

	excel, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Excel: %w", err)
	}
	return excel, nil
}

func enrichExcelWithVivino(db *db.Store, excel *excelize.File) error {
	sheetName := excel.GetSheetName(0)
	rows, err := excel.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("reading rows failed: %w", err)
	}

	firstEmptyColumn, err := firstEmptyColumn(excel)
	if err != nil {
		return fmt.Errorf("first empty column: %w", err)
	}
	col1Name, _ := excelize.ColumnNumberToName(firstEmptyColumn)
	col2Name, _ := excelize.ColumnNumberToName(firstEmptyColumn + 1)
	col3Name, _ := excelize.ColumnNumberToName(firstEmptyColumn + 2)
	outColumns := OutputColumns{
		RatingCol: col1Name,
		URLCol:    col2Name,
		ConfCol:   col3Name,
	}

	columnIndexes, err := findColumnIndexes(excel)
	if err != nil {
		return fmt.Errorf("couldn't find columns: %w", err)
	}
	for i, row := range rows[1:] {
		if err := enrichRow(db, excel, columnIndexes, sheetName, row, i+2, outColumns); err != nil {
			log.Printf("Row %d failed: %v", i+2, err)
		}
	}
	return nil
}

func enrichRow(db *db.Store, excel *excelize.File, columnIndexes ColumnIndexes, sheetName string, row []string, rowNum int, outCols OutputColumns) error {
	name := row[columnIndexes.NameCol]
	producer := row[columnIndexes.ProducerCol]
	vintage := resolveVintage(row, columnIndexes.VintageCol, name)

	match, err := vivino.FindMatch(db, name, producer, vintage)
	if err != nil {
		return fmt.Errorf("find match: %w", err)
	}

	cell := func(col string) string {
		return fmt.Sprintf("%s%d", col, rowNum)
	}

	if !similarity.HighEnough(match.Similarity) {
		return nil
	}

	if match.RatingsAverage != nil {
		excel.SetCellValue(sheetName, cell(outCols.RatingCol), *match.RatingsAverage)
	} else {
		excel.SetCellValue(sheetName, cell(outCols.RatingCol), "n/a")
	}

	excel.SetCellValue(sheetName, cell(outCols.URLCol), match.Url)
	excel.SetCellValue(sheetName, cell(outCols.ConfCol), fmt.Sprintf("%.2f", match.Similarity))

	return nil
}

func firstEmptyColumn(f *excelize.File) (int, error) {
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil || len(rows) == 0 {
		return 0, fmt.Errorf("Failed to read rows or empty sheet %w", err)
	}
	for col := 1; col < 1000; col++ {
		colName, err := excelize.ColumnNumberToName(col)
		if err != nil {
			return 0, fmt.Errorf("Failed column name: %w", err)
		}
		cellRef := fmt.Sprintf("%s%d", colName, 1)
		val, err := f.GetCellValue(sheetName, cellRef)
		if err != nil {
			return 0, fmt.Errorf("Failed get cell value: %w", err)
		}
		if val == "" {
			return col, nil
		}
	}
	return 0, fmt.Errorf("no empty column found in first 1000 row 0")
}

type ColumnIndexes struct {
	NameCol     int
	ProducerCol int
	VintageCol  *int
}

func resolveVintage(row []string, vintageCol *int, name string) *int {
	if vintageCol != nil {
		vintage, err := strconv.Atoi(row[*vintageCol])
		if err != nil {
			return &vintage
		}
	}
	if vintage, ok := extractYear(name); ok {
		return &vintage
	}
	return nil
}

func findColumnIndexes(f *excelize.File) (ColumnIndexes, error) {
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil || len(rows) == 0 {
		return ColumnIndexes{}, fmt.Errorf("Failed to read rows or empty sheet %w", err)
	}

	header := rows[0]
	var result ColumnIndexes

	foundName := false
	foundProducer := false
	for i, col := range header {
		normalized := strings.ToLower(strings.TrimSpace(col))

		switch normalized {
		case "artikkelnavn", "produktnavn":
			result.NameCol = i
			foundName = true
		case "produsent":
			result.ProducerCol = i
			foundProducer = true
		case "årgang":
			result.VintageCol = new(int)
			*result.VintageCol = i
		}
	}
	if !foundName {
		return ColumnIndexes{}, fmt.Errorf("couldn't find name column")
	}
	if !foundProducer {
		return ColumnIndexes{}, fmt.Errorf("couldn't find producer column")
	}

	return result, nil
}

func extractYear(text string) (int, bool) {
	re := regexp.MustCompile(`\b\d{4}\b`)
	matches := re.FindAllString(text, -1)
	for _, m := range matches {
		year, _ := strconv.Atoi(m)
		if year >= 1500 && year <= 2100 {
			return year, true
		}
	}
	return 0, false
}
