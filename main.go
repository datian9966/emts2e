package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)
import _ "github.com/go-sql-driver/mysql"

type Structure struct {
	columnName             sql.NullString
	DataType               sql.NullString
	CharacterMaximumLength sql.NullString
	IsNullable             sql.NullString
	ColumnDefault          sql.NullString
	ColumnComment          sql.NullString
}

func main() {

	var (
		url      string
		dbName   string
		userName string
		password string
		fileName string
	)

	//获取参数
	u := flag.String("url", "", "数据库url")
	d := flag.String("db", "", "库名")
	un := flag.String("u", "", "用户名")
	p := flag.String("p", "", "密码")
	f := flag.String("f", "", "导出文件名")

	flag.Parse()

	url = *u
	dbName = *d
	userName = *un
	password = *p
	fileName = *f

	if url == "" || dbName == "" || userName == "" || password == "" {
		//如果flag模式的参数没有输入全 则启动交互模式
		fmt.Println("请输入数据库url: (必填 例:localhost:3306)")
		fmt.Scanln(&url)
		fmt.Println("请输入数据库名称: (必填)")
		fmt.Scanln(&dbName)
		fmt.Println("请输入数据库用户名: (必填)")
		fmt.Scanln(&userName)
		fmt.Println("请输入数据库密码: (必填)")
		fmt.Scanln(&password)
		fmt.Println("请输入导出文件名称: (选填 默认 structure.xlsx 如自定义名字 .xlsx后缀也要写全)")
		fmt.Scanln(&fileName)
	}

	if url == "" || dbName == "" || userName == "" || password == "" {
		fmt.Println("参数不全,无法启动程序")
		return
	}

	if fileName == "" {
		fileName = "structure.xlsx"
	}

	fmt.Println("参数校验通过,程序开始开始执行")

	db, _ := sql.Open("mysql", ""+userName+":"+password+"@tcp("+url+")/"+dbName+"?charset=utf8")

	xlsx := excelize.NewFile()
	xlsx.SetActiveSheet(0)
	style, _ := xlsx.NewStyle(`{"border":[{"type":"left","color":" 000000","style":1},{"type":"top","color":"000000","style":1},{"type":"bottom","color":"000000","style":1},{"type":"right","color":"000000","style":1}]}`)
	hyperLinkStyle, _ := xlsx.NewStyle(`{"border":[{"type":"left","color":" 000000","style":1},{"type":"top","color":"000000","style":1},{"type":"bottom","color":"000000","style":1},{"type":"right","color":"000000","style":1}],"font":{"color":"#1265BE","underline":"single"}}`)
	xlsx.SetSheetName("Sheet1", "概览")
	xlsx.SetCellValue("概览", "A1", "序号")
	xlsx.SetCellStyle("概览", "A1", "A1", style)
	xlsx.SetCellValue("概览", "B1", "备注")
	xlsx.SetCellStyle("概览", "B1", "B1", style)
	xlsx.SetCellValue("概览", "C1", "表名")
	xlsx.SetCellStyle("概览", "C1", "C1", style)
	xlsx.SetColWidth("概览", "B", "C", 50)
	tableRows, _ := db.Query("SELECT TABLE_NAME,TABLE_COMMENT FROM information_schema.TABLES WHERE table_schema = '" + dbName + "'")
	var count float64
	db.QueryRow("select count(*) FROM information_schema.TABLES WHERE table_schema = '" + dbName + "'").Scan(&count)
	tableRowIndex := 2
	defer tableRows.Close()
	for tableRows.Next() {
		var TableName string
		var TableComment string
		tableRows.Scan(&TableName, &TableComment)
		rows, _ := db.Query("SELECT COLUMN_NAME,DATA_TYPE,CHARACTER_MAXIMUM_LENGTH,IS_NULLABLE,COLUMN_DEFAULT,COLUMN_COMMENT FROM INFORMATION_SCHEMA.COLUMNS WHERE table_schema = '" + dbName + "' AND table_name = '" + TableName + "'")
		defer rows.Close()
		xlsx.NewSheet(TableName)
		xlsx.SetColWidth(TableName, "A", "A", 30)
		xlsx.SetColWidth(TableName, "F", "F", 80)
		xlsx.SetCellValue(TableName, "A1", "列名")
		xlsx.SetCellStyle(TableName, "A1", "A1", style)
		xlsx.SetCellValue(TableName, "B1", "类型")
		xlsx.SetCellStyle(TableName, "B1", "B1", style)
		xlsx.SetCellValue(TableName, "C1", "长度")
		xlsx.SetCellStyle(TableName, "C1", "C1", style)
		xlsx.SetCellValue(TableName, "D1", "是否为空")
		xlsx.SetCellStyle(TableName, "D1", "D1", style)
		xlsx.SetCellValue(TableName, "E1", "默认值")
		xlsx.SetCellStyle(TableName, "E1", "E1", style)
		xlsx.SetCellValue(TableName, "F1", "备注")
		xlsx.SetCellStyle(TableName, "F1", "F1", style)
		rowIndex := 2
		for rows.Next() {
			rowIndexS := strconv.Itoa(rowIndex)
			s := Structure{}
			rows.Scan(&s.columnName, &s.DataType, &s.CharacterMaximumLength, &s.IsNullable, &s.ColumnDefault, &s.ColumnComment)
			xlsx.SetCellValue(TableName, "A"+rowIndexS, s.columnName.String)
			xlsx.SetCellStyle(TableName, "A"+rowIndexS, "A"+rowIndexS, style)
			xlsx.SetCellValue(TableName, "B"+rowIndexS, s.DataType.String)
			xlsx.SetCellStyle(TableName, "B"+rowIndexS, "B"+rowIndexS, style)
			xlsx.SetCellValue(TableName, "C"+rowIndexS, s.CharacterMaximumLength.String)
			xlsx.SetCellStyle(TableName, "C"+rowIndexS, "C"+rowIndexS, style)
			xlsx.SetCellValue(TableName, "D"+rowIndexS, s.IsNullable.String)
			xlsx.SetCellStyle(TableName, "D"+rowIndexS, "D"+rowIndexS, style)
			xlsx.SetCellValue(TableName, "E"+rowIndexS, s.ColumnDefault.String)
			xlsx.SetCellStyle(TableName, "E"+rowIndexS, "E"+rowIndexS, style)
			xlsx.SetCellValue(TableName, "F"+rowIndexS, s.ColumnComment.String)
			xlsx.SetCellStyle(TableName, "F"+rowIndexS, "F"+rowIndexS, style)
			rowIndex++
		}
		tableRowIndexS := strconv.Itoa(tableRowIndex)
		xlsx.SetCellValue("概览", "A"+tableRowIndexS, tableRowIndex-1)
		xlsx.SetCellStyle("概览", "A"+tableRowIndexS, "A"+tableRowIndexS, style)
		xlsx.SetCellValue("概览", "B"+tableRowIndexS, TableComment)
		xlsx.SetCellStyle("概览", "B"+tableRowIndexS, "B"+tableRowIndexS, style)
		xlsx.SetCellValue("概览", "C"+tableRowIndexS, TableName)
		xlsx.SetCellStyle("概览", "C"+tableRowIndexS, "C"+tableRowIndexS, hyperLinkStyle)
		xlsx.SetCellHyperLink("概览", "C"+tableRowIndexS, TableName+"!A1", "Location")
		fmt.Println("当前已完成" + strconv.Itoa(int((float64(tableRowIndex-1)/count)*100)) + "%...")
		tableRowIndex++
	}
	err := xlsx.SaveAs(getCurrentDirectory() + "/" + fileName)
	fmt.Println("表结构导出完毕完毕,再见")
	if err != nil {
		fmt.Println(err)
	}
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
