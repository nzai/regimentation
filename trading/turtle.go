package trading

import (
	"log"
	"math"
	"time"

	"github.com/nzai/regimentation/data"
)

//	海龟设定
type TurtleSetting struct {
	Holding int //	头寸划分
	N       int //	实波动幅度区间
	Enter   int //	入市区间
	Exit    int //	退市区间
	Stop    int //	止损区间
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

	log.Print("初始化结束")

	return nil
}

//	演算
func (t *TurtleSystem) Simulate() {

	log.Print("[Simulate]\t演算开始")

	t.EndAmount = t.StartAmount
	t.BestProfit = -math.MaxFloat32

	go func() {
		//	每天运行一次
		ticker := time.NewTicker(time.Second * 120)

		for _ = range ticker.C {

			log.Printf("[Current]\tHolding:%d N:%d Enter:%d Exit:%d Stop:%d Profit:%f ProfitPercent:%.3f\t%d %d %.2f",
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
	}()

	for stop := t.MinSetting.Stop; stop <= t.MaxSetting.Stop; stop++ {
		for exit := t.MinSetting.Exit; exit <= t.MaxSetting.Exit; exit++ {
			for enter := t.MinSetting.Enter; enter <= t.MaxSetting.Enter; enter++ {
				for n := t.MinSetting.N; n <= t.MaxSetting.N; n++ {
					for holding := t.MinSetting.Holding; holding <= t.MaxSetting.Holding; holding++ {

						t.CurrentSetting = &TurtleSetting{
							Holding: holding,
							N:       n,
							Enter:   enter,
							Exit:    exit,
							Stop:    stop}

						//	新的测试
						newTest := &TurtleSystemTest{
							TurtleSystem:  t,
							TurtleSetting: t.CurrentSetting,
							StartAmount:   t.StartAmount,
							EndAmount:     t.StartAmount,
							Trends:        []Trend{}}

						//	演算
						err := newTest.Simulate()
						if err != nil {
							log.Printf("[Error]\tHolding:%d N:%d Enter:%d Exit:%d Stop:%d Error:%s",
								t.CurrentSetting.Holding,
								t.CurrentSetting.N,
								t.CurrentSetting.Enter,
								t.CurrentSetting.Exit,
								t.CurrentSetting.Stop,
								err.Error())
						}

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
					}
				}
			}
		}
	}

	log.Printf("[Simulate]\t演算结束,最佳设定Holding:%d N:%d Enter:%d Exit:%d Stop:%d Profit:%f ProfitPercent:%.3f",
		t.BestSetting.Holding,
		t.BestSetting.N,
		t.BestSetting.Enter,
		t.BestSetting.Exit,
		t.BestSetting.Stop,
		t.BestProfit,
		t.BestProfitPercent*100)
}
