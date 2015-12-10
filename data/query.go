package data

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/nzai/go-utility/net"
	"github.com/nzai/regimentation/config"
	m "github.com/nzai/stockrecorder/market"
	"github.com/nzai/stockrecorder/server/result"
)

//	查询分时数据
func QueryPeroids(market, code, start, end string) ([]m.Peroid60, error) {

	log.Print("ServerAddress:", config.Get().ServerAddress)
	//	url := path.Join(config.Get().ServerAddress, market, code, start, end, "1m")
	//	url := "http://52.69.228.175:602/america/aapl/20151101/20151111/1m"
	url := "http://localhost:602/america/aapl/20151101/20151111/1m"
	log.Print("url:", url)
	content, err := net.DownloadString(url)
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

	peroids := make([]m.Peroid60, 0)
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

		peroids = append(peroids, m.Peroid60{
			Market: upperMarket,
			Code:   upperCode,
			Time:   _time,
			Open:   float32(values[1].(float64)) / 1000,
			Close:  float32(values[2].(float64)) / 1000,
			High:   float32(values[3].(float64)) / 1000,
			Low:    float32(values[4].(float64)) / 1000,
			Volume: int64(values[5].(float64))})
	}

	return peroids, nil
}
