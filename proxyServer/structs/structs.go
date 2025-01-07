package structs

import (
	"github.com/google/uuid"
	"net"
	"sync"
	"time"
)

type TrafficMonitor struct {
	Lock          sync.Mutex
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
			tm.Lock.Lock()
			tm.DownlinkInSecond = tm.DownlinkRecord
			tm.UplinkInSecond = tm.UplinkRecord
			tm.DownlinkRecord = 0
			tm.UplinkRecord = 0
			tm.Lock.Unlock()
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
	tm.Lock.Lock()
	tm.DownlinkRecord += traffic
	tm.DownlinkTotal += traffic
	tm.Lock.Unlock()
	if tm.ParentTrafficMonitor != nil {
		tm.ParentTrafficMonitor.Downlink(traffic)
	}
}
func (tm *TrafficMonitor) Uplink(traffic uint64) {
	tm.Lock.Lock()
	tm.UplinkRecord += traffic
	tm.UplinkTotal += traffic
	tm.Lock.Unlock()
	if tm.ParentTrafficMonitor != nil {
		tm.ParentTrafficMonitor.Uplink(traffic)
	}
}

type Link struct {
	UUID       string
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
	UUID           string
	Name           string
	LocalAddress   string
	RemoteAddress  string
	Listener       net.Listener
	MaxLink        uint
	Links          sync.Map
	LinksCount     uint
	LinksCountLock sync.Mutex
	Traffic        *TrafficMonitor
}

func (proxy *Proxy) AddLink(localConn net.Conn, remoteAddr string) *Link {
	uid := uuid.New().String()
	link := Link{
		UUID:     uid,
		Start:    time.Now(),
		Conn:     localConn,
		RemoteIP: remoteAddr,
		Proxy:    proxy,
		Traffic:  &TrafficMonitor{},
	}
	proxy.Links.Store(uid, &link)
	proxy.LinksCountLock.Lock()
	proxy.LinksCount++
	proxy.LinksCountLock.Unlock()
	link.Traffic.ParentTrafficMonitor = proxy.Traffic
	link.Traffic.Start()
	exitSignal := make(chan bool)
	link.ExitSignal = &exitSignal
	return &link
}
func (proxy *Proxy) ReleaseLink(link *Link) {
	link.Traffic.Stop()
	proxy.Links.Delete(link.UUID)
	proxy.LinksCountLock.Lock()
	proxy.LinksCount--
	proxy.LinksCountLock.Unlock()
}

type ProxyManager map[string]*Proxy

func (pm *ProxyManager) CloseAllProxy() {
	for _, proxy := range *pm {
		for _, conn := range proxy.Links.Range {
			conn.(*Link).Traffic.Stop()
			_ = conn.(*Link).Conn.Close()
		}
		proxy.Traffic.Stop()
		_ = proxy.Listener.Close()
	}
	*pm = make(map[string]*Proxy)
}

type MiddleWare func(proxy *Proxy, conn net.Conn) error

func (proxy *Proxy) ProxyMiddlewareChain(conn net.Conn, middlewares ...MiddleWare) error {
	for _, middleware := range middlewares {
		if err := middleware(proxy, conn); err != nil {
			return err
		}
	}
	return nil
}
