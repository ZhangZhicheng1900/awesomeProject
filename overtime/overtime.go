package main

import (
	"awesomeProject/xlsxop"
	"fmt"
	"os"
	"sort"
	"time"
)

const (
	InOutSheet            = "出入记录查询"
	LineIndexBegin        = 2
	FlagOut               = "出"
	StandOverTimeBeginStr = "18:00:00"
)

const (
	ColIndexName = iota
	ColIndexCardNo
	ColIndexDate
	ColIndexTime
	ColIndexLocation
	ColIndexInOutFlag
)

func parseTime(dateStr, timeStr string) time.Time {
	t, err := time.Parse("2006010215:04:05Z07:00", fmt.Sprintf("%s%s+08:00", dateStr, timeStr))
	if err != nil {
		panic(err)
	}
	return t
}

type OutTime struct {
	time.Time
	Date               string
	standOverTimeBegin time.Time
}

func (ot OutTime) OverTimeInMinute() int {
	if ot.Before(ot.standOverTimeBegin) {
		return 0
	}
	return int(ot.Sub(ot.standOverTimeBegin) / time.Minute)
}

func (ot OutTime) OverTimeInHalfHour() int {
	if ot.Before(ot.standOverTimeBegin) {
		return 0
	}
	return int(ot.Sub(ot.standOverTimeBegin) / (30 * time.Minute))
}

func (ot OutTime) IsOverTime() bool {
	return ot.After(ot.standOverTimeBegin)
}

func NewOutTime(date, t string) OutTime {
	outTime := parseTime(date, t)
	standOverTimeBegin := parseTime(date, StandOverTimeBeginStr)
	return OutTime{Time: outTime, Date: date, standOverTimeBegin: standOverTimeBegin}
}

func findLatestOut(sheet [][]string) map[string]OutTime {
	var latest = make(map[string]OutTime)

	for i := LineIndexBegin; i < len(sheet); i++ {
		line := sheet[i]
		if line[ColIndexInOutFlag] != FlagOut {
			continue
		}
		outTime := NewOutTime(line[ColIndexDate], line[ColIndexTime])
		latestTime, found := latest[line[ColIndexDate]]
		if found && outTime.Before(latestTime.Time) {
			continue
		}
		latest[line[ColIndexDate]] = outTime
	}

	return latest
}

func main() {
	var filePath = "出入记录.xlsx"
	if len(os.Args) >= 2 {
		filePath = os.Args[1]
	}
	data := xlsxop.Read(filePath)
	sheet, found := data[InOutSheet]
	if !found {
		panic(fmt.Errorf("sheet %s not found\n", InOutSheet))
	}

	var minuteCount, halfHourCount int
	latest := findLatestOut(sheet)

	var dates = make([]string, len(latest))
	i := 0
	for k, _ := range latest {
		dates[i] = k
		i++
	}
	sort.Strings(dates)

	for _, date := range dates {
		t := latest[date]
		if !t.IsOverTime() {
			continue
		}
		minutes := t.OverTimeInMinute()
		halfHours := t.OverTimeInHalfHour()
		fmt.Printf("%s overtime %d in minute, %.1f in hour(count by half hour)\n", date, minutes, float64(halfHours)*0.5)
		minuteCount += minutes
		halfHourCount += halfHours
	}

	fmt.Printf("total overtime %d in minute, %.1f in hour(count by half hour)\n", minuteCount, float64(halfHourCount)*0.5)

	fmt.Println("Press enter key to quit...")
	fmt.Scanln()
}
