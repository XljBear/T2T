package config

import (
	"os"

	"github.com/spf13/viper"
)

type ProxyAddressrRecord struct {
	LocalAddress  string `mapstructure:"local_address" yaml:"local_address" json:"local_address"`
	RemoteAddress string `mapstructure:"remote_address" yaml:"remote_address" json:"remote_address"`
	Name          string `mapstructure:"name" yaml:"name" json:"name"`
	Status        bool   `mapstructure:"status" yaml:"status" json:"status"`
}
type Config struct {
	Proxy              []ProxyAddressrRecord `mapstructure:"proxy" yaml:"proxy" json:"proxy"`
	EnablePanel        bool                  `mapstructure:"enable_panel" yaml:"enable_panel" json:"enable_panel"`
	PanelListenAddress string                `mapstructure:"panel_listen_address" yaml:"panel_listen_address" json:"panel_listen_address"`
}

const configFileName = "proxy"

var Cfg = Config{}

func Init() {
	viper.SetConfigName(configFileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err := CreateDefaultConfig()
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	e := viper.Unmarshal(&Cfg)
	if e != nil {
		panic(e)
	}
}

func CreateDefaultConfig() error {
	f, err := os.Create("./" + configFileName + ".yaml")
	if err != nil {
		return err
	}
	f.Close()

	paList := []ProxyAddressrRecord{}
	viper.Set("proxy", paList)
	viper.Set("enable_panel", true)
	viper.Set("panel_listen_address", ":8080")
	err = viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}
