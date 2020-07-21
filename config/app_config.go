package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

type FundConfig struct {
	Code string `json:"code"`
}

type AppConfig struct {
	Font  string       `json:"font"`
	Funds []FundConfig `json:"funds"`
}

var configOnce sync.Once
var appConfig *AppConfig

func GetAppConfig() *AppConfig {
	configOnce.Do(func() {
		bytes, err := ioutil.ReadFile("./config.json")
		if err != nil {
			panic(err)
		}
		configJson := string(bytes)
		err = json.Unmarshal([]byte(configJson), &appConfig)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	})
	return appConfig
}

func UpdateAppConfig() {
	jsonStr, err := json.Marshal(appConfig)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("./config.json", jsonStr, 0644)
	if err != nil {
		panic(err)
	}
}
