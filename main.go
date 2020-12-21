package main

import (
	"github.com/modeltool/generate"
)

func main() {
	//generate.Generate() //生成所有表信息
	generate.Generate(
		"t_swap_order")
}
