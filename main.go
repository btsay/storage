package main

import (
	"github.com/btsay/storage/crawl"
	"github.com/btsay/storage/utils"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	utils.Init()
	crawl.Run()
}
