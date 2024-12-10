package proxyServer

import (
	"T2T/config"
	"fmt"
	"io"
	"net"
)

var listenerList []net.Listener
var connList []net.Conn

func handleConnection(localConn net.Conn, remoteAddr string) {
	defer localConn.Close()
	remoteConn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		panic(err)
	}
	defer remoteConn.Close()
	go io.Copy(remoteConn, localConn)
	io.Copy(localConn, remoteConn)
}
func StartProxyServer() (success bool) {
	if len(listenerList) != 0 {
		fmt.Println("Restarting T2T server")
		for _, conn := range connList {
			conn.Close()
		}
		for _, listener := range listenerList {
			listener.Close()
		}
		connList = []net.Conn{}
		listenerList = []net.Listener{}
	} else {
		fmt.Println("Starting T2T server")
	}
	if len(config.Cfg.Proxy) == 0 {
		fmt.Println("No proxy configured, exiting")
		return false
	}
	for _, proxy := range config.Cfg.Proxy {
		if !proxy.Status {
			continue
		}
		fmt.Printf("[%s]Proxying %s to %s\n", proxy.Name, proxy.LocalAddress, proxy.RemoteAddress)
		listener, err := net.Listen("tcp", proxy.LocalAddress)
		if err != nil {
			fmt.Println("Error listening on", proxy.LocalAddress, err)
			continue
		}
		listenerList = append(listenerList, listener)
		fmt.Printf("[%s]Proxying started on %s\n", proxy.Name, proxy.LocalAddress)
		go func() {
			for {
				localConn, err := listener.Accept()
				connList = append(connList, localConn)
				if err != nil {
					return
				}
				go handleConnection(localConn, proxy.RemoteAddress)
			}
		}()
	}
	if len(listenerList) == 0 {
		fmt.Println("No proxy configured, exiting")
		return false
	}
	return true
}
