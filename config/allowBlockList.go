package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"sync"
	"time"
)

type IPItem struct {
	UUID      string     `mapstructure:"uuid" yaml:"uuid" json:"uuid"`
	IP        string     `mapstructure:"ip" yaml:"ip" json:"ip"`
	Port      string     `mapstructure:"port" yaml:"port" json:"port"` // empty means all ports
	StartTime time.Time  `mapstructure:"start_time" yaml:"start_time" json:"start_time"`
	EndTime   *time.Time `mapstructure:"end_time" yaml:"end_time" json:"end_time"`
	Reason    string     `mapstructure:"reason" yaml:"reason" json:"reason"`
}
type Block struct {
	BlockIPs []IPItem `mapstructure:"block_ips" yaml:"block_ips" json:"block_ips"`
}
type Allow struct {
	AllowIPs []IPItem `mapstructure:"allow_ips" yaml:"allow_ips" json:"allow_ips"`
}

type AllowBlock struct {
	Block     Block      `mapstructure:"block" yaml:"block" json:"block"`
	BlockLock sync.Mutex `mapstructure:"-" yaml:"-" json:"-"`
	Allow     Allow      `mapstructure:"allow" yaml:"allow" json:"allow"`
	AllowLock sync.Mutex `mapstructure:"-" yaml:"-" json:"-"`
	Mode      int        `mapstructure:"mode" yaml:"mode" json:"mode"` // 0: disable, 1: block mode, 2: allow mode
}
type ABCfg struct {
	AllowBlock AllowBlock `mapstructure:"allow-block" yaml:"allow-block" json:"allow-block"`
}

const allowBlockFileName = "allow-block"

var AllowBlockCfg *ABCfg
var allowBlockViper *viper.Viper
var allowBlockLock = sync.Mutex{}
var ipCleaner *IPCleaner

func InitBlockIPs() {
	allowBlockLock.Lock()
	defer allowBlockLock.Unlock()
	if ipCleaner == nil {
		ipCleaner = &IPCleaner{}
		ipCleaner.ExitSignal = make(chan bool)
		ipCleaner.Start()
	}
	AllowBlockCfg = &ABCfg{
		AllowBlock: AllowBlock{},
	}
	allowBlockViper = viper.New()
	allowBlockViper.SetConfigName(allowBlockFileName)
	allowBlockViper.SetConfigType("yaml")
	allowBlockViper.AddConfigPath(".")
	if err := allowBlockViper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err := CreateDefaultAllowBlock()
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	allowBlockViper.Get("")
	e := allowBlockViper.Unmarshal(&AllowBlockCfg)
	if e != nil {
		panic(e)
	}
	log.Println("IP Ruler initialized.")
}
func StopIPCleaner() {
	ipCleaner.Stop()
}

func CreateDefaultAllowBlock() error {
	f, err := os.Create("./" + allowBlockFileName + ".yaml")
	if err != nil {
		return err
	}
	f.Close()

	allowBlockViper.Set("allow-block", AllowBlockCfg)
	err = allowBlockViper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

func SaveAllowBlock() error {
	allowBlockLock.Lock()
	defer allowBlockLock.Unlock()
	allowBlockViper.Set("allow-block", &AllowBlockCfg.AllowBlock)
	err := allowBlockViper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

func (ab *AllowBlock) GetAllowIPs() []IPItem {
	return ab.Allow.AllowIPs
}
func (ab *AllowBlock) GetBlockIPs() []IPItem {
	return ab.Block.BlockIPs
}

func (ab *AllowBlock) GetAllowIPsByUUID(uuid string) *IPItem {
	for _, ip := range ab.Allow.AllowIPs {
		if ip.UUID == uuid {
			return &ip
		}
	}
	return nil
}
func (ab *AllowBlock) GetBlockIPsByUUID(uuid string) *IPItem {
	for _, ip := range ab.Block.BlockIPs {
		if ip.UUID == uuid {
			return &ip
		}
	}
	return nil
}
func (ab *AllowBlock) DeleteAllowIPByUUID(uuid string) {
	ab.AllowLock.Lock()
	defer ab.AllowLock.Unlock()
	for index, ip := range ab.Allow.AllowIPs {
		if ip.UUID == uuid {
			ab.Allow.AllowIPs = append(ab.Allow.AllowIPs[:index], ab.Allow.AllowIPs[index+1:]...)
			go SaveAllowBlock()
			return
		}
	}
}

func (ab *AllowBlock) DeleteBlockIPByUUID(uuid string) {
	ab.BlockLock.Lock()
	defer ab.BlockLock.Unlock()
	for index, ip := range ab.Block.BlockIPs {
		if ip.UUID == uuid {
			ab.Block.BlockIPs = append(ab.Block.BlockIPs[:index], ab.Block.BlockIPs[index+1:]...)
			go SaveAllowBlock()
			return
		}
	}
}

func (ab *AllowBlock) GetAllowIPsByPort(port string) []IPItem {
	var allowIPs []IPItem
	for _, ip := range ab.Allow.AllowIPs {
		if ip.Port == port || ip.Port == "" {
			allowIPs = append(allowIPs, ip)
		}
	}
	return allowIPs
}

func (ab *AllowBlock) GetBlockIPsByPort(port string) []IPItem {
	var blockIPs []IPItem
	for _, ip := range ab.Block.BlockIPs {
		if ip.Port == port || ip.Port == "" {
			blockIPs = append(blockIPs, ip)
		}
	}
	return blockIPs
}

type IPCleaner struct {
	ExitSignal chan bool
}

func (ic *IPCleaner) Start() {
	go func() {
		for {
			select {
			case <-ic.ExitSignal:
				return
			case <-time.After(time.Second):
			}
			allowIps := AllowBlockCfg.AllowBlock.GetAllowIPs()
			for _, ip := range allowIps {
				if ip.EndTime != nil && time.Now().After(*ip.EndTime) {
					AllowBlockCfg.AllowBlock.DeleteAllowIPByUUID(ip.UUID)
				}
			}
			blockIps := AllowBlockCfg.AllowBlock.GetBlockIPs()
			for _, ip := range blockIps {
				if ip.EndTime != nil && time.Now().After(*ip.EndTime) {
					AllowBlockCfg.AllowBlock.DeleteBlockIPByUUID(ip.UUID)
				}
			}
		}
	}()
}
func (ic *IPCleaner) Stop() {
	ic.ExitSignal <- true
	log.Println("IP cleaner stopped.")
}
