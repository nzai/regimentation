package trading

import "github.com/nzai/regimentation/data"

//	交易系统接口
type TradingSystem interface {
	Enter(data.Peroid) (bool, error)
	Increase(data.Peroid) (bool, error)
	Exit(data.Peroid) (bool, error)
	Stop(data.Peroid) (bool, error)
}

var (
	systems []TradingSystem = make([]TradingSystem, 0)
)

func Register(ts TradingSystem) {
	systems = append(systems, ts)
}

//	头寸
type Holding struct {
	StartTime  string  //	起始时间
	StartPrice float32 //	起始价格
	Long       bool    //	做多
	Quantity   int     //	数量
	EndTime    string  //	结束时间
	EndPrice   float32 //	结束价格
	Profit     float32 //	利润
}

//	趋势
type Trend struct {
	StartTime   string    //	起始时间
	StartPrice  float32   //	起始价格
	StartReason string    //	起始原因
	Long        bool      //	做多
	EndTime     string    //	结束时间
	EndPrice    float32   //	结束价格
	EndReason   string    //	结束原因
	Profit      float32   //	利润
	Holdings    []Holding //	头寸
}

//	增加头寸
func (t *Trend) AddHolding(time string, price float32, long bool, quantity int) *Holding {

	holding := Holding{
		StartTime:  time,
		StartPrice: price,
		Long:       long,
		Quantity:   quantity}

	t.Holdings = append(t.Holdings, holding)

	return &holding
}

//	结束趋势
func (t *Trend) End(time string, price float32, reason string) {

	t.EndTime = time
	t.EndPrice = price
	t.EndReason = reason

	for _, holding := range t.Holdings {
		holding.EndTime = time
		holding.EndPrice = price
		if holding.Long {
			//	做多
			holding.Profit = (holding.EndPrice - holding.StartPrice) * float32(holding.Quantity)
		} else {
			//	做空
			holding.Profit = (holding.StartPrice - holding.EndPrice) * float32(holding.Quantity)
		}

		t.Profit += holding.Profit
	}
}

//	模拟
type Simulate struct {
	Market       string  //	市场
	Code         string  //	上市公司
	StartTime    string  //	起始时间
	StartAmount  float32 //	起始资金
	EndTime      string  //	结束时间
	EndAmount    float32 //	结束资金
	Profit       float32 //	利润
	CurrentTrend *Trend  //	当前趋势
	Trends       []Trend //	趋势
}

func (s *Simulate) Start(f func(*Simulate, data.Peroid)) error {
	
	//	获取分时数据
	
	
	return nil
}

//	启动趋势
func (s *Simulate) StartTrend(time string, price float32, reason string, long bool, quantity int) *Trend {

	trend := Trend{
		StartTime:   time,
		StartPrice:  price,
		StartReason: reason,
		Long:        long,
		Holdings:    []Holding{}}

	//	趋势启动时的第一个头寸
	trend.AddHolding(time, price, long, quantity)

	s.Trends = append(s.Trends, trend)
	s.CurrentTrend = &trend

	return &trend
}
