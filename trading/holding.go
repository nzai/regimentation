package trading

import "time"

//	头寸
type Holding struct {
	StartTime  time.Time //	起始时间
	StartPrice float32   //	起始价格
	Direction  bool      //	做多
	Quantity   int       //	数量
	//	EndTime    time.Time //	结束时间
	//	EndPrice   float32   //	结束价格
	//	Profit     float32   //	利润
}
