package main

import (
	"T2T/config"
	"T2T/panelServer"
	"T2T/proxyServer"
	"T2T/storages"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("Starting T2T server...")
	config.InitConfig()
	config.InitBlockIPs()
	storages.Init()

	systemSignal := make(chan os.Signal, 1)
	signal.Notify(systemSignal, syscall.SIGINT, syscall.SIGTERM)

	proxyServer.ProxyServerInstance = &proxyServer.ProxyServer{}
	proxySuccess := proxyServer.ProxyServerInstance.Start()

	var panelSuccess bool
	if config.Cfg.EnablePanel {
		panelSuccess = panelServer.StartPanelServer(config.Cfg.PanelListenAddress)
		if panelSuccess {
			log.Println("Panel server started successfully.")
		} else {
			log.Println("Panel server failed to start.")
		}
	}

	if !proxySuccess && !panelSuccess {
		return
	}

	for {
		select {
		case <-systemSignal:
			log.Println("Received system signal, exiting...")
			config.StopIPCleaner()
			storages.Release()
			proxyServer.ProxyServerInstance.Stop()
			log.Println("T2T server stopped.")
			return
		case <-time.After(time.Millisecond * 100):
		}
	}
}
