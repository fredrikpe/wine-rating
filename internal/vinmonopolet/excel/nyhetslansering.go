package excel

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

type ColumnIndexes struct {
	NameCol     int
	ProducerCol int
	VintageCol  *int
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
