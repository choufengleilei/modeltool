package conf

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

type DbConfig struct {
	DataSourceConfig `json:"data_source"`
	Config           `json:"config"`
}

type Config struct {
	ModelPath        string   `json:"model_path"`
	ModelReplace     bool     `json:"model_replace"`
	FieldFormatTypes []string `json:"field_format_types"`
}

type DataSourceConfig struct {
	DriverName string `json:"driver_name"`
	Addr       string `json:"addr"`
	Database   string `json:"database"`
	User       string `json:"user"`
	Password   string `json:"password"`
}

var dbConfig DbConfig
var configOnce sync.Once

func GetConfig() *DbConfig {
	configOnce.Do(func() {
		bytes, err := ioutil.ReadFile("conf.json")
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(bytes, &dbConfig)
		if err != nil {
			panic(err)
		}
	})
	return &dbConfig
}
