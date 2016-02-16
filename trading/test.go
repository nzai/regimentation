package trading

import (
	"fmt"

	"github.com/nzai/regimentation/data"
)

type TurtleSystemTest struct {
	TurtleSystem  *TurtleSystem  //	海龟系统
	TurtleSetting *TurtleSetting //	海龟设定

	StartAmount   float32 //	起始资金
	EndAmount     float32 //	结束资金
	Profit        float32 //	利润
	ProfitPercent float32 //	利润率

	CurrentTrend *Trend  //	当前趋势
	Trends       []Trend //	趋势
}

//	是否入市
func (t *TurtleSystemTest) Enter(peroid data.MinuteHistory) (bool, bool, string, error) {
	//	当前在趋势中就不重复进入了
	if t.CurrentTrend != nil {
		return false, false, "", nil
	}

	//	查上一交易日的区间极值
	peroid_exterma_last_trading_day, err := t.TurtleSystem.PeroidExtermaIndexes.Get(t.TurtleSetting.Enter, peroid.Time)
	if err != nil {
		return false, false, "", fmt.Errorf("[TutleSystem]\t%s", err.Error())
	}

	//	是否超过极值
	if peroid.High > peroid_exterma_last_trading_day.Max {
		return true, true, fmt.Sprintf("分时价格%.2f突破%d日最大值%.2f", peroid.High, peroid_exterma_last_trading_day.Peroid, peroid_exterma_last_trading_day.Max), nil
	}

	if peroid.Low < peroid_exterma_last_trading_day.Min {
		return true, false, fmt.Sprintf("分时价格%.2f突破%d日最小值%.2f", peroid.Low, peroid_exterma_last_trading_day.Peroid, peroid_exterma_last_trading_day.Min), nil
	}

	return false, false, "", nil
}

//	是否增持
func (t *TurtleSystemTest) Increase(peroid data.MinuteHistory) (bool, string, error) {
	//	没有趋势就没有增持
	if t.CurrentTrend == nil {
		return false, "", nil
	}

	//	满仓了就不增持了
	if len(t.CurrentTrend.Holdings) >= t.TurtleSetting.Holding {
		return false, "", nil
	}

	//	海龟指标
	turtle_last_trading_date, err := t.TurtleSystem.TurtleIndexes.Get(t.TurtleSetting.N, peroid.Time)
	if err != nil {
		return false, "", fmt.Errorf("[TutleSystem]\t%s", err.Error())
	}

	lastHolding := t.CurrentTrend.Holdings[len(t.CurrentTrend.Holdings)-1]

	//	是否超过海龟加仓指标
	if t.CurrentTrend.Direction && peroid.High > lastHolding.StartPrice+turtle_last_trading_date.N/2 {
		return true, fmt.Sprintf("分时价格%.2f突破做多加仓最小值%.2f", peroid.High, lastHolding.StartPrice+turtle_last_trading_date.N/2), nil
	}

	if !t.CurrentTrend.Direction && peroid.Low < lastHolding.StartPrice-turtle_last_trading_date.N/2 {
		return true, fmt.Sprintf("分时价格%.2f突破做空加仓最小值%.2f", peroid.Low, lastHolding.StartPrice-turtle_last_trading_date.N/2), nil
	}

	return false, "", nil
}

//	是否止盈
func (t *TurtleSystemTest) Exit(peroid data.MinuteHistory) (bool, string, error) {
	//	没有趋势就没有止盈
	if t.CurrentTrend == nil {
		return false, "", nil
	}

	//	查昨天的区间极值
	peroid_exterma_last_trading_day, err := t.TurtleSystem.PeroidExtermaIndexes.Get(t.TurtleSetting.Exit, peroid.Time)
	if err != nil {
		return false, "", fmt.Errorf("[TutleSystem]\t%s", err.Error())
	}

	if t.CurrentTrend.Direction && peroid.Low < peroid_exterma_last_trading_day.Min {
		return true, fmt.Sprintf("趋势结束 分时价格%.2f突破%d日最小值%.2f", peroid.Low, t.TurtleSetting.Exit, peroid_exterma_last_trading_day.Min), nil
	}

	if !t.CurrentTrend.Direction && peroid.High > peroid_exterma_last_trading_day.Max {
		return true, fmt.Sprintf("趋势结束 分时价格%.2f突破%d日最大值%.2f", peroid.High, t.TurtleSetting.Exit, peroid_exterma_last_trading_day.Max), nil
	}

	return false, "", nil
}

//	是否止损
func (t *TurtleSystemTest) Stop(peroid data.MinuteHistory) (bool, string, error) {
	//	没有趋势就没有止损
	if t.CurrentTrend == nil {
		return false, "", nil
	}

	lastHolding := t.CurrentTrend.Holdings[len(t.CurrentTrend.Holdings)-1]

	//	海龟指标
	turtle_last_trading_day, err := t.TurtleSystem.TurtleIndexes.Get(t.TurtleSetting.N, peroid.Time)
	if err != nil {
		return false, "", fmt.Errorf("[TutleSystem]\t%s", err.Error())
	}

	if t.CurrentTrend.Direction && peroid.Low < lastHolding.StartPrice-turtle_last_trading_day.N*float32(t.TurtleSetting.Stop) {
		return true, fmt.Sprintf("止损 分时价格%.2f突破%d日止损价%.2f", peroid.Low, t.TurtleSetting.Stop, lastHolding.StartPrice-turtle_last_trading_day.N*float32(t.TurtleSetting.Stop)), nil
	}

	if !t.CurrentTrend.Direction && peroid.High > lastHolding.StartPrice+turtle_last_trading_day.N*float32(t.TurtleSetting.Stop) {
		return true, fmt.Sprintf("止损 分时价格%.2f突破%d日止损价%.2f", peroid.High, t.TurtleSetting.Stop, lastHolding.StartPrice+turtle_last_trading_day.N*float32(t.TurtleSetting.Stop)), nil
	}

	return false, "", nil
}

//	入市
func (t *TurtleSystemTest) DoEnter(peroid data.MinuteHistory, direction bool, reason string) {

	startPrice := peroid.High
	if !direction {
		startPrice = peroid.Low
	}

	quantity := int(t.EndAmount / (startPrice * float32(t.TurtleSetting.Holding)))

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

	t.CurrentTrend = &trend
}

//	增持
func (t *TurtleSystemTest) DoIncrease(peroid data.MinuteHistory, reason string) {

	startPrice := peroid.High
	if !t.CurrentTrend.Direction {
		startPrice = peroid.Low
	}

	quantity := int(t.EndAmount / (startPrice * float32(t.TurtleSetting.Holding)))

	t.CurrentTrend.Holdings = append(t.CurrentTrend.Holdings, Holding{
		StartTime:  peroid.Time,
		StartPrice: startPrice,
		Direction:  t.CurrentTrend.Direction,
		Quantity:   quantity})
}

//	止盈
func (t *TurtleSystemTest) DoExit(peroid data.MinuteHistory, reason string) {

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

	//	log.Printf("profit:%f  q:%d  s:%f e:%f", profit, len(t.CurrentTrend.Holdings), t.CurrentTrend.StartPrice, t.CurrentTrend.EndPrice)
	t.EndAmount += profit
	t.Profit = t.EndAmount - t.StartAmount
	t.ProfitPercent = t.Profit / t.StartAmount

	t.Trends = append(t.Trends, *t.CurrentTrend)

	//	趋势结束
	t.CurrentTrend = nil
}

//	演算
func (t *TurtleSystemTest) Simulate() error {

	for _, peroid := range t.TurtleSystem.MinutePeroids {
		//	是否入市
		enter, direction, reason, err := t.Enter(peroid)
		if err != nil {
			return err
		}

		//	入市
		if enter {
			t.DoEnter(peroid, direction, reason)
			continue
		}

		//	是否增持
		increase, reason, err := t.Increase(peroid)
		if err != nil {
			return err
		}

		//	增持
		if increase {
			t.DoIncrease(peroid, reason)
			continue
		}

		//	是否止盈
		exit, reason, err := t.Exit(peroid)
		if err != nil {
			return err
		}

		//	止盈
		if exit {
			t.DoExit(peroid, reason)
			continue
		}

		//	是否止损
		stop, reason, err := t.Stop(peroid)
		if err != nil {
			return err
		}

		//	止损
		if stop {
			t.DoExit(peroid, reason)
			continue
		}
	}

	if t.CurrentTrend != nil {
		t.DoExit(t.TurtleSystem.MinutePeroids[len(t.TurtleSystem.MinutePeroids)-1], "演算结束")
	}

	return nil
}
