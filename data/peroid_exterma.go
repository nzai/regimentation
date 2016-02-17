package data

import (
	"fmt"
	"time"
)

type PeroidExtermaIndex struct {
	High float32 //	最大值
	Low  float32 //	最小值
}

type PeroidExtermaIndexes struct {
	indexes map[int]map[time.Time]PeroidExtermaIndex //	字典
}

func (t *PeroidExtermaIndexes) Init(histories []PeroidHistory, peroids ...int) error {
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
func (t *PeroidExtermaIndexes) calculate(histories []PeroidHistory, peroid int) (map[time.Time]PeroidExtermaIndex, error) {

	if peroid < 2 {
		return nil, fmt.Errorf("peroid必须大于等于2")
	}

	start := 0
	dict := make(map[time.Time]PeroidExtermaIndex)
	for index, history := range histories {

		if index >= peroid {
			start = index - peroid
		}

		high := histories[start].High
		low := histories[start].Low

		for hi := start + 1; hi < index; hi++ {
			if histories[hi].High > high {
				high = histories[hi].High
			}

			if histories[hi].Low < low {
				low = histories[hi].Low
			}
		}

		dict[history.Time] = PeroidExtermaIndex{High: high, Low: low}
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
