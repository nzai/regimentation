package main

import (
	"log"
	"time"

	"github.com/nzai/go-utility/path"
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

	root, err := path.GetStartupDir()
	if err != nil {
		log.Fatal("获取起始目录失败: ", err)
		return
	}

	//	读取配置文件
	err = config.SetRootDir(root)
	if err != nil {
		log.Fatal("读取配置文件错误: ", err)
	}

	start, err := time.Parse("20060102", "20151101")
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
		MinSetting:  trading.TurtleSetting{Holding: 1, N: 2, Enter: 2, Exit: 2, Stop: 2},
		MaxSetting:  trading.TurtleSetting{Holding: 5, N: 20, Enter: 20, Exit: 20, Stop: 20}}

	//	初始化
	err = system.Init()
	if err != nil {
		log.Fatal(err)
	}

	//	演算
	system.Simulate()

	//	新的测试
	//	newTest := &trading.TurtleSystemTest{
	//		TurtleSystem: &system,
	//		TurtleSetting: &trading.TurtleSetting{
	//			Holding: 2,
	//			N:       2,
	//			Enter:   2,
	//			Exit:    2,
	//			Stop:    2},
	//		StartAmount: system.StartAmount,
	//		EndAmount:   system.StartAmount,
	//		Trends:      []trading.Trend{}}

	//	//	演算
	//	err = newTest.Simulate()
	//	if err != nil {
	//		log.Printf("[Error]\tHolding:%d N:%d Enter:%d Exit:%d Stop:%d Error:%s",
	//			newTest.TurtleSetting.Holding,
	//			newTest.TurtleSetting.N,
	//			newTest.TurtleSetting.Enter,
	//			newTest.TurtleSetting.Exit,
	//			newTest.TurtleSetting.Stop,
	//			err.Error())
	//	}

	//	for _, trend := range newTest.Trends {

	//		direction := "做多"
	//		if !trend.Direction {
	//			direction = "做空"
	//		}

	//		var quantity int = 0
	//		for _, holding := range trend.Holdings {
	//			quantity += holding.Quantity
	//		}

	//		log.Printf("%s %s %s %s %s 数量:%d 利润:%.2f",
	//			trend.StartTime.Format("2006-01-02 15:04"),
	//			trend.StartReason,
	//			direction,
	//			trend.EndTime.Format("2006-01-02 15:04"),
	//			trend.EndReason,
	//			quantity,
	//			trend.Profit)
	//	}

	//	log.Print(len(newTest.Trends))

	//	log.Printf("[Current]\tHolding:%d N:%d Enter:%d Exit:%d Stop:%d Profit:%f ProfitPercent:%.3f",
	//		newTest.TurtleSetting.Holding,
	//		newTest.TurtleSetting.N,
	//		newTest.TurtleSetting.Enter,
	//		newTest.TurtleSetting.Exit,
	//		newTest.TurtleSetting.Stop,
	//		newTest.Profit,
	//		newTest.ProfitPercent*100)
}
