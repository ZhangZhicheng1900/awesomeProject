package readxls

import (
	"github.com/extrame/xls"
)

func ReadXls(filePath string) (res map[string][][]string) {
	res = make(map[string][][]string)
	xlFile, err := xls.Open(filePath, "utf-8")
	if err != nil {
		panic(err)
	}
	sheetCount := xlFile.NumSheets()

	for sheetI := 0; sheetI < sheetCount; sheetI++ {
		sheet := xlFile.GetSheet(sheetI)
		res[sheet.Name] = make([][]string, 0)

		rowsCount := sheet.RowsCount()

		if rowsCount != 0 {
			res[sheet.Name] = make([][]string, rowsCount)
			for i := 0; i < rowsCount; i++ {
				row := sheet.Row(i)
				data := make([]string, 0)
				for j := 0; j < row.LastCol(); j++ {
					col := row.Col(j)
					data = append(data, col)
				}
				res[sheet.Name][i] = data
			}
		}
	}
	return res
}
