package generate

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/generator"
	"github.com/modeltool/conf"
)

// Generate 生成model
func Generate(tableNames ...string) {
	tableNamesStr := ""
	for _, name := range tableNames {
		if tableNamesStr != "" {
			tableNamesStr += ","
		}
		tableNamesStr += "'" + name + "'"
	}
	tables := getTables(tableNamesStr) //生成所有表信息
	//tables := getTables("admin_info","video_info") //生成指定表信息，可变参数可传入过个表名
	for _, table := range tables {
		fields := getFields(table.Name)
		generateModel(table, fields)
	}
}

//获取表信息
func getTables(tableNames string) []Table {
	db := SharedStore().db
	var tables []Table
	if tableNames == "" {
		db.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE table_schema='" + conf.GetConfig().Database + "';").Find(&tables)
	} else {
		db.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE TABLE_NAME IN (" + tableNames + ") AND table_schema='" + conf.GetConfig().Database + "';").Find(&tables)
	}
	return tables
}

//获取所有字段信息
func getFields(tableName string) []Field {
	db := SharedStore().db
	var fields []Field
	db.Raw("show FULL COLUMNS from " + tableName + ";").Find(&fields)
	return fields
}

//生成Model
func generateModel(table Table, fields []Field) {
	content := "package model\n\n"
	//表注释
	if len(table.Comment) > 0 {
		content += "// " + table.Comment + "\n"
	}
	content += "type " + generator.CamelCase(table.Name) + " struct {\n"
	//生成字段
	for _, field := range fields {
		fieldName := generator.CamelCase(field.Field)
		fieldType := getFiledType(field)
		fieldFormatType := getFieldFormatType(field)
		fieldComment := getFieldComment(field)
		content += "	" + fieldName + " " + fieldType + " " + fieldFormatType + " " + fieldComment + "\n"
	}
	content += "}"

	filename := conf.GetConfig().ModelPath + table.Name + ".go"
	var f *os.File
	var err error
	if checkFileIsExist(filename) {
		if !conf.GetConfig().ModelReplace {
			fmt.Println(generator.CamelCase(table.Name) + " 已存在，需删除才能重新生成...")
			return
		}
		f, err = os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0666) //打开文件
		if err != nil {
			panic(err)
		}
	} else {
		f, err = os.Create(filename)
		if err != nil {
			panic(err)
		}
	}
	defer f.Close()
	_, err = io.WriteString(f, content)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(generator.CamelCase(table.Name) + " 已生成...")
	}
}

//获取字段类型
func getFiledType(field Field) string {
	typeArr := strings.Split(field.Type, "(")

	switch typeArr[0] {
	case "int":
		return "int"
	case "int32":
		return "int32"
	case "int64":
		return "int64"
	case "integer":
		return "int"
	case "mediumint":
		return "int"
	case "bit":
		return "int"
	case "year":
		return "int"
	case "smallint":
		return "int"
	case "tinyint":
		return "int"
	case "bigint":
		return "int64"
	case "decimal":
		return "float64"
	case "double":
		return "float64"
	case "float":
		return "float32"
	case "real":
		return "float32"
	case "numeric":
		return "float32"
	case "timestamp":
		return "string"
	case "datetime":
		return "time.Time"
	case "time":
		return "time.Time"
	default:
		return "string"
	}
}

//获取字段json描述
func getFieldJSON(field Field) string {
	return `json:"` + field.Field + `"`
}

func getFieldOrm(field Field) string {
	return `orm:"` + field.Field + `"`
}

func getFieldGconv(field Field) string {
	return `gconv:"` + field.Field + `"`
}

func getFieldFormatType(field Field) string {
	var fieldFormatType string
	for i, v := range conf.GetConfig().FieldFormatTypes {
		if i == len(conf.GetConfig().FieldFormatTypes)-1 {
			fieldFormatType += fmt.Sprintf("%v:\"%v\"", v, field.Field)
		} else {
			fieldFormatType += fmt.Sprintf("%v:\"%v\" ", v, field.Field)
		}
	}
	return fmt.Sprintf("`%v`", fieldFormatType)
}

//获取字段说明
func getFieldComment(field Field) string {
	if len(field.Comment) > 0 {
		return "// " + field.Comment
	}
	return ""
}

//检查文件是否存在
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
