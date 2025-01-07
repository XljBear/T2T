package middlewares

import (
	"T2T/config"
	"T2T/proxyServer/structs"
	"errors"
	"fmt"
	"net"
	"strings"
)

func IPAllowBlock(proxy *structs.Proxy, conn net.Conn) error {
	proxyPort := strings.Split(proxy.RemoteAddress, ":")[1]
	connIP := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	switch config.AllowBlockCfg.AllowBlock.Mode {
	case 1:
		// Block Mode
		IPList := config.AllowBlockCfg.AllowBlock.Block.BlockIPs
		for _, ip := range IPList {
			if ip.IP == connIP && (ip.Port == "" || ip.Port == proxyPort) {
				return errors.New(fmt.Sprintf("IP %s was blocked.", connIP))
			}
		}
	case 2:
		// Allow Mode
		IPList := config.AllowBlockCfg.AllowBlock.Allow.AllowIPs
		exist := false
		for _, ip := range IPList {
			if ip.IP == connIP && (ip.Port == "" || ip.Port == proxy.RemoteAddress) {
				exist = true
				break
			}
		}
		if !exist {
			return errors.New(fmt.Sprintf("IP %s is not allowed.", connIP))
		}
	default:
		// Disable Mode
		return nil
	}
	return nil
}
