package xlsxop

import (
	"fmt"
	"sort"

	excellizev2 "github.com/xuri/excelize/v2"
)

const defaultSheet1 = "Sheet1"

func Read(filePath string) map[string][][]string {
	var res = make(map[string][][]string)
	excelEntryF, excelEntryErr := excellizev2.OpenFile(filePath)
	if excelEntryErr != nil {
		panic(excelEntryErr)
	}

	sheets := excelEntryF.GetSheetList()
	for _, sheetName := range sheets {
		sheetRows, err := excelEntryF.GetRows(sheetName)
		if err != nil {
			panic(err)
		}
		res[sheetName] = sheetRows
	}

	return res
}

type ExcelOpFunc func(*excellizev2.File) error

func Write(filePath string, content map[string][][]string, fn ExcelOpFunc) error {
	var sheets = make([]string, 0)
	for name, _ := range content {
		sheets = append(sheets, name)
	}
	sort.Strings(sheets)

	xlsxFile := excellizev2.NewFile()
	for _, sheetName := range sheets {
		sheetContent := content[sheetName]
		xlsxFile.NewSheet(sheetName)

		for i := 0; i < len(sheetContent); i++ {
			if err := xlsxFile.SetSheetRow(sheetName, fmt.Sprintf("A%d", i+1), &sheetContent[i]); err != nil {
				return fmt.Errorf("SetSheetRow error, %v", err)
			}
		}
	}
	if _, found := content[defaultSheet1]; !found && len(content) != 0 {
		xlsxFile.DeleteSheet(defaultSheet1)
	}

	if fn != nil {
		if err := fn(xlsxFile); err != nil {
			return err
		}
	}

	return xlsxFile.SaveAs(filePath)
}
