package data

import (
	"fmt"
	"time"
)

type TurtleIndex struct {
	Market string    //	市场
	Code   string    //	上市公司
	Peroid int       //	区间
	Date   time.Time //	日期
	N      float32   //	波动性均值
	TR     float32   //	真实波动性
}

type TurtleIndexes struct {
	indexes map[int]map[time.Time]TurtleIndex //	字典
}

//	初始化
func (t *TurtleIndexes) Init(histories []DayHistory, peroids ...int) error {

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
func (t *TurtleIndexes) calculate(histories []DayHistory, peroid int) (map[time.Time]TurtleIndex, error) {

	if peroid < 2 {
		return nil, fmt.Errorf("peroid必须大于等于2")
	}

	var yesterdayN float32 = 0
	var n float32 = 0
	var yesterdayClose float32 = 0

	dict := make(map[time.Time]TurtleIndex)
	for index, history := range histories {
		if index > 0 {
			yesterdayClose = histories[index-1].Close
		}

		//	真实波动幅度TR = Max(high – low，history - yesterdayClose，yesterdayClose - low）
		tr := history.High - history.Low
		if history.High-yesterdayClose > tr {
			tr = history.High - yesterdayClose
		}

		if yesterdayClose-history.Low > tr {
			tr = yesterdayClose - history.Low
		}

		//	真实波动幅度的20日指数移动平均值 N = (19 * PDN + TR) / 20
		if index > 0 {
			n = tr / float32(peroid)
		} else {
			n = (float32(peroid-1)*yesterdayN + tr) / float32(peroid)
		}

		dict[history.Date] = TurtleIndex{
			Market: history.Market,
			Code:   history.Code,
			Peroid: peroid,
			Date:   history.Date,
			N:      n,
			TR:     tr}

		yesterdayN = n
	}

	return dict, nil
}

//	查询
func (t *TurtleIndexes) Get(peroid int, date time.Time) (*TurtleIndex, error) {

	indexes, found := t.indexes[peroid]
	if !found {
		return nil, fmt.Errorf("没有找到区间%d在%s的海龟指标", peroid, date.Format("2006-01-02"))
	}

	index, found := indexes[date]
	if !found {
		return nil, fmt.Errorf("没有找到区间%d在%s的海龟指标", peroid, date.Format("2006-01-02"))
	}

	return &index, nil
}
