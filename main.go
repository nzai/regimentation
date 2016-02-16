package main

import (
	"log"
	"time"

	"github.com/nzai/regimentation/config"

	"github.com/nzai/regimentation/trading"
)

func main() {

	defer func() {
		// 捕获panic异常
		if err := recover(); err != nil {
			log.Print("发生了致命错误:", err)
		}
	}()

	//	读取配置文件
	err := config.ReadConfig()
	if err != nil {
		log.Fatal("读取配置文件错误: ", err)
	}

	start, err := time.Parse("20060102", "20150801")
	if err != nil {
		log.Fatal(err)
	}
	end, err := time.Parse("20060102", "20160201")
	if err != nil {
		log.Fatal(err)
	}

	//	海龟交易系统演算范围
	system := trading.TurtleSystem{
		Market:      "America",
		Code:        "AAPL",
		StartTime:   start,
		StartAmount: 100000,
		EndTime:     end,
		MinSetting:  trading.TurtleSetting{Holding: 1, N: 60, Enter: 60, Exit: 60, Stop: 60},
		MaxSetting:  trading.TurtleSetting{Holding: 8, N: 6000, Enter: 6000, Exit: 6000, Stop: 6000},
		Step:        60}

	//	初始化
	err = system.Init()
	if err != nil {
		log.Fatal(err)
	}

	//	演算
	system.Simulate()

	//	go system.SimulateResultProcess()

	//	system.SimulateOnce(system.MinSetting)
	//	log.Printf("Profit:%.3f", system.CurrentProfit)

	//	time.Sleep(time.Second * 5)
	//	极值测试
	//	for index, history := range system.MinutePeroids {
	//		if index >= 40 {
	//			break
	//		}

	//		//		pe, err := system.PeroidExtermaIndexes.Get(4, history.Time)
	//		//		if err != nil {
	//		//			log.Fatal(err)
	//		//		}

	//		//		log.Printf("%s %.3f %.3f  [4]  %.3f %.3f", history.Time.Format("2006-01-02 15:04"), history.High, history.Low, pe.Max, pe.Min)

	//		tu, err := system.TurtleIndexes.Get(4, history.Time)
	//		if err != nil {
	//			log.Fatal(err)
	//		}

	//		log.Printf("%s %.3f %.3f %.3f  [4]  %.3f %.3f",
	//			history.Time.Format("2006-01-02 15:04"),
	//			history.High,
	//			history.Low,
	//			history.Close,
	//			tu.TR,
	//			tu.N)
	//	}
}
