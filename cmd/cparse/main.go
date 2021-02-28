package main

import (
	"discord_autosender/pkg/curlparser"
	"fmt"
)

func main() {
	cp := curlparser.CurlParser{}

	cp.GetConf()

	fmt.Println(cp.GetHeader("authority"))
}
