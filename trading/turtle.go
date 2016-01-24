package trading

import (
	"fmt"
	"time"

	gtime "github.com/nzai/go-utility/time"
	"github.com/nzai/regimentation/data"
)

type TutleSystem struct {
	*TrendTradingSystem
	DailyPeroids         map[time.Time]data.MinuteHistory //	每日分时指标
	PeroidExtermaIndexes *data.PeroidExtermaIndexes       //	区间极值
	TurtleIndexes        *data.TurtleIndexes              //	海龟指标
	Holding              int                              //	头寸
	TurtleN              int                              //	实波动幅度区间
	TurtleEnter          int                              //	入市区间
	TurtleExit           int                              //	退市区间
	TurtleStop           int                              //	止损区间
}

//	是否进入
func (t *TutleSystem) Enter(peroid data.MinuteHistory) (bool, bool, string, error) {
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
func (t *TutleSystem) Increase(peroid data.MinuteHistory) (bool, string, error) {
	//	没有趋势就没有增持
	if t.CurrentTrend == nil {
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
func (t *TutleSystem) Exit(peroid data.MinuteHistory) (bool, string, error) {
	//	没有趋势就没有止盈
	if t.CurrentTrend == nil {
		return false, "", nil
	}

	return false, "", nil
}

//	是否止损
func (t *TutleSystem) Stop(peroid data.MinuteHistory) (bool, error) {
	return false, nil
}

//	昨天0点
func getYesterdayZero(_time time.Time) time.Time {
	return gtime.BeginOfDay(_time).Add(-time.Hour * 24)
}
