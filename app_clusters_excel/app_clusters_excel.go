package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	excellizev2 "github.com/xuri/excelize/v2"
)

func appAddClusters(appTemplateClustersMap map[string]map[string]struct{}, app, clusters string) {
	fmt.Printf("app: %s clusters %s\n", app, clusters)
	appClusterMap, found := appTemplateClustersMap[app]
	if !found {
		appClusterMap = make(map[string]struct{})
		appTemplateClustersMap[app] = appClusterMap
	}

	for _, cluster := range strings.Split(clusters, ",") {
		if cluster == "" {
			continue
		}
		appClusterMap[cluster] = struct{}{}
	}
}

func getAppClustersFromXlsx(appTemplateClustersMap map[string]map[string]struct{},
	excelFilePath string, isTemplate bool) {
	f, err := excellizev2.OpenFile(excelFilePath)
	if err != nil {
		panic(err)
	}

	sheets := f.GetSheetList()
	for _, sheet := range sheets {
		fmt.Printf("------------开始分析sheet '%s'\n", sheet)
		rows, getRowsErr := f.GetRows(sheet)
		if getRowsErr != nil {
			panic(getRowsErr)
		}
		getAppTemplateClustersFromRows(appTemplateClustersMap, rows, isTemplate)
		fmt.Printf("------------完成分析sheet '%s'\n", sheet)
	}
}

func getAppClustersFromCsv(appTemplateClustersMap map[string]map[string]struct{},
	excelFilePath string, isTemplate bool) {
	csvF, csvErr := os.OpenFile(excelFilePath, os.O_RDONLY, 0777)
	if csvErr != nil {
		panic(csvErr)
	}
	r := csv.NewReader(bufio.NewReader(csvF))
	r.LazyQuotes = true
	rows, csvReadErr := r.ReadAll()
	if csvReadErr != nil {
		panic(csvReadErr)
	}

	getAppTemplateClustersFromRows(appTemplateClustersMap, rows, isTemplate)
}

func getAppTemplateClustersFromRows(appClustersMap map[string]map[string]struct{}, rows [][]string, isTemplate bool) {
	var appNameIndex = -1
	var clusterIndex = -1

	if len(rows) == 0 {
		return
	}

	fmt.Printf("rows count %d\n", len(rows))
	if len(rows) > 1 {
		fmt.Printf("%+v\n", rows[0])
		fmt.Printf("%+v\n", rows[1])
	}

	clusterColTitle := "DEPLOY_CLUSTER"
	if !isTemplate {
		clusterColTitle = "CLUSTER_IDS"
	}

	for index, col := range rows[0] {
		if col == "APP_NAME" {
			appNameIndex = index
		}
		if col == clusterColTitle {
			clusterIndex = index
		}
		if clusterIndex != -1 && appNameIndex != -1 {
			break
		}
	}
	if clusterIndex == -1 || appNameIndex == -1 {
		return
	}
	fmt.Printf("appNameIndex: %d clusterIndex %d\n", appNameIndex, clusterIndex)
	for index, row := range rows {
		if index == 0 {
			continue
		}
		appAddClusters(appClustersMap, row[appNameIndex], row[clusterIndex])
	}
}

func getAppClustersFromExcel(appClustersMap map[string]map[string]struct{}, excelFilePath string, isTemplate bool) {
	switch {
	case strings.HasSuffix(excelFilePath, ".xlsx"):
		getAppClustersFromXlsx(appClustersMap, excelFilePath, isTemplate)
	case strings.HasSuffix(excelFilePath, ".csv"):
		getAppClustersFromCsv(appClustersMap, excelFilePath, isTemplate)
	}
}

func help() {
	fmt.Printf("Usage:\t%s templatesDir appDir outputDir\n", os.Args[0])
	fmt.Printf("Example: %s D:\\appUselessClusterFilter\\templates D:\\appUselessClusterFilter\\apps D:\\appUselessClusterFilter\\output\n", os.Args[0])
	fmt.Println("Important: if you are using csv excel, dos2unix them first !!!")
	os.Exit(1)
}

func walkDir(appClustersMap map[string]map[string]struct{}, excelDir string, isTemplate bool) {
	if isTemplate {
		fmt.Println("开始分析模板")
		defer fmt.Println("完成分析模板")
	} else {
		fmt.Println("开始分析应用")
		defer fmt.Println("完成分析应用")
	}

	_ = filepath.Walk(excelDir, func(itemPath string, itemInfo fs.FileInfo, itemErr error) error {
		if itemErr != nil {
			return nil
		}
		if itemInfo.IsDir() {
			return nil
		}
		if !strings.HasSuffix(itemPath, ".xlsx") && !strings.HasSuffix(itemPath, ".csv") {
			return nil
		}
		fmt.Printf("------开始分析excel '%s'\n", itemPath)
		getAppClustersFromExcel(appClustersMap, itemPath, isTemplate)
		fmt.Printf("------结束分析excel '%s'\n", itemPath)
		return nil
	})
}

const summary = "汇总"

func clustersSummary(appTemplateClustersMap, appClustersMap map[string]map[string]struct{}, outputDir string) {
	fmt.Println("开始分析差异")
	defer fmt.Println("完成分析差异")

	outputExcel := excellizev2.NewFile()
	outputExcel.NewSheet(summary)
	outputExcel.DeleteSheet("Sheet1")
	err := outputExcel.SetSheetRow(summary, "A1", &[]string{
		"APP_NAME",
		"APP_BIND_CLUSTER_COUNT",
		"APP_BIND_CLUSTER",
		"TEMPLATE_USED_CLUSTER_COUNT",
		"TEMPLATE_USED_CLUSTER",
		"NOT_USED_CLUSTER_COUNT",
		"NOT_USED_CLUSTER"})
	if err != nil {
		panic(err)
	}

	var apps = make([]string, 0, len(appClustersMap))

	for app, _ := range appClustersMap {
		apps = append(apps, app)
	}
	sort.Strings(apps)

	contentRowNum := 2
	for _, appItem := range apps {
		appBindClustersMap, _ := appClustersMap[appItem]
		appBindClustersArr := make([]string, 0, len(appBindClustersMap))
		for app, _ := range appBindClustersMap {
			appBindClustersArr = append(appBindClustersArr, app)
		}
		sort.Strings(appBindClustersArr)

		oneAppTemplateClustersMap, foundOneAppTemplateClustersMap := appTemplateClustersMap[appItem]
		appTemplateUsedClustersArr := make([]string, 0, len(oneAppTemplateClustersMap))
		if foundOneAppTemplateClustersMap {
			for app, _ := range oneAppTemplateClustersMap {
				appTemplateUsedClustersArr = append(appTemplateUsedClustersArr, app)
			}
		}
		sort.Strings(appTemplateUsedClustersArr)

		notUsedClusters := make([]string, 0)
		for cluster, _ := range appBindClustersMap {
			if !foundOneAppTemplateClustersMap {
				notUsedClusters = append(notUsedClusters, cluster)
			} else {
				_, found := oneAppTemplateClustersMap[cluster]
				if !found {
					notUsedClusters = append(notUsedClusters, cluster)
				}
			}
		}
		setRowErr := outputExcel.SetSheetRow(summary, fmt.Sprintf("A%d", contentRowNum), &[]interface{}{
			appItem,
			len(appBindClustersArr),
			strings.Join(appBindClustersArr, ","),
			len(appTemplateUsedClustersArr),
			strings.Join(appTemplateUsedClustersArr, ","),
			len(notUsedClusters),
			strings.Join(notUsedClusters, ",")})
		if setRowErr != nil {
			panic(setRowErr)
		}
		contentRowNum++
	}

	saveErr := outputExcel.SaveAs(path.Join(outputDir, "汇总.xlsx"))
	if saveErr != nil {
		panic(saveErr)
	}
	fmt.Printf("输出分析excel %s\n", path.Clean(strings.Join([]string{outputDir, "汇总.xlsx"}, "\\")))
}

func main() {
	if len(os.Args) != 4 {
		help()
	}
	templateDir := os.Args[1]
	appDir := os.Args[2]
	outputDir := os.Args[3]
	appTemplateClustersMap := make(map[string]map[string]struct{})
	appClustersMap := make(map[string]map[string]struct{})

	walkDir(appTemplateClustersMap, templateDir, true)
	walkDir(appClustersMap, appDir, false)

	clustersSummary(appTemplateClustersMap, appClustersMap, outputDir)
}
