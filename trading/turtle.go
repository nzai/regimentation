package trading

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/nzai/regimentation/data"
)

const (
	SimulateGCCount     = 32
	ProgressDelaySecond = 5
)

//	海龟设定
type TurtleSetting struct {
	Holding int //	头寸划分
	N       int //	实波动幅度区间
	Enter   int //	入市区间
	Exit    int //	退市区间
	Stop    int //	止损区间
}

//	海龟演算结果
type TurtleSimulateResult struct {
	Setting       TurtleSetting //	海龟设定
	Profit        float32       //	利润
	ProfitPercent float32       //	利润率
}

//	海龟系统
type TurtleSystem struct {
	Market      string    //	市场
	Code        string    //	上市公司
	StartTime   time.Time //	起始时间
	StartAmount float32   //	起始资金
	EndTime     time.Time //	结束时间
	EndAmount   float32   //	结束资金

	MinSetting TurtleSetting //	最小设定
	MaxSetting TurtleSetting //	最大设定

	CurrentSetting       *TurtleSetting //	当前设定
	CurrentProfit        float32        //	当前利润
	CurrentProfitPercent float32        //	当前利润率

	Caculated   int //	当前的计算量
	TotalAmount int //	总的计算量

	BestSetting       *TurtleSetting //	最佳设定
	BestProfit        float32        //	最佳利润
	BestProfitPercent float32        //	最佳利润率

	MinutePeroids []data.MinuteHistory //	分时数据

	PeroidExtermaIndexes *data.PeroidExtermaIndexes //	区间极值
	TurtleIndexes        *data.TurtleIndexes        //	海龟指标

	simulateResultChannel chan TurtleSimulateResult //	验算结果发送通道
}

//	初始化
func (t *TurtleSystem) Init() error {

	log.Print("初始化开始")

	//	分时数据
	mhs, err := data.QueryMinuteHistories(t.Market, t.Code, t.StartTime, t.EndTime)
	if err != nil {
		return err
	}

	t.MinutePeroids = mhs

	//	区间极值
	minPeroid := t.MinSetting.Enter
	if t.MinSetting.Exit < minPeroid {
		minPeroid = t.MinSetting.Exit
	}

	maxPeroid := t.MaxSetting.Enter
	if t.MaxSetting.Exit > maxPeroid {
		maxPeroid = t.MaxSetting.Exit
	}

	peroids := make([]int, 0)
	for peroid := minPeroid; peroid <= maxPeroid; peroid++ {
		peroids = append(peroids, peroid)
	}
	//	log.Printf("peroids:%v", peroids)
	t.PeroidExtermaIndexes = &data.PeroidExtermaIndexes{}
	t.PeroidExtermaIndexes.Init(t.MinutePeroids, peroids...)

	//	海龟指标
	ns := make([]int, 0)
	for n := t.MinSetting.N; n <= t.MaxSetting.N; n++ {
		ns = append(ns, n)
	}
	//	log.Printf("ns:%v", ns)
	t.TurtleIndexes = &data.TurtleIndexes{}
	t.TurtleIndexes.Init(t.MinutePeroids, ns...)
	t.TotalAmount = (t.MaxSetting.Holding - t.MinSetting.Holding + 1) *
		(t.MaxSetting.N - t.MinSetting.N + 1) *
		(t.MaxSetting.Enter - t.MinSetting.Enter + 1) *
		(t.MaxSetting.Exit - t.MinSetting.Exit + 1) *
		(t.MaxSetting.Stop - t.MinSetting.Stop + 1)
	t.Caculated = 0

	//	初始化演算结果通道
	t.simulateResultChannel = make(chan TurtleSimulateResult)

	log.Print("初始化结束")

	return nil
}

//	演算
func (t *TurtleSystem) Simulate() {

	log.Print("[Simulate]\t演算开始")

	t.EndAmount = t.StartAmount
	t.BestProfit = -math.MaxFloat32

	//	进度显示
	go t.simulateProgress()

	//	演算结果处理
	go t.SimulateResultProcess()

	defer close(t.simulateResultChannel)

	chanSend := make(chan int, SimulateGCCount)
	defer close(chanSend)

	var wg sync.WaitGroup
	wg.Add(t.TotalAmount)

	for stop := t.MinSetting.Stop; stop <= t.MaxSetting.Stop; stop++ {
		for exit := t.MinSetting.Exit; exit <= t.MaxSetting.Exit; exit++ {
			for enter := t.MinSetting.Enter; enter <= t.MaxSetting.Enter; enter++ {
				for n := t.MinSetting.N; n <= t.MaxSetting.N; n++ {
					for holding := t.MinSetting.Holding; holding <= t.MaxSetting.Holding; holding++ {

						//	并发演算
						go func(setting TurtleSetting) {
							_err := t.SimulateOnce(setting)
							if _err != nil {
								log.Print(_err.Error())
							}

							<-chanSend
							wg.Done()
						}(TurtleSetting{Holding: holding, N: n, Enter: enter, Exit: exit, Stop: stop})

						chanSend <- 1
					}
				}
			}
		}
	}

	//	阻塞，直到演算完所有组合
	wg.Wait()

	log.Printf("[Simulate]\t演算结束,最佳设定Holding:%d N:%d Enter:%d Exit:%d Stop:%d Profit:%f ProfitPercent:%.3f",
		t.BestSetting.Holding,
		t.BestSetting.N,
		t.BestSetting.Enter,
		t.BestSetting.Exit,
		t.BestSetting.Stop,
		t.BestProfit,
		t.BestProfitPercent*100)
}

//	演算一种配置
func (t *TurtleSystem) SimulateOnce(setting TurtleSetting) error {

	//	新的测试
	newTest := &TurtleSystemTest{
		TurtleSystem:  t,
		TurtleSetting: &setting,
		StartAmount:   t.StartAmount,
		EndAmount:     t.StartAmount,
		Trends:        []Trend{}}

	//	演算
	err := newTest.Simulate()
	if err != nil {
		return fmt.Errorf("[Error]\tHolding:%d N:%d Enter:%d Exit:%d Stop:%d Error:%s",
			t.CurrentSetting.Holding,
			t.CurrentSetting.N,
			t.CurrentSetting.Enter,
			t.CurrentSetting.Exit,
			t.CurrentSetting.Stop,
			err.Error())
	}

	//	发送演算结果
	t.simulateResultChannel <- TurtleSimulateResult{
		Setting:       setting,
		Profit:        newTest.Profit,
		ProfitPercent: newTest.ProfitPercent}

	return nil
}

//	演算进度
func (t *TurtleSystem) simulateProgress() {
	//	定时任务
	ticker := time.NewTicker(time.Second * ProgressDelaySecond)

	for _ = range ticker.C {
		if t.CurrentSetting != nil {
			log.Printf("[Current]\tHolding:%d N:%d Enter:%d Exit:%d Stop:%d Profit:%f ProfitPercent:%.3f\t%d %d %03.2f%%",
				t.CurrentSetting.Holding,
				t.CurrentSetting.N,
				t.CurrentSetting.Enter,
				t.CurrentSetting.Exit,
				t.CurrentSetting.Stop,
				t.CurrentProfit,
				t.CurrentProfitPercent*100,
				t.Caculated,
				t.TotalAmount,
				float32(t.Caculated)/float32(t.TotalAmount)*100)
		}

		if t.BestSetting != nil {
			log.Printf("[Best]\tHolding:%d N:%d Enter:%d Exit:%d Stop:%d Profit:%f ProfitPercent:%.3f",
				t.BestSetting.Holding,
				t.BestSetting.N,
				t.BestSetting.Enter,
				t.BestSetting.Exit,
				t.BestSetting.Stop,
				t.BestProfit,
				t.BestProfitPercent*100)
		}

	}
}

//	演算结果处理
func (t *TurtleSystem) SimulateResultProcess() {

	for {
		//	从通道中读取演算结果
		result := <-t.simulateResultChannel

		t.CurrentSetting = &result.Setting
		t.CurrentProfit = result.Profit
		t.CurrentProfitPercent = result.ProfitPercent

		if t.CurrentProfit > t.BestProfit {
			t.BestSetting = t.CurrentSetting
			t.BestProfit = t.CurrentProfit
			t.BestProfitPercent = t.CurrentProfitPercent

			log.Printf("[Best]\tHolding:%d N:%d Enter:%d Exit:%d Stop:%d Profit:%f ProfitPercent:%.3f",
				t.BestSetting.Holding,
				t.BestSetting.N,
				t.BestSetting.Enter,
				t.BestSetting.Exit,
				t.BestSetting.Stop,
				t.BestProfit,
				t.BestProfitPercent*100)
		}

		t.Caculated++
	}
}
