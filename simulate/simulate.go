package simulate

import "time"

//	模拟
type Simulate struct {
	Market      string    //	市场
	Code        string    //	上市公司
	StartTime   time.Time //	起始时间
	StartAmount float32   //	起始资金
	EndTime     time.Time //	结束时间
	EndAmount   float32   //	结束资金
	Profit      float32   //	利润

}

func (s *Simulate) Start(f func(*Simulate)) error {

	//	获取分时数据

	return nil
}

//	启动趋势
//func (s *Simulate) StartTrend(startTime time.Time, price float32, reason string, long bool, quantity int) *Trend {

//	trend := Trend{
//		StartTime:   startTime,
//		StartPrice:  price,
//		StartReason: reason,
//		Long:        long,
//		Holdings:    []Holding{}}

//	//	趋势启动时的第一个头寸
//	trend.AddHolding(startTime, price, long, quantity)

//	s.Trends = append(s.Trends, trend)
//	s.CurrentTrend = &trend

//	return &trend
//}

////	模拟
//type SimulateSetting struct {
//	Companies   m.CompanyList //	上市公司
//	StartTime   time.Time     //	起始时间
//	StartAmount float32       //	起始资金
//	EndTime     time.Time     //	结束时间
//	EndAmount   float32       //	结束资金

//}
