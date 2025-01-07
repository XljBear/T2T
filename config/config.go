package config

import (
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type ProxyAddressRecord struct {
	UUID          string `mapstructure:"uuid" yaml:"uuid" json:"uuid"`
	LocalAddress  string `mapstructure:"local_address" yaml:"local_address" json:"local_address"`
	RemoteAddress string `mapstructure:"remote_address" yaml:"remote_address" json:"remote_address"`
	Name          string `mapstructure:"name" yaml:"name" json:"name"`
	Status        bool   `mapstructure:"status" yaml:"status" json:"status"`
	MaxLink       uint   `mapstructure:"max_link" yaml:"max_link" json:"max_link"`
	TotalDownlink uint64 `mapstructure:"total_downlink" yaml:"total_downlink" json:"total_downlink"`
	TotalUplink   uint64 `mapstructure:"total_uplink" yaml:"total_uplink" json:"total_uplink"`
}
type Config struct {
	Proxy              []ProxyAddressRecord `mapstructure:"proxy" yaml:"proxy" json:"proxy"`
	EnablePanel        bool                 `mapstructure:"enable_panel" yaml:"enable_panel" json:"enable_panel"`
	PanelListenAddress string               `mapstructure:"panel_listen_address" yaml:"panel_listen_address" json:"panel_listen_address"`
	PanelPassword      string               `mapstructure:"panel_password" yaml:"panel_password" json:"panel_password"`
	CaptchaType        uint                 `mapstructure:"captcha_type" yaml:"captcha_type" json:"captcha_type"`
	DarkMode           bool                 `mapstructure:"dark_mode" yaml:"dark_mode" json:"dark_mode"`
}

const configFileName = "proxy"

var Cfg = Config{}
var cfgLock sync.Mutex

func InitConfig() {
	cfgLock.Lock()
	defer cfgLock.Unlock()
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
	log.Println("Config loaded.")
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
	viper.Set("captcha_type", 1)
	viper.Set("dark_mode", false)
	err = viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

func SaveProxy() error {
	cfgLock.Lock()
	defer cfgLock.Unlock()
	viper.Set("proxy", Cfg.Proxy)
	err := viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

func FindProxyByUUID(uuid string) *ProxyAddressRecord {
	for i, proxy := range Cfg.Proxy {
		if proxy.UUID == uuid {
			return &Cfg.Proxy[i]
		}
	}
	return nil
}

func findProxyIndex(uuid string) int {
	for i, proxy := range Cfg.Proxy {
		if proxy.UUID == uuid {
			return i
		}
	}
	return -1
}

func DeleteProxyByUUID(uuid string) bool {
	deleteIndex := findProxyIndex(uuid)
	if deleteIndex < 0 || deleteIndex >= len(Cfg.Proxy) {
		return false
	}
	Cfg.Proxy = append(Cfg.Proxy[:deleteIndex], Cfg.Proxy[deleteIndex+1:]...)
	err := SaveProxy()
	if err != nil {
		return false
	}
	return true
}

func ReloadConfig() error {
	InitConfig()
	log.Println("Config reloaded.")
	return nil
}
