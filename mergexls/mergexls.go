package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"awesomeProject/xlsxop"

	excellizev2 "github.com/xuri/excelize/v2"
)

var versionStr = "202204261128"

var pythonXls2XlsxScript = `
import pyexcel as p
import sys
xlsFile = sys.argv[1]
xlsxFile = sys.argv[2]
p.save_book_as(file_name=xlsFile,
               library='pyexcel-xls',
               skip_hidden_row_and_column=False,
               dest_file_name=xlsxFile)
`

const (
	xls2xlsx = "xls2xlsx.py"
)

func Fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func WritePythonScript() {
	err := os.WriteFile(path.Join(tempDir, xls2xlsx), []byte(pythonXls2XlsxScript), 0777)
	if err != nil {
		Fatal("create xls2xlsx.py error, %v", err)
	}
}

func ExcelTranslate(pythonFile, srcPath, newPath string) error {
	// 删除已经存在同名xlsx文件
	if _, err := os.Stat(newPath); err == nil {
		os.Remove(newPath)
	}
	// 执行命令（阻塞式调用）
	err := exec.Command(*pythonCMD, path.Join(tempDir, pythonFile), srcPath, newPath).Run()
	if err != nil {
		return err
	}
	// 查看xlsx文件是否成功生成
	if _, err := os.Stat(newPath); err != nil {
		return fmt.Errorf("using %s translate %s failed", pythonFile, srcPath)
	}
	return nil
}

func getEntries(dirPath string) []os.DirEntry {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		Fatal("read dir error, %v", err)
	}
	return entries
}

func TranslateDirXls2Xlsx(dirPath *string) {
	entries := getEntries(*dirPath)
	for _, item := range entries {
		if !item.Type().IsRegular() || !strings.HasSuffix(item.Name(), ".xls") {
			continue
		}
		filePath := path.Join(*dirPath, item.Name())
		err := ExcelTranslate(xls2xlsx, filePath, path.Join(tempDir, item.Name()+"x"))
		if err != nil {
			Fatal("%v", err)
		}
	}
}

func MergeXlsxSheets() map[string][][]string {
	var merged = make(map[string][][]string)
	var sheetsIndex = make(map[string]string)
	entries := getEntries(tempDir)
	for _, item := range entries {
		if !item.Type().IsRegular() || !strings.HasSuffix(item.Name(), ".xlsx") {
			continue
		}
		filePath := path.Join(tempDir, item.Name())
		xlsxContent := xlsxop.Read(filePath)

		for name, sheet := range xlsxContent {
			hisXlsxFile, found := sheetsIndex[name]
			if found {
				Fatal("merging %s: sheet %s already exist in %s", filePath, name, hisXlsxFile)
			}
			merged[name] = sheet
			sheetsIndex[name] = filePath
		}
	}
	return merged
}

func GenerateMerged(xlsxFilePath string, content map[string][][]string) {
	styleFunc := func(xlsxFile *excellizev2.File) error {
		styleID, err := xlsxFile.NewStyle(&excellizev2.Style{
			Alignment: &excellizev2.Alignment{
				WrapText:   true,
				Horizontal: "left",
				Vertical:   "top",
			},
		})
		if err != nil {
			return fmt.Errorf("NewStyle error, %v", err)
		}

		sheets := xlsxFile.GetSheetList()
		for _, sheetName := range sheets {
			if err := xlsxFile.SetRowHeight(sheetName, 2, 100); err != nil {
				return fmt.Errorf("SetRowHeight error, %v", err)
			}
			if err := xlsxFile.SetColStyle(sheetName, "A:D", styleID); err != nil {
				return fmt.Errorf("SetColStyle error, %v", err)
			}
			if err := xlsxFile.SetColWidth(sheetName, "A", "D", 60); err != nil {
				return fmt.Errorf("SetColWidth error, %v", err)
			}
		}
		return nil
	}

	if err := xlsxop.Write(xlsxFilePath, content, styleFunc); err != nil {
		Fatal("generate merged %s error, %v", xlsxFilePath, err)
	}
}

var pythonCMD *string

var tempDir = "tmp_merge"

func main() {
	dirPath := flag.String("to-merge-xls-dir", "./", "the dir where xls files need to be merged")
	oFilePath := flag.String("output-file-name", "merged.xlsx", "output xlsx file name, will created "+
		"at 'to-merge-xls-dir'/tmp_merge")
	pythonCMD = flag.String("python-cmd", "python", "python cmd, maybe 'python' or 'python3'")
	version := flag.Bool("version", false, "version of this tool")
	flag.Parse()

	if *version {
		fmt.Println(versionStr)
		return
	}

	tempDir = path.Join(*dirPath, "tmp_merge")
	os.RemoveAll(tempDir)
	if err := os.Mkdir(tempDir, 0777); err != nil {
		Fatal("mkdir %s error, %v", path.Join(*dirPath, "tmp"), err)
	}

	WritePythonScript()
	TranslateDirXls2Xlsx(dirPath)
	GenerateMerged(path.Join(tempDir, *oFilePath), MergeXlsxSheets())
}
