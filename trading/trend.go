package trading

import (
	"time"

	"github.com/nzai/regimentation/data"
)

//	趋势
type Trend struct {
	StartTime   time.Time //	起始时间
	StartPrice  float32   //	起始价格
	StartReason string    //	起始原因
	Direction   bool      //	做多
	EndTime     time.Time //	结束时间
	EndPrice    float32   //	结束价格
	EndReason   string    //	结束原因
	Profit      float32   //	利润
	Holdings    []Holding //	头寸
}

//	增加头寸
func (t *Trend) AddHolding(_time time.Time, price float32, direction bool, quantity int) *Holding {

	holding := Holding{
		StartTime:  _time,
		StartPrice: price,
		Direction:  direction,
		Quantity:   quantity}

	t.Holdings = append(t.Holdings, holding)

	return &holding
}

//	结束趋势
func (t *Trend) End(endTime time.Time, price float32, reason string) {

	t.EndTime = endTime
	t.EndPrice = price
	t.EndReason = reason

	for _, holding := range t.Holdings {
		holding.EndTime = endTime
		holding.EndPrice = price
		if holding.Direction {
			//	做多
			holding.Profit = (holding.EndPrice - holding.StartPrice) * float32(holding.Quantity)
		} else {
			//	做空
			holding.Profit = (holding.StartPrice - holding.EndPrice) * float32(holding.Quantity)
		}

		t.Profit += holding.Profit
	}
}

type ITrendTradingSystem interface {
	Enter(data.MinuteHistory) (bool, bool, string, error) // 是否进入
	Increase(data.MinuteHistory) (bool, string, error)    // 是否增持
	Exit(data.MinuteHistory) (bool, error)                // 是否止盈
	Stop(data.MinuteHistory) (bool, error)                // 是否止损
}

//	趋势交易系统
type TrendTradingSystem struct {
	CurrentTrend *Trend               //	当前趋势
	Trends       []Trend              //	趋势
	System       *ITrendTradingSystem //	系统
}
