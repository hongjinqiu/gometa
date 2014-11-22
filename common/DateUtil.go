package common

import (
	"fmt"
	"strconv"
	"time"
)

type DateUtil struct{}

func (o DateUtil) GetCurrentYyyyMMdd() int {
	date, err := strconv.Atoi(time.Now().Format("20060102"))
	if err != nil {
		panic(err)
	}
	return date
}

func (o DateUtil) GetDateByFormat(format string) string {
	return time.Now().Format(format)
}

func (o DateUtil) GetCurrentYyyyMMddHHmmss() int64 {
	createTime, err := strconv.ParseInt(time.Now().Format("20060102150405"), 10, 64)
	if err != nil {
		panic(err)
	}
	return createTime
}

func (o DateUtil) ConvertDate2String(date string, sourcePattern string, targetPattern string) string {
	dDate, err := time.Parse(sourcePattern, fmt.Sprint(date))
	if err != nil {
		panic(err)
	}
	return dDate.Format(targetPattern)
}

func (o DateUtil) GetNextDate(date int) int {
	dDate, err := time.Parse("20060102", fmt.Sprint(date))
	if err != nil {
		panic(err)
	}
	dDate = dDate.AddDate(0, 0, 1)
	result, err := strconv.Atoi(dDate.Format("20060102"))
	if err != nil {
		panic(err)
	}
	return result
}

func (o DateUtil) GetPreDate(date int) int {
	dDate, err := time.Parse("20060102", fmt.Sprint(date))
	if err != nil {
		panic(err)
	}
	dDate = dDate.AddDate(0, 0, -1)
	result, err := strconv.Atoi(dDate.Format("20060102"))
	if err != nil {
		panic(err)
	}
	return result
}
