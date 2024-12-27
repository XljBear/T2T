package proxyServer

import (
	"T2T/config"
	"fmt"
	"github.com/google/uuid"
	"net"
	"sync"
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
	tm.BreakSignal = make(chan bool)
	go func() {
		for {
			tm.DownlinkInSecond = tm.DownlinkRecord
			tm.UplinkInSecond = tm.UplinkRecord
			tm.DownlinkRecord = 0
			tm.UplinkRecord = 0
			select {
			case <-tm.BreakSignal:
				return
			case <-time.After(time.Second):
				continue
			}
		}
	}()
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
	Conn       net.Conn
	Start      time.Time
	RemoteIP   string
	Proxy      *Proxy
	Traffic    *TrafficMonitor
	ExitSignal *chan bool
}

func (link *Link) Close() {
	*link.ExitSignal <- true
}

type Proxy struct {
	UUID          string
	Name          string
	LocalAddress  string
	RemoteAddress string
	Listener      net.Listener
	MaxLink       uint
	Links         sync.Map
	LinksCount    uint
	Traffic       *TrafficMonitor
}

var ProxyManager map[string]*Proxy

func handleConnection(proxy *Proxy, localConn net.Conn, remoteAddr string) {
	defer localConn.Close()
	remoteConn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		panic(err)
	}
	remoteTcpAddr := localConn.RemoteAddr().(*net.TCPAddr)
	link := Link{
		Start:    time.Now(),
		Conn:     localConn,
		RemoteIP: remoteTcpAddr.IP.String(),
		Proxy:    proxy,
		Traffic:  &TrafficMonitor{},
	}
	uid := uuid.New().String()
	proxy.Links.Store(uid, &link)
	proxy.LinksCount++
	link.Traffic.ParentTrafficMonitor = proxy.Traffic
	link.Traffic.Start()
	brokenSignal := make(chan bool, 1)
	exitSignal := make(chan bool, 1)
	link.ExitSignal = &exitSignal
	go proxyTransform(remoteConn, localConn, link.Traffic, "Downlink", brokenSignal)
	go proxyTransform(localConn, remoteConn, link.Traffic, "Uplink", brokenSignal)
	select {
	case <-brokenSignal:
	case <-exitSignal:
	}
	link.Traffic.Stop()
	proxy.Links.Delete(uid)
	proxy.LinksCount--
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
	for _, proxy := range ProxyManager {
		for _, conn := range proxy.Links.Range {
			conn.(*Link).Traffic.Stop()
			_ = conn.(*Link).Conn.Close()
		}
		proxy.Traffic.Stop()
		_ = proxy.Listener.Close()
	}
	ProxyManager = make(map[string]*Proxy)
}
func StartProxyServer() (success bool) {
	if len(ProxyManager) != 0 {
		fmt.Println("Restarting T2T server")
		CloseAllProxy()
	} else {
		ProxyManager = make(map[string]*Proxy)
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
		proxy := Proxy{
			UUID:          proxyAddressRecord.UUID,
			Name:          proxyAddressRecord.Name,
			LocalAddress:  proxyAddressRecord.LocalAddress,
			RemoteAddress: proxyAddressRecord.RemoteAddress,
			Links:         sync.Map{},
			MaxLink:       proxyAddressRecord.MaxLink,
			Traffic:       &TrafficMonitor{},
		}
		fmt.Printf("[%s]Proxying %s to %s\n", proxy.Name, proxy.LocalAddress, proxy.RemoteAddress)
		listener, err := net.Listen("tcp", proxy.LocalAddress)
		if err != nil {
			fmt.Println("Error listening on", proxy.LocalAddress, err)
			continue
		}
		proxy.Listener = listener
		ProxyManager[proxy.UUID] = &proxy
		proxy.Traffic.Start()
		fmt.Printf("[%s]Proxying started on %s\n", proxy.Name, proxy.LocalAddress)
		go func(proxy *Proxy) {
			for {
				localConn, err := listener.Accept()
				if err != nil {
					return
				}
				if proxy.MaxLink > 0 && proxy.LinksCount >= proxy.MaxLink {
					_ = localConn.Close()
					fmt.Printf("[%s]Max link reached, rejecting connection\n", proxy.Name)
					continue
				}
				go handleConnection(proxy, localConn, proxy.RemoteAddress)
			}
		}(&proxy)
	}
	if len(ProxyManager) == 0 {
		fmt.Println("No proxy configured, exiting")
		return false
	}
	return true
}
