package trading

import (
	"time"

	"github.com/nzai/regimentation/data"
)

//	交易系统接口
type ITradingSystem interface {
	//	Init([]m.Peroid60) error           // 初始化
	Enter(data.MinuteHistory) (bool, error)    // 是否进入
	Increase(data.MinuteHistory) (bool, error) // 是否增持
	Exit(data.MinuteHistory) (bool, error)     // 是否止盈
	Stop(data.MinuteHistory) (bool, error)     // 是否止损
}

var (
	systems []ITradingSystem = make([]ITradingSystem, 0)
)

func Register(ts ITradingSystem) {
	systems = append(systems, ts)
}

//	交易系统
type TradingSystem struct {
	Market        string               //	市场
	Code          string               //	上市公司
	StartTime     time.Time            //	起始时间
	StartAmount   float32              //	起始资金
	EndTime       time.Time            //	结束时间
	EndAmount     float32              //	结束资金
	Peroids       []data.MinuteHistory //	分时数据
	System        ITradingSystem       //	交易法则
	Profit        float32              //	利润
	ProfitPercent float32              //	利润率
}
