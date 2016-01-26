package trading

import (
	"fmt"
	"time"

	gtime "github.com/nzai/go-utility/time"
	"github.com/nzai/regimentation/data"
)

type TurtleSystem struct {
	Market        string               //	市场
	Code          string               //	上市公司
	StartTime     time.Time            //	起始时间
	StartAmount   float32              //	起始资金
	EndTime       time.Time            //	结束时间
	EndAmount     float32              //	结束资金
	Peroids       []data.MinuteHistory //	分时数据
	Profit        float32              //	利润
	ProfitPercent float32              //	利润率

	DailyPeroids         map[time.Time]data.MinuteHistory //	每日分时指标
	PeroidExtermaIndexes *data.PeroidExtermaIndexes       //	区间极值
	TurtleIndexes        *data.TurtleIndexes              //	海龟指标

	TurlteHolding int //	头寸划分
	TurtleN       int //	实波动幅度区间
	TurtleEnter   int //	入市区间
	TurtleExit    int //	退市区间
	TurtleStop    int //	止损区间

	CurrentTrend *Trend  //	当前趋势
	Trends       []Trend //	趋势
}

//	昨天0点
func getYesterdayZero(_time time.Time) time.Time {
	return gtime.BeginOfDay(_time).Add(-time.Hour * 24)
}

//	是否入市
func (t *TurtleSystem) Enter(peroid data.MinuteHistory) (bool, bool, string, error) {
	//	当前在趋势中就不重复进入了
	if t.CurrentTrend != nil {
		return false, false, "", nil
	}

	//	昨天0点
	yesterdayZero := getYesterdayZero(peroid.Time)

	//	查昨天的区间极值
	peroid_exterma_yesterday, err := t.PeroidExtermaIndexes.Get(t.TurtleEnter, yesterdayZero)
	if err != nil {
		return false, false, "", fmt.Errorf("[TutleSystem]\t%s", err.Error())
	}

	//	是否超过极值
	if peroid.High > peroid_exterma_yesterday.Max {
		return true, true, fmt.Sprintf("分时价格%.2d突破%d日最大值%.2d", peroid.High, peroid_exterma_yesterday.Peroid, peroid_exterma_yesterday.Max), nil
	}

	if peroid.Low < peroid_exterma_yesterday.Min {
		return true, false, fmt.Sprintf("分时价格%.2d突破%d日最小值%.2d", peroid.High, peroid_exterma_yesterday.Peroid, peroid_exterma_yesterday.Min), nil
	}

	return false, false, "", nil
}

//	是否增持
func (t *TurtleSystem) Increase(peroid data.MinuteHistory) (bool, string, error) {
	//	没有趋势就没有增持
	if t.CurrentTrend == nil {
		return false, "", nil
	}

	//	满仓了就不增持了
	if len(t.CurrentTrend.Holdings) >= t.TurlteHolding {
		return false, "", nil
	}

	//	昨天0点
	yesterdayZero := getYesterdayZero(peroid.Time)

	//	海龟指标
	turtleYesterday, err := t.TurtleIndexes.Get(t.TurtleN, yesterdayZero)
	if err != nil {
		return false, "", fmt.Errorf("[TutleSystem]\t%s", err.Error())
	}

	lastHolding := t.CurrentTrend.Holdings[len(t.CurrentTrend.Holdings)-1]

	//	是否超过海龟加仓指标
	if t.CurrentTrend.Direction && peroid.High > lastHolding.StartPrice+turtleYesterday.N/2 {
		return true, fmt.Sprintf("分时价格%.2d突破做多加仓最小值%.2d", peroid.High, lastHolding.StartPrice+turtleYesterday.N/2), nil
	}

	if !t.CurrentTrend.Direction && peroid.Low < lastHolding.StartPrice-turtleYesterday.N/2 {
		return true, fmt.Sprintf("分时价格%.2d突破做空加仓最小值%.2d", peroid.Low, lastHolding.StartPrice-turtleYesterday.N/2), nil
	}

	return false, "", nil
}

//	是否止盈
func (t *TurtleSystem) Exit(peroid data.MinuteHistory) (bool, string, error) {
	//	没有趋势就没有止盈
	if t.CurrentTrend == nil {
		return false, "", nil
	}

	//	昨天0点
	yesterdayZero := getYesterdayZero(peroid.Time)

	//	查昨天的区间极值
	peroid_exterma_yesterday, err := t.PeroidExtermaIndexes.Get(t.TurtleExit, yesterdayZero)
	if err != nil {
		return false, "", fmt.Errorf("[TutleSystem]\t%s", err.Error())
	}

	if (t.CurrentTrend.Direction && peroid.Low > peroid_exterma_yesterday.Min) ||
		(!t.CurrentTrend.Direction && peroid.High < peroid_exterma_yesterday.Max) {
		return true, "趋势结束", nil
	}

	return false, "", nil
}

//	是否止损
func (t *TurtleSystem) Stop(peroid data.MinuteHistory) (bool, string, error) {
	//	没有趋势就没有止损
	if t.CurrentTrend == nil {
		return false, "", nil
	}

	lastHolding := t.CurrentTrend.Holdings[len(t.CurrentTrend.Holdings)-1]
	//	昨天0点
	yesterdayZero := getYesterdayZero(lastHolding.StartTime)

	//	海龟指标
	turtleYesterday, err := t.TurtleIndexes.Get(t.TurtleN, yesterdayZero)
	if err != nil {
		return false, "", fmt.Errorf("[TutleSystem]\t%s", err.Error())
	}

	if (t.CurrentTrend.Direction && peroid.Low < lastHolding.StartPrice-turtleYesterday.N*float32(t.TurtleStop)) ||
		(!t.CurrentTrend.Direction && peroid.High > lastHolding.StartPrice-turtleYesterday.N*float32(t.TurtleStop)) {
		return true, "止损", nil
	}

	return false, "", nil
}

//	入市
func (t *TurtleSystem) DoEnter(peroid data.MinuteHistory, direction bool, reason string) {

	startPrice := peroid.High
	if !direction {
		startPrice = peroid.Low
	}

	quantity := int(t.EndAmount / float32(t.TurlteHolding))

	//	启动新趋势
	trend := Trend{
		StartTime:   peroid.Time,
		StartPrice:  startPrice,
		StartReason: reason,
		Direction:   direction,
		Holdings:    make([]Holding, 0)}

	//	第一个头寸
	trend.Holdings = append(trend.Holdings, Holding{
		StartTime:  trend.StartTime,
		StartPrice: trend.StartPrice,
		Direction:  trend.Direction,
		Quantity:   quantity})

	t.Trends = append(t.Trends, trend)
	t.CurrentTrend = &trend
}

//	增持
func (t *TurtleSystem) DoIncrease(peroid data.MinuteHistory, reason string) {
	startPrice := peroid.High
	if !t.CurrentTrend.Direction {
		startPrice = peroid.Low
	}

	quantity := int(t.EndAmount / float32(t.TurlteHolding))

	t.CurrentTrend.Holdings = append(t.CurrentTrend.Holdings, Holding{
		StartTime:  peroid.Time,
		StartPrice: startPrice,
		Direction:  t.CurrentTrend.Direction,
		Quantity:   quantity})
}

//	止盈
func (t *TurtleSystem) DoExit(peroid data.MinuteHistory, reason string) {
	endPrice := peroid.Low
	if !t.CurrentTrend.Direction {
		endPrice = peroid.High
	}

	var profit float32 = 0
	for _, holding := range t.CurrentTrend.Holdings {
		if holding.Direction {
			profit += (endPrice - holding.StartPrice) * float32(holding.Quantity)
		} else {
			profit += (holding.StartPrice - endPrice) * float32(holding.Quantity)
		}
	}

	t.CurrentTrend.EndPrice = endPrice
	t.CurrentTrend.EndTime = peroid.Time
	t.CurrentTrend.EndReason = reason
	t.CurrentTrend.Profit = profit
	t.CurrentTrend.ProfitPercent = profit / t.EndAmount

	t.EndTime = peroid.Time
	t.EndAmount += profit

	//	趋势结束
	t.CurrentTrend = nil
}

func (t *TurtleSystem) Init() {

}
