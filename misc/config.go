package misc

import (
	"encoding/json"
	"os"
	"time"
)

type BackendConfig struct {
	Mongo struct {
		Hosts   string        `json:"hosts"`
		Timeout time.Duration `json:"timeout"`
		Db      string        `json:"db"`
		User    string        `json:"user"`
		Passwd  string        `json:"passwd"`
	} `json:"mongo"`
}

type Config struct {
	Addr string `json:"addr"`
	BackendConfig
}

var (
	Conf *Config
)

func LoadConf(filename string) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(r)
	Conf = &Config{}
	err = decoder.Decode(Conf)
	if err != nil {
		return err
	}
	return nil
}
