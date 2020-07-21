package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type StockData struct {
	Name           string  `json:"name"`
	TodayPrice     float64 `json:"price"`
	YesterdayPrice float64 `json:"yestclose"`
	UpdateTime     string  `json:"update"`
}

func GetStockData(stockCode string) (StockData, error) {
	res, err := http.Get(fmt.Sprintf("http://api.money.126.net/data/feed/%s,money.api", stockCode))
	if err != nil {
		return StockData{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return StockData{}, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return StockData{}, err
	}
	result := string(bytes)
	jsonStr := strings.ReplaceAll(strings.ReplaceAll(result, "_ntes_quote_callback({\""+stockCode+"\":", ""), "});", "")
	var stockData StockData
	err = json.Unmarshal([]byte(jsonStr), &stockData)
	if err != nil {
		return StockData{}, err
	}
	return stockData, nil
}
