package utils

import (
	"fmt"
	"io"
	"reflect"

	"github.com/xuri/excelize/v2"
)

var HEADER = []string{"年级", "班级", "姓名", "联系方式", "学科", "考试名称", "得分"}

func ReadExcel(excelFile io.Reader) ([][]string, []string) {
	var allRows [][]string
	var respErr []string

	f, err := excelize.OpenReader(excelFile)
	if err != nil {
		// 读取文件错误
		respErr = append(respErr, "读取文件错误")
		return allRows, respErr
	}

	// 全部行数据
	// 获取全部工作表
	allSheets := f.GetSheetList()
	for _, sheet := range allSheets {
		rows, err := f.GetRows(sheet)
		if err != nil {
			// 打开工作表失败
			respErr = append(respErr, fmt.Sprintf("打开工作表 %s 失败", sheet))
			continue
		}
		// 检查表头
		if !reflect.DeepEqual(rows[0][:7], HEADER) {
			respErr = append(respErr, fmt.Sprintf("工作表 %s 表头不正确", sheet))
			continue
		}

		allRows = append(allRows, rows[1:]...)
	}

	return allRows, respErr
}
