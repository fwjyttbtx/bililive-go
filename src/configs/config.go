package configs

import (
	"crypto/tls"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type TLS struct {
	Enable   bool   `yaml:"enable"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}
type RPC struct {
	Enable bool   `yaml:"enable"`
	Port   uint   `yaml:"port"`
	Path   string `yaml:"path"`
	TLS    TLS    `yaml:"tls"`
}
type Config struct {
	RPC        RPC      `yaml:"rpc"`
	LogLevel   string   `yaml:"log_level"`
	Interval   int      `yaml:"interval"`
	OutPutPath string   `yaml:"out_put_path"`
	LiveRooms  []string `yaml:"live_rooms"`
}

func VerifyConfig(config *Config) error {
	if config.Interval <= 0 {
		return errors.New(fmt.Sprintf(`the interval can not <= 0`))
	}
	if _, err := os.Stat(config.OutPutPath); err != nil {
		return errors.New(fmt.Sprintf(`the out put path: "%s" is not exist`, config.OutPutPath))
	}
	if config.RPC.Enable {
		if config.RPC.Port == 0 {
			return errors.New("rpc listen port can not be null or '0'")
		}
		if config.RPC.TLS.Enable {
			if _, err := tls.LoadX509KeyPair(config.RPC.TLS.CertFile, config.RPC.TLS.KeyFile); err != nil {
				return err
			}
		}
	}
	return nil
}

func NewConfigWithFile(configFilePath string) (*Config, error) {
	b, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can`t open file: %s", configFilePath))
	}
	config := new(Config)
	err = yaml.Unmarshal(b, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func NewConfig() *Config {
	config := new(Config)
	config.RPC.Enable = false
	config.LogLevel = "info"
	config.Interval = 30
	config.OutPutPath = "./"
	return config
}
