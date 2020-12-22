package main

import (
	"github.com/modeltool/generate"
)

func main() {
	//TODO 加入flag命令行执行

	//generate.Generate() //生成所有表信息
	generate.Generate(
		"t_swap_order")
}
