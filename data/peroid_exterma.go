package data

import (
	"fmt"
	"math"
	"time"
)

type PeroidExtermaIndex struct {
	Market    string    //	市场
	Code      string    //	上市公司
	StartDate time.Time //	起始日期
	EndDate   time.Time //	结束日期
	Peroid    int       //	区间
	Min       float32   //	波动性均值
	Max       float32   //	真实波动性
}

type PeroidExtermaIndexes struct {
	indexes map[int]map[time.Time]PeroidExtermaIndex //	字典
}

func (t *PeroidExtermaIndexes) Init(histories []DayHistory, peroids ...int) error {
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
func (t *PeroidExtermaIndexes) calculate(histories []DayHistory, peroid int) (map[time.Time]PeroidExtermaIndex, error) {

	if peroid < 2 {
		return nil, fmt.Errorf("peroid必须大于等于2")
	}

	queue := make([]DayHistory, 0)
	dict := make(map[time.Time]PeroidExtermaIndex)
	for index, history := range histories {
		if index >= peroid {
			queue = queue[1:]
		}

		queue = append(queue, history)

		min := float32(math.MaxFloat32)
		max := float32(-math.MaxFloat32)

		for _, _history := range queue {
			if _history.High > max {
				max = _history.High
			}

			if _history.Low < min {
				min = _history.Low
			}
		}

		dict[history.Date] = PeroidExtermaIndex{
			Market:    history.Market,
			Code:      history.Code,
			StartDate: queue[0].Date,
			EndDate:   history.Date,
			Peroid:    peroid,
			Min:       min,
			Max:       max}
	}

	return dict, nil
}

func (t *PeroidExtermaIndexes) Get(peroid int, date time.Time) (*PeroidExtermaIndex, error) {

	indexes, found := t.indexes[peroid]
	if !found {
		return nil, fmt.Errorf("没有找到区间%d在%s的海龟", peroid, date.Format("2006-01-02"))
	}

	index, found := indexes[date]
	if !found {
		return nil, fmt.Errorf("没有找到区间%d在%s的极值", peroid, date.Format("2006-01-02"))
	}

	return &index, nil
}
