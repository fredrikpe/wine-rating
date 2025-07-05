package excel

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type ColumnIndexes struct {
	NameCol     int
	ProducerCol int
	VintageCol  *int
}

func ResolveVintage(row []string, vintageCol *int, name string) *int {
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

func FirstEmptyColumn(f *excelize.File) (int, error) {
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

func FindColumnIndexes(f *excelize.File) (ColumnIndexes, error) {
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
