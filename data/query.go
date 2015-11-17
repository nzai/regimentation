package data

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nzai/regimentation/config"
	"github.com/nzai/stockrecorder/io"
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

	content, err := io.DownloadString(url)
	if err != nil {
		return nil, err
	}

	//log.Print("content:", content)

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

	var timeString string
	var open, _close, high, low float32
	var volume int64

	peroids := make([]m.Peroid60, 0)
	for _, obj := range objs {

		str, ok := obj.(string)
		if !ok {
			return nil, fmt.Errorf("转换Data出错:%v", obj)
		}

		_, err := fmt.Sscanf(str, "%s %f %f %f %f %d", &timeString, &open, &_close, &high, &low, &volume)
		if err != nil {
			return nil, fmt.Errorf("转换Peroid60出错:%s", err.Error())
		}

		peroids = append(peroids, m.Peroid60{
			Market: market,
			Code:   code,
			Time:   timeString,
			Open:   open,
			Close:  _close,
			High:   high,
			Low:    low,
			Volume: volume})
	}

	return peroids, nil
}
