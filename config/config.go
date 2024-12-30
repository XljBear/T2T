package config

import (
	"os"

	"github.com/spf13/viper"
)

type ProxyAddressRecord struct {
	UUID          string `mapstructure:"uuid" yaml:"uuid" json:"uuid"`
	LocalAddress  string `mapstructure:"local_address" yaml:"local_address" json:"local_address"`
	RemoteAddress string `mapstructure:"remote_address" yaml:"remote_address" json:"remote_address"`
	Name          string `mapstructure:"name" yaml:"name" json:"name"`
	Status        bool   `mapstructure:"status" yaml:"status" json:"status"`
	MaxLink       uint   `mapstructure:"max_link" yaml:"max_link" json:"max_link"`
}
type Config struct {
	Proxy              []ProxyAddressRecord `mapstructure:"proxy" yaml:"proxy" json:"proxy"`
	EnablePanel        bool                 `mapstructure:"enable_panel" yaml:"enable_panel" json:"enable_panel"`
	PanelListenAddress string               `mapstructure:"panel_listen_address" yaml:"panel_listen_address" json:"panel_listen_address"`
	PanelPassword      string               `mapstructure:"panel_password" yaml:"panel_password" json:"panel_password"`
}

const configFileName = "proxy"

var Cfg = Config{}

func Init() {
	Cfg = Config{}
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

	paList := []ProxyAddressRecord{}
	viper.Set("proxy", paList)
	viper.Set("enable_panel", true)
	viper.Set("panel_listen_address", ":8080")
	viper.Set("panel_password", "admin")
	err = viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}
