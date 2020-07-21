package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type FundData struct {
	FundName string `json:"name"`
	FundCode string `json:"fundcode"`
	FundUd   string `json:"gszzl"`
}

func GetFundData(fundCode string) (FundData, error) {
	res, err := http.Get(fmt.Sprintf("http://fundgz.1234567.com.cn/js/%s.js?rt=%d", fundCode, time.Now().Unix()*1000))
	if err != nil {
		return FundData{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return FundData{}, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}
	//f, err := os.Open("/Users/chenyitao/Desktop/180012.html")
	//if err != nil {
	//	panic(err)
	//}
	//defer f.Close()
	//var r io.Reader
	//r = f

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return FundData{}, err
	}
	jsoup := string(bytes)
	jsonStr := strings.ReplaceAll(strings.ReplaceAll(jsoup, "jsonpgz(", ""), ");", "")
	var fundData FundData
	err = json.Unmarshal([]byte(jsonStr), &fundData)
	if err != nil {
		return FundData{}, err
	}
	return fundData, nil
}
