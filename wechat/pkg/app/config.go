package app

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug            bool   `json:"debug"`
	HTTPAddr         string `json:"http_addr" default:":80"`
	Mysql            string `json:"mysql" default:"root:shengdian@tcp(mysql:3306)/wechat?charset=utf8mb4"`
	Redis            string `json:"redis" default:"redis:6379"`
	InternalToken    string `json:"internal_token"`
	Gateway          string `json:"gateway" default:"http://gateway"`
	Version          string `json:"version" default:"0.0.1"`
	WechatTemplateID string `json:"wechat_template_id"`
	WechatAppID      string `json:"wechat_app_id"`
	WechatSecret     string `json:"wechat_secret"`
}

func NewConfig() (Config, error) {
	godotenv.Load()
	var config Config
	err := envconfig.Process("", &config)
	return config, err
}
