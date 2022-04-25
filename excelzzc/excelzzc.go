package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"awesomeProject/readxls"

	excellizev2 "github.com/xuri/excelize/v2"
)

const aggregationSheetName = "汇总"
const firstRowAxis = "A1"
const uselessSheet = "Sheet1"

type ExtendedExcel struct {
	f             *excellizev2.File
	rowIndex      int64
	hasTitleRow   bool
	excelsDirPath string
	workers       chan struct{}
	workDone      chan string
	workCount     chan int
}

func NewExtendExcel(excelsDirPath string, hasTitleRow bool, workers int) *ExtendedExcel {
	if workers <= 0 {
		panic("invalid workers count")
	}

	f := excellizev2.NewFile()
	fmt.Println(f.GetSheetList())
	f.NewSheet(aggregationSheetName)
	f.DeleteSheet(uselessSheet)

	return &ExtendedExcel{
		f:             f,
		workers:       make(chan struct{}, workers),
		workDone:      make(chan string, 0),
		workCount:     make(chan int, 0),
		hasTitleRow:   hasTitleRow,
		excelsDirPath: excelsDirPath,
	}
}

func (ee *ExtendedExcel) addRows(excelName, sheetName string, sheetRows [][]string) {
	defer func() {
		fmt.Printf("worker finish excel %s sheet %s\n", excelName, sheetName)
		<-ee.workers
		ee.workDone <- fmt.Sprintf("excel: %s sheet %s", excelName, sheetName)
	}()

	fmt.Printf("begin add excel %s sheet %s\n", excelName, sheetName)

	rowsCount := len(sheetRows)
	if rowsCount <= 0 {
		return
	}

	sheetStartRowIndex := 0
	if ee.hasTitleRow {
		sheetStartRowIndex = 1

		axis := fmt.Sprintf("A%d", atomic.AddInt64(&ee.rowIndex, 1))
		if axis == firstRowAxis {
			err := ee.f.SetSheetRow(aggregationSheetName, axis, &sheetRows[0])
			if err != nil {
				panic(err)
			}
		}
	} else {
		sheetStartRowIndex = 0
	}

	for i := sheetStartRowIndex; i < rowsCount; i++ {
		axis := fmt.Sprintf("A%d", atomic.AddInt64(&ee.rowIndex, 1))
		err := ee.f.SetSheetRow(aggregationSheetName, axis, &sheetRows[i])
		if err != nil {
			panic(err)
		}
	}

}

func (ee *ExtendedExcel) AddRowsWorker(excelName, sheetName string, sheetRows [][]string) {
	ee.workers <- struct{}{}
	go ee.addRows(excelName, sheetName, sheetRows)
}

func (ee *ExtendedExcel) WaitAllExcelsJoined() {
	totalCount := -1
	count := 0
	for {
		select {
		case taskName := <-ee.workDone:
			count++
			fmt.Printf("%s added\n", taskName)
		case totalCount = <-ee.workCount:
			fmt.Printf("total sheet count %d\n", totalCount)
		default:
			time.Sleep(time.Second)
		}
		if totalCount >= 0 && count >= totalCount {
			fmt.Printf("all works done\n")
			break
		}
	}

	ee.CleanBlankRow()
}

func (ee *ExtendedExcel) joinXlsx(excelsDir string, entry os.DirEntry) int {
	excelEntryF, excelEntryErr := excellizev2.OpenFile(path.Join(excelsDir, entry.Name()))
	if excelEntryErr != nil {
		panic(excelEntryErr)
	}

	sheets := excelEntryF.GetSheetList()
	for _, sheetName := range sheets {
		sheetRows, err := excelEntryF.GetRows(sheetName)
		if err != nil {
			panic(err)
		}
		ee.AddRowsWorker(entry.Name(), sheetName, sheetRows)
	}

	return len(sheets)
}

func (ee *ExtendedExcel) joinXls(excelsDir string, entry os.DirEntry) int {
	sheetsMap := readxls.ReadXls(path.Join(excelsDir, entry.Name()))
	for sheetName, sheetRows := range sheetsMap {
		ee.AddRowsWorker(entry.Name(), sheetName, sheetRows)
	}

	return len(sheetsMap)
}

func (ee *ExtendedExcel) joinExcel(excelsDir string, entry os.DirEntry) int {
	switch {
	case strings.HasSuffix(entry.Name(), "xls"):
		return ee.joinXls(excelsDir, entry)
	case strings.HasSuffix(entry.Name(), "xlsx"):
		return ee.joinXlsx(excelsDir, entry)
	}
	return 0
}

func (ee *ExtendedExcel) joinExcels() {
	totalWorkCount := 0
	dirF, dirErr := os.OpenFile(ee.excelsDirPath, os.O_RDONLY, 0777)
	if dirErr != nil {
		panic(dirErr)
	}

	dirEntries, readDirErr := dirF.ReadDir(0)
	if readDirErr != nil {
		panic(readDirErr)
	}

	for _, entry := range dirEntries {
		if !entry.Type().IsRegular() {
			continue
		}

		totalWorkCount += ee.joinExcel(ee.excelsDirPath, entry)
	}

	ee.workCount <- totalWorkCount
}

func (ee *ExtendedExcel) JoinExcelsInBackground() {
	go ee.joinExcels()
}

func (ee *ExtendedExcel) CleanBlankRow() {
	blankRows := make([]int, 0)
	sheetValues, err := ee.f.GetRows(aggregationSheetName)
	if err != nil {
		panic(err)
	}
	for index, rowValues := range sheetValues {
		if len(rowValues) == 0 {
			blankRows = append(blankRows, index+1)
			continue
		}
		hasValue := false
		for _, value := range rowValues {
			hasValue = len(value) != 0
		}
		if !hasValue {
			blankRows = append(blankRows, index+1)
		}
	}

	for i := len(blankRows) - 1; i >= 0; i-- {
		err := ee.f.RemoveRow(aggregationSheetName, blankRows[i])
		if err != nil {
			panic(err)
		}
	}
}

func (ee *ExtendedExcel) SaveToFile(filePath string) {
	if err := ee.f.SaveAs(path.Clean(filePath)); err != nil {
		panic(err)
	}
}

func help() {
	fmt.Printf("%s D:\\abc\\小组统计汇总目录  true,false有无标题  D:\\输出汇总文件.xlsx\n", os.Args[0])
}

func main() {
	if len(os.Args) != 4 {
		help()
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}

	excelsDirPath := os.Args[1]
	hasExcelTitleRow, boolErr := strconv.ParseBool(os.Args[2])
	if boolErr != nil {
		panic(boolErr)
	}
	mergedExcelFilePath := os.Args[3]

	ee := NewExtendExcel(excelsDirPath, hasExcelTitleRow, runtime.NumCPU())
	ee.JoinExcelsInBackground()
	ee.WaitAllExcelsJoined()

	ee.SaveToFile(mergedExcelFilePath)
}
