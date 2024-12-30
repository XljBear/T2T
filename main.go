package main

import (
	"T2T/config"
	"T2T/panelServer"
	"T2T/proxyServer"
	"T2T/storages"
	"fmt"
	"time"
)

func main() {

	config.Init()
	storages.Init()
	defer storages.Release()

	proxySuccess := proxyServer.StartProxyServer()

	var panelSuccess bool
	if config.Cfg.EnablePanel {
		panelSuccess = panelServer.StartPanelServer(config.Cfg.PanelListenAddress)
		if panelSuccess {
			fmt.Println("Panel server started successfully")
		} else {
			fmt.Println("Panel server failed to start")
		}
	}

	if !proxySuccess && !panelSuccess {
		return
	}

	for {
		time.Sleep(time.Millisecond * 100)
	}
}
