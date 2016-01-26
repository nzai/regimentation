package trading

import (
	"time"

	"github.com/nzai/regimentation/data"
)

//	趋势
type Trend struct {
	StartTime     time.Time //	起始时间
	StartPrice    float32   //	起始价格
	StartReason   string    //	起始原因
	Direction     bool      //	做多
	EndTime       time.Time //	结束时间
	EndPrice      float32   //	结束价格
	EndReason     string    //	结束原因
	Profit        float32   //	利润
	ProfitPercent float32   //	利润率
	Holdings      []Holding //	头寸
}

type ITrendTradingSystem interface {
	Enter(data.MinuteHistory) (bool, bool, string, error) // 是否进入
	Increase(data.MinuteHistory) (bool, string, error)    // 是否增持
	Exit(data.MinuteHistory) (bool, error)                // 是否止盈
	Stop(data.MinuteHistory) (bool, error)                // 是否止损
}

//	趋势交易系统
type TrendTradingSystem struct {
	*TradingSystem
	CurrentTrend *Trend               //	当前趋势
	Trends       []Trend              //	趋势
	System       *ITrendTradingSystem //	系统
}
