package data

import (
	"fmt"
	"time"
)

type PeroidExtermaIndex struct {
	Market string    //	市场
	Code   string    //	上市公司
	Time   time.Time //	时间(区间后的下一分钟)
	Peroid int       //	区间
	Min    float32   //	波动性均值
	Max    float32   //	真实波动性
}

type PeroidExtermaIndexes struct {
	indexes map[int]map[time.Time]PeroidExtermaIndex //	字典
}

func (t *PeroidExtermaIndexes) Init(histories []MinuteHistory, peroids ...int) error {
	if t.indexes == nil {
		t.indexes = make(map[int]map[time.Time]PeroidExtermaIndex)
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
func (t *PeroidExtermaIndexes) calculate(histories []MinuteHistory, peroid int) (map[time.Time]PeroidExtermaIndex, error) {

	if peroid < 2 {
		return nil, fmt.Errorf("peroid必须大于等于2")
	}

	start := 0
	dict := make(map[time.Time]PeroidExtermaIndex)
	for index, history := range histories {

		if index >= peroid {
			start = index - peroid
		}

		min := histories[start].Low
		max := histories[start].High

		for hi := start + 1; hi < index; hi++ {
			if histories[hi].High > max {
				max = histories[hi].High
			}

			if histories[hi].Low < min {
				min = histories[hi].Low
			}
		}

		dict[history.Time] = PeroidExtermaIndex{
			Market: history.Market,
			Code:   history.Code,
			Time:   history.Time,
			Peroid: peroid,
			Min:    min,
			Max:    max}
	}

	return dict, nil
}

//	查询
func (t *PeroidExtermaIndexes) Get(peroid int, _time time.Time) (*PeroidExtermaIndex, error) {

	indexes, found := t.indexes[peroid]
	if !found {
		return nil, fmt.Errorf("没有找到区间%d在%s的极值", peroid, _time.Format("2006-01-02"))
	}

	index, found := indexes[_time]
	if !found {
		return nil, fmt.Errorf("没有找到区间%d在%s的极值", peroid, _time.Format("2006-01-02"))
	}

	return &index, nil
}
