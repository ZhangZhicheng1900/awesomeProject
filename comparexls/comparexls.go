package main

import (
	"fmt"
	"os"
	"reflect"

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
			for i := 0; i < int(rowsCount); i++ {
				row := sheet.Row(i)
				data := make([]string, 0)
				if row.LastCol() > 0 {
					for j := 0; j < row.LastCol(); j++ {
						col := row.Col(j)
						data = append(data, col)
					}
					res[sheet.Name][i] = data
				}
			}
		}
	}
	return res
}

func fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	os.Exit(1)
}

func main() {
	xls1 := os.Args[1]
	xls2 := os.Args[2]
	xls1Map := ReadXls(xls1)
	xls2Map := ReadXls(xls2)

	if len(xls1Map) != len(xls2Map) {
		fatalf("xls sheets len %s %d, %s %d", xls1, len(xls1Map), xls2, len(xls2Map))
	}
	for k, v := range xls1Map {
		v2, found := xls2Map[k]
		if !found {
			fatalf("sheet %s not found in xls %s", k, xls2)
		}

		lines1 := len(v)
		lines2 := len(v2)
		minLine := 0
		if lines1 >= lines2 {
			minLine = lines2
		} else {
			minLine = lines1
		}
		for i := 0; i < minLine; i++ {
			if !reflect.DeepEqual(v[i], v2[i]) {
				fatalf("sheet %s line %d not equal", k, i+1)
			}
		}
		for i := minLine; i < lines1; i++ {
			fatalf("sheet %s line %d exist in %s, but not in xls2", k, i+1, xls1)
		}
		for i := minLine; i < lines2; i++ {
			fatalf("sheet %s line %d exist in %s, but not in xls1", k, i+1, xls2)
		}
	}

	fmt.Println("ok, same")
}
