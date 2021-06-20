package cconf

import (
	"github.com/larspensjo/config"
)

var runmode string
var cfg *config.Config

func init() {
	var err error
	cfg, err = config.ReadDefault("./conf/app.conf") //读取配置文件，并返回其Config
	if err != nil {
		panic("config.ReadDefault " + err.Error())
	}

	runmode, err = cfg.String(config.DEFAULT_SECTION, "runmode")
	if err != nil {
		runmode = config.DEFAULT_SECTION //默认select
	}
}

func String(key string) string {
	value, _ := cfg.String(runmode, key)
	return value
}

func Int(key string) (int, error) {
	return cfg.Int(runmode, key)
}
