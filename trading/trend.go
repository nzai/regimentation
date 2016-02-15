package trading

import "time"

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
