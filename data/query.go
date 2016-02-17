package data

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/nzai/go-utility/net"
	utime "github.com/nzai/go-utility/time"
	"github.com/nzai/regimentation/config"
	"github.com/nzai/stockrecorder/server/result"
)

//	每分钟历史
type MinuteHistory struct {
	Market string
	Code   string
	Time   time.Time
	Open   float32
	Close  float32
	High   float32
	Low    float32
	Volume int64
}

//	查询分时数据
func QueryMinuteHistories(market, code string, start, end time.Time) ([]MinuteHistory, error) {

	log.Print("ServerAddress:", config.Get().ServerAddress)
	//		url := path.Join(config.Get().ServerAddress, market, code, start, end, "1m")
	url := fmt.Sprintf("%s/%s/%s/%s/%s/1m",
		config.Get().ServerAddress,
		strings.ToLower(market),
		strings.ToLower(code),
		start.Format("20060102"),
		end.Format("20060102"))
	//	url := "http://52.69.228.175:602/america/aapl/20151101/20151111/1m"
	//	url := "http://localhost:602/america/aapl/20151101/20151111/1m"
	//	log.Print("url:", url)
	content, err := net.DownloadStringRetry(url, 3, 10)
	if err != nil {
		return nil, err
	}

	r := result.HttpResult{}
	err = json.Unmarshal([]byte(content), &r)
	if err != nil {
		return nil, err
	}

	if !r.Success {
		return nil, fmt.Errorf("从服务器查询分时数据出错:%s", r.Message)
	}

	objs, ok := r.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("转换Data出错:%v", r.Data)
	}

	upperMarket := strings.Title(market)
	upperCode := strings.ToUpper(code)

	var currentTradingDate time.Time
	histories := make([]MinuteHistory, 0)
	for _, obj := range objs {
		//log.Print(reflect.TypeOf(obj).String())
		values, ok := obj.([]interface{})
		if !ok {
			return nil, fmt.Errorf("转换Data item出错:%v", obj)
		}

		_time, err := time.Parse("0601021504", strconv.FormatInt(int64(values[0].(float64)), 10))
		if err != nil {
			return nil, err
		}

		date := utime.BeginOfDay(_time)
		if !date.Equal(currentTradingDate) {
			currentTradingDate = date
		}

		histories = append(histories, MinuteHistory{
			Market: upperMarket,
			Code:   upperCode,
			Time:   _time,
			Open:   float32(values[1].(float64)) / 1000,
			Close:  float32(values[2].(float64)) / 1000,
			High:   float32(values[3].(float64)) / 1000,
			Low:    float32(values[4].(float64)) / 1000,
			Volume: int64(values[5].(float64))})
	}

	return histories, nil
}

//	区间历史
type PeroidHistory struct {
	Time    time.Time
	Open    float32
	Close   float32
	High    float32
	Low     float32
	Volume  int64
	Minutes []MinuteHistory
}

//	转化区间历史
func ParsePeroidHistory(histories []MinuteHistory, peroid int) ([]PeroidHistory, error) {
	if len(histories) == 0 {
		return []PeroidHistory{}, nil
	}

	var ph *PeroidHistory = nil
	end := histories[0].Time.Add(-time.Minute)

	phs := make([]PeroidHistory, 0)
	for _, history := range histories {

		if history.Time.After(end) {
			if ph != nil {
				phs = append(phs, *ph)
			}

			end = history.Time.Add(time.Minute * time.Duration(peroid))

			ph = &PeroidHistory{
				Time:    history.Time,
				Open:    history.Open,
				Close:   history.Close,
				High:    history.High,
				Low:     history.Low,
				Volume:  history.Volume,
				Minutes: make([]MinuteHistory, 0)}

			continue
		}

		if history.High > ph.High {
			ph.High = history.High
		}

		if history.Low < ph.Low {
			ph.Low = history.Low
		}

		ph.Close = history.Close
		ph.Volume += history.Volume
		ph.Minutes = append(ph.Minutes, history)
	}

	if ph != nil {
		phs = append(phs, *ph)
	}

	return phs, nil
}
