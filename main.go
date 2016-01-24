package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/nzai/regimentation/config"
	"github.com/nzai/regimentation/data"
)

func main() {

	defer func() {
		// 捕获panic异常
		if err := recover(); err != nil {
			log.Print("发生了致命错误:", err)
		}
	}()

	//	当前目录
	rootDir := filepath.Dir(os.Args[0])

	//	读取配置文件
	err := config.SetRootDir(rootDir)
	if err != nil {
		log.Fatal("读取配置文件错误: ", err)
		return
	}

	start, err := time.Parse("20060102", "20151101")
	if err != nil {
		log.Fatal(err)
		return
	}
	end, err := time.Parse("20060102", "20151111")
	if err != nil {
		log.Fatal(err)
		return
	}

	peroids, err := data.QueryPeroids("america", "aapl", start, end)
	if err != nil {
		log.Fatal("查询服务器数据出错: ", err)
		return
	}

	log.Printf("peroids: %d", len(peroids))
	log.Print(peroids[0])
}
