package main

import (
	"T2T/config"
	"T2T/panelServer"
	"T2T/proxyServer"
	"time"
)

func main() {

	config.Init()

	success := proxyServer.StartProxyServer()
	if config.Cfg.EnablePanel {
		panelServer.StartPanelServer(config.Cfg.PanelListenAddress)
	} else if !success {
		return
	}

	for {
		time.Sleep(time.Millisecond)
	}
}
