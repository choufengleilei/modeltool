package main

import (
	"flag"
	"github.com/modeltool/generate"
	"strings"
)

func main() {
	tables := flag.String("table", "", "mysql tableNames example: table_name1,table_name2")
	flag.Parse()
	var ts []string
	if tables != nil {
		ts = strings.Split(*tables, ",")
	}
	//generate.Generate() //生成所有表信息
	generate.Generate(ts...)
}
