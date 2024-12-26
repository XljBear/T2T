package proxyServer

import (
	"T2T/config"
	"fmt"
	"github.com/google/uuid"
	"net"
	"time"
)

type TrafficMonitor struct {
	DownlinkTotal uint64
	UplinkTotal   uint64

	DownlinkRecord   uint64
	UplinkRecord     uint64
	DownlinkInSecond uint64
	UplinkInSecond   uint64
	BreakSignal      chan bool

	ParentTrafficMonitor *TrafficMonitor
}

func (tm *TrafficMonitor) Start() {
	for {
		tm.DownlinkInSecond = tm.DownlinkRecord
		tm.UplinkInSecond = tm.UplinkRecord
		tm.DownlinkRecord = 0
		tm.UplinkRecord = 0
		//if tm.DownlinkInSecond != 0 || tm.UplinkInSecond != 0 {
		//	if tm.ParentTrafficMonitor != nil {
		//		fmt.Printf("Sub-Traffic Monitor: Downlink: %d, Uplink: %d\n", tm.DownlinkInSecond, tm.UplinkInSecond)
		//	} else {
		//		fmt.Printf("Main-Traffic Monitor: Downlink: %d, Uplink: %d\n", tm.DownlinkInSecond, tm.UplinkInSecond)
		//	}
		//}
		select {
		case <-tm.BreakSignal:
			return
		case <-time.After(time.Second):
			continue
		}
	}
}
func (tm *TrafficMonitor) Stop() {
	tm.BreakSignal <- true
}
func (tm *TrafficMonitor) Downlink(traffic uint64) {
	tm.DownlinkRecord += traffic
	tm.DownlinkTotal += traffic
	if tm.ParentTrafficMonitor != nil {
		tm.ParentTrafficMonitor.Downlink(traffic)
	}
}
func (tm *TrafficMonitor) Uplink(traffic uint64) {
	tm.UplinkRecord += traffic
	tm.UplinkTotal += traffic
	if tm.ParentTrafficMonitor != nil {
		tm.ParentTrafficMonitor.Uplink(traffic)
	}
}

type Link struct {
	Conn    net.Conn
	Start   time.Time
	Traffic *TrafficMonitor
}
type ProxyInfo struct {
	Listener net.Listener
	MaxLink  uint
	Links    map[string]Link
	Traffic  *TrafficMonitor
}

var proxyManager map[string]ProxyInfo

func handleConnection(proxyInfo *ProxyInfo, localConn net.Conn, remoteAddr string) {
	defer localConn.Close()
	remoteConn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		panic(err)
	}
	link := Link{
		Start:   time.Now(),
		Conn:    localConn,
		Traffic: &TrafficMonitor{},
	}
	uid := uuid.New().String()
	proxyInfo.Links[uid] = link
	link.Traffic.ParentTrafficMonitor = proxyInfo.Traffic
	go link.Traffic.Start()
	brokenSignal := make(chan bool)
	go proxyTransform(remoteConn, localConn, link.Traffic, "Downlink", brokenSignal)
	go proxyTransform(localConn, remoteConn, link.Traffic, "Uplink", brokenSignal)
	<-brokenSignal
	link.Traffic.Stop()
	delete(proxyInfo.Links, uid)
	remoteConn.Close()
}
func proxyTransform(dst net.Conn, src net.Conn, trafficMonitor *TrafficMonitor, trafficType string, brokenSignal chan bool) {
	buffer := make([]byte, 4096)
	for {
		bytesRead, err := src.Read(buffer)
		if err != nil {
			brokenSignal <- true
			return
		}
		if trafficType == "Downlink" {
			trafficMonitor.Downlink(uint64(bytesRead))
		} else if trafficType == "Uplink" {
			trafficMonitor.Uplink(uint64(bytesRead))
		}
		_, err = dst.Write(buffer[:bytesRead])
		if err != nil {
			brokenSignal <- true
			return
		}
	}
}
func CloseAllProxy() {
	for _, proxy := range proxyManager {
		for _, conn := range proxy.Links {
			conn.Traffic.Stop()
			_ = conn.Conn.Close()
		}
		proxy.Traffic.Stop()
		_ = proxy.Listener.Close()
	}
	proxyManager = make(map[string]ProxyInfo)
}
func StartProxyServer() (success bool) {
	if len(proxyManager) != 0 {
		fmt.Println("Restarting T2T server")
		CloseAllProxy()
	} else {
		proxyManager = make(map[string]ProxyInfo)
		fmt.Println("Starting T2T server")
	}
	if len(config.Cfg.Proxy) == 0 {
		fmt.Println("No proxy configured, exiting")
		return false
	}
	for _, proxyAddressRecord := range config.Cfg.Proxy {
		if !proxyAddressRecord.Status {
			continue
		}
		proxyInfo := ProxyInfo{
			Links:   make(map[string]Link),
			MaxLink: proxyAddressRecord.MaxLink,
			Traffic: &TrafficMonitor{},
		}
		fmt.Printf("[%s]Proxying %s to %s\n", proxyAddressRecord.Name, proxyAddressRecord.LocalAddress, proxyAddressRecord.RemoteAddress)
		listener, err := net.Listen("tcp", proxyAddressRecord.LocalAddress)
		if err != nil {
			fmt.Println("Error listening on", proxyAddressRecord.LocalAddress, err)
			continue
		}
		proxyInfo.Listener = listener
		proxyManager[proxyAddressRecord.UUID] = proxyInfo
		go proxyInfo.Traffic.Start()
		fmt.Printf("[%s]Proxying started on %s\n", proxyAddressRecord.Name, proxyAddressRecord.LocalAddress)
		go func(proxyAddressRecord *config.ProxyAddressRecord, proxyInfo *ProxyInfo) {
			for {
				localConn, err := listener.Accept()
				if err != nil {
					return
				}
				if proxyInfo.MaxLink > 0 && uint(len(proxyInfo.Links)) >= proxyInfo.MaxLink {
					_ = localConn.Close()
					fmt.Printf("[%s]Max link reached, rejecting connection\n", proxyAddressRecord.Name)
					continue
				}
				go handleConnection(proxyInfo, localConn, proxyAddressRecord.RemoteAddress)
			}
		}(&proxyAddressRecord, &proxyInfo)
	}
	if len(proxyManager) == 0 {
		fmt.Println("No proxy configured, exiting")
		return false
	}
	return true
}
