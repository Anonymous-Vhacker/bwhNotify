package config

import (
	"bwhNotify/util"
	"fmt"

	"github.com/jinzhu/configor"
)

type BWHostConf struct {
	Veid   string `yaml:"veid"`   // 搬瓦工主机识别id
	ApiKey string `yaml:"apiKey"` // 搬瓦工主机api key
}

type DingTalkConf struct {
	AccessToken string `yaml:"accessToken"` // 请求token
	Secret      string `yaml:"secret"`      // 加签密钥
}

type Config struct {
	BWHosts  []BWHostConf `yaml:"BWHosts"`  // 搬瓦工主机信息
	DingTalk DingTalkConf `yaml:"DingTalk"` // 钉钉机器人信息
}

var Conf = &Config{}

func InitConf(configPath string) error {
	exist, _ := util.FileExists(configPath)
	if !exist {
		return fmt.Errorf("config file %s not exist or not regular file", configPath)
	}
	return configor.Load(Conf, configPath)
}
