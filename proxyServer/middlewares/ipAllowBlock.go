package middlewares

import (
	"T2T/config"
	"T2T/proxyServer/structs"
	"errors"
	"fmt"
	"net"
	"slices"
	"strconv"
)

func IPAllowBlock(proxy *structs.Proxy, conn net.Conn) error {
	localPort := strconv.Itoa(conn.LocalAddr().(*net.TCPAddr).Port)
	connIP := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	switch config.AllowBlockCfg.AllowBlock.Mode {
	case 1:
		// Block Mode
		IPList := config.AllowBlockCfg.AllowBlock.Block.BlockIPs
		if matchRule(connIP, localPort, IPList) {
			return errors.New(fmt.Sprintf("IP %s was blocked.", connIP))
		}
	case 2:
		// Allow Mode
		IPList := config.AllowBlockCfg.AllowBlock.Allow.AllowIPs
		if !matchRule(connIP, localPort, IPList) {
			return errors.New(fmt.Sprintf("IP %s is not allowed.", connIP))
		}
	default:
		// Disable Mode
		return nil
	}
	return nil
}
func matchRule(connIP string, port string, ipList []config.IPItem) bool {
	for _, ip := range ipList {
		if ip.IP == connIP && (len(ip.Port) == 0 || slices.Contains(ip.Port, port)) {
			return true
		}
	}
	return false
}
