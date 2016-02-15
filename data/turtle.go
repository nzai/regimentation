package data

import (
	"fmt"
	"time"
)

type TurtleIndex struct {
	Market string    //	市场
	Code   string    //	上市公司
	Peroid int       //	区间
	Time   time.Time //	时间(区间后的下一分钟)
	N      float32   //	波动性均值
	TR     float32   //	真实波动性
}

type TurtleIndexes struct {
	indexes map[int]map[time.Time]TurtleIndex //	字典
}

//	初始化
func (t *TurtleIndexes) Init(histories []MinuteHistory, peroids ...int) error {
	if t.indexes == nil {
		t.indexes = make(map[int]map[time.Time]TurtleIndex)
	}

	for _, peroid := range peroids {
		dict, err := t.calculate(histories, peroid)
		if err != nil {
			return err
		}

		t.indexes[peroid] = dict
	}

	return nil
}

//	计算
func (t *TurtleIndexes) calculate(histories []MinuteHistory, peroid int) (map[time.Time]TurtleIndex, error) {

	if peroid < 2 {
		return nil, fmt.Errorf("peroid必须大于等于2")
	}

	var yesterdayN float32 = 0
	var tr float32 = 0
	var n float32 = 0

	dict := make(map[time.Time]TurtleIndex)
	for index, history := range histories {

		//	真实波动幅度TR = Max(high – low，history - yesterdayClose，yesterdayClose - low）

		if index <= 1 {
			tr = histories[0].High - histories[0].Low
			n = tr / float32(peroid)
		} else {
			tr = histories[index-1].High - histories[index-1].Low

			if histories[index-1].High-histories[index-2].Close > tr {
				tr = histories[index-1].High - histories[index-2].Close
			}

			if histories[index-2].Close-histories[index-1].Low > tr {
				tr = histories[index-2].Close - histories[index-1].Low
			}

			//	真实波动幅度的20日指数移动平均值 N = (19 * PDN + TR) / 20
			n = (float32(peroid-1)*yesterdayN + tr) / float32(peroid)
		}

		dict[history.Time] = TurtleIndex{
			Market: history.Market,
			Code:   history.Code,
			Peroid: peroid,
			Time:   history.Time,
			N:      n,
			TR:     tr}

		yesterdayN = n
	}

	return dict, nil
}

//	查询
func (t *TurtleIndexes) Get(peroid int, _time time.Time) (*TurtleIndex, error) {

	indexes, found := t.indexes[peroid]
	if !found {
		return nil, fmt.Errorf("没有找到区间%d在%s的海龟指标", peroid, _time.Format("2006-01-02"))
	}

	index, found := indexes[_time]
	if !found {
		return nil, fmt.Errorf("没有找到区间%d在%s的海龟指标", peroid, _time.Format("2006-01-02"))
	}

	return &index, nil
}
