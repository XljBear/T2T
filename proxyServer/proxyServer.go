package proxyServer

import (
	"T2T/config"
	"T2T/proxyServer/middlewares"
	"T2T/proxyServer/structs"
	"log"
	"net"
	"sync"
	"time"
)

type ProxyServer struct {
	ProxyManager structs.ProxyManager
	StopSignal   chan bool
}

var ProxyServerInstance *ProxyServer

func handleConnection(proxy *structs.Proxy, localConn net.Conn, remoteAddr string) {
	defer localConn.Close()
	remoteConn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Printf("[%s]Error connecting to remote server %s: %s\n", proxy.Name, remoteAddr, err)
		return
	}
	remoteTcpAddr := localConn.RemoteAddr().(*net.TCPAddr)
	link := proxy.AddLink(localConn, remoteTcpAddr.IP.String())
	brokenSignal := make(chan bool)
	go proxyTransform(remoteConn, localConn, link.Traffic, "Downlink", brokenSignal)
	go proxyTransform(localConn, remoteConn, link.Traffic, "Uplink", brokenSignal)
	select {
	case <-brokenSignal:
	case <-*link.ExitSignal:
	}
	proxy.ReleaseLink(link)
	remoteConn.Close()
}
func proxyTransform(dst net.Conn, src net.Conn, trafficMonitor *structs.TrafficMonitor, trafficType string, brokenSignal chan bool) {
	buffer := make([]byte, 4096)
	for {
		bytesRead, err := src.Read(buffer)
		if err != nil {
			brokenSignal <- true
			return
		}
		if trafficType == "Downlink" {
			go trafficMonitor.Downlink(uint64(bytesRead))
		} else if trafficType == "Uplink" {
			go trafficMonitor.Uplink(uint64(bytesRead))
		}
		_, err = dst.Write(buffer[:bytesRead])
		if err != nil {
			brokenSignal <- true
			return
		}
	}
}

func (ps *ProxyServer) Start() (success bool) {
	if len(ps.ProxyManager) != 0 {
		log.Println("Restarting proxy server.")
		ps.Stop()
	} else {
		ps.ProxyManager = make(map[string]*structs.Proxy)
		log.Println("Starting proxy server.")
	}
	ps.StopSignal = make(chan bool)
	if len(config.Cfg.Proxy) == 0 {
		log.Println("No proxy configured, exiting.")
		return false
	}
	for _, proxyAddressRecord := range config.Cfg.Proxy {
		if !proxyAddressRecord.Status {
			continue
		}
		proxy := structs.Proxy{
			UUID:          proxyAddressRecord.UUID,
			Name:          proxyAddressRecord.Name,
			LocalAddress:  proxyAddressRecord.LocalAddress,
			RemoteAddress: proxyAddressRecord.RemoteAddress,
			Links:         sync.Map{},
			MaxLink:       proxyAddressRecord.MaxLink,
			Traffic: &structs.TrafficMonitor{
				DownlinkTotal: proxyAddressRecord.TotalDownlink,
				UplinkTotal:   proxyAddressRecord.TotalUplink,
			},
		}
		log.Printf("[%s]Proxying %s to %s\n", proxy.Name, proxy.LocalAddress, proxy.RemoteAddress)
		listener, err := net.Listen("tcp", proxy.LocalAddress)
		if err != nil {
			log.Println("Error listening on", proxy.LocalAddress, err)
			continue
		}
		proxy.Listener = listener
		ps.ProxyManager[proxy.UUID] = &proxy
		proxy.Traffic.Start()
		log.Printf("[%s]Proxying started on %s\n", proxy.Name, proxy.LocalAddress)
		go func(proxy *structs.Proxy) {
			for {
				localConn, err := listener.Accept()
				if err != nil {
					return
				}
				err = proxy.ProxyMiddlewareChain(localConn, middlewares.IPAllowBlock, middlewares.MaxLinks)
				if err != nil {
					_ = localConn.Close()
					log.Println(err.Error())
					continue
				}
				go handleConnection(proxy, localConn, proxy.RemoteAddress)
			}
		}(&proxy)
	}
	go func() {
		for {
			ps.writeProxyTraffic()
			select {
			case <-ps.StopSignal:
				log.Println("Proxy server stopped.")
				return
			case <-time.After(time.Second * 10):
				continue
			}
		}
	}()
	if len(ps.ProxyManager) == 0 {
		log.Println("No proxy configured, exiting.")
		return false
	}
	return true
}
func (ps *ProxyServer) Stop() {
	ps.ProxyManager.CloseAllProxy()
	ps.StopSignal <- true
}
func (ps *ProxyServer) writeProxyTraffic() {
	if len(ps.ProxyManager) == 0 {
		return
	}
	for _, proxy := range ps.ProxyManager {
		proxyCfg := config.FindProxyByUUID(proxy.UUID)
		if proxyCfg == nil {
			continue
		}
		proxyCfg.TotalDownlink = proxy.Traffic.DownlinkTotal
		proxyCfg.TotalUplink = proxy.Traffic.UplinkTotal
		_ = config.SaveProxy()
	}
}
