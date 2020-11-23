package main

import (
	"github.com/modeltool/generate"
)

func main() {
	//generate.Genertate() //生成所有表信息
	generate.Genertate(
		"t_wallet_token_address")
}
