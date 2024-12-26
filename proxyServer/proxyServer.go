package proxyServer

import (
	"T2T/config"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net"
	"time"
)

//var listenerList []net.Listener
//var connList []net.Conn

type Link struct {
	Conn      net.Conn
	Start     time.Time
	InBounds  uint64
	OutBounds uint64
}
type ProxyInfo struct {
	Listener  net.Listener
	MaxLink   uint
	Links     map[string]Link
	InBounds  uint64
	OutBounds uint64
}

var proxyManager map[string]ProxyInfo

func handleConnection(proxyInfo *ProxyInfo, localConn net.Conn, remoteAddr string) {
	defer localConn.Close()
	remoteConn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		panic(err)
	}
	link := Link{
		Start: time.Now(),
		Conn:  localConn,
	}
	uid := uuid.New().String()
	proxyInfo.Links[uid] = link

	var outBounds int64
	var inBounds int64
	go func() {
		inBounds, _ = io.Copy(remoteConn, localConn)
	}()
	outBounds, _ = io.Copy(localConn, remoteConn)

	proxyInfo.OutBounds += uint64(outBounds)
	proxyInfo.InBounds += uint64(inBounds)
	delete(proxyInfo.Links, uid)
	remoteConn.Close()
}
func CloseAllProxy() {
	for _, proxy := range proxyManager {
		for _, conn := range proxy.Links {
			_ = conn.Conn.Close()
		}
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
		}
		fmt.Printf("[%s]Proxying %s to %s\n", proxyAddressRecord.Name, proxyAddressRecord.LocalAddress, proxyAddressRecord.RemoteAddress)
		listener, err := net.Listen("tcp", proxyAddressRecord.LocalAddress)
		if err != nil {
			fmt.Println("Error listening on", proxyAddressRecord.LocalAddress, err)
			continue
		}
		proxyInfo.Listener = listener
		proxyManager[proxyAddressRecord.Name] = proxyInfo
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
