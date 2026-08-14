package utils

import (
	"time"
)

func CalcDate() int {
	flagDay := "08-01"
	// 解析指定日期
	location, _ := time.LoadLocation("Local")
	specifiedDate, _ := time.ParseInLocation("01-02", flagDay, location)

	// 获取当前日期
	currentDate := time.Now().In(location)

	// 提取年、月、日信息
	currentYear, currentMonth, currentDay := currentDate.Date()
	_, specifiedMonth, specifiedDay := specifiedDate.Date()

	if currentMonth < specifiedMonth {
		return currentYear - 3
	} else if currentMonth > specifiedMonth {
		return currentYear - 2
	}

	if currentDay >= specifiedDay {
		return currentYear - 2
	}

	return 0
}
