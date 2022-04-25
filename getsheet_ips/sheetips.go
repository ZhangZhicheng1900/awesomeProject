package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"awesomeProject/readxls"
)

func main() {
	xlsPath := flag.String("xls-path", "", "xls excel file path")
	sheets := flag.String("sheets", "Sheet1\nSheet2", "sheet names, one line for one sheet")
	needFirstLine := flag.Bool("need-first-line", false, "need first line")
	raw := flag.Bool("raw", false, "raw data")
	flag.Parse()

	res := readxls.ReadXls(*xlsPath)
	if *raw {
		bytes, _ := json.MarshalIndent(res, "", "\t")
		fmt.Printf("%s", string(bytes))
		return
	}

	sheetsArr := strings.Split(*sheets, "\n")
	for _, sheet := range sheetsArr {
		sheetContent, found := res[sheet]
		if !found {
			panic(fmt.Sprintf("sheet %s not found", sheet))
		}

		count := 0
		for i, line := range sheetContent {
			if !*needFirstLine && i == 0 {
				continue
			}
			if len(line) >= 2 {
				fmt.Println(line[1])
				count++
			}
		}
		fmt.Printf("------sheet %s lines count %d (include first line), ip count %d-----\n", sheet, len(sheetContent), count)

		fmt.Println()
	}
}
