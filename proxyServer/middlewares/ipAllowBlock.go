package middlewares

import (
	"T2T/config"
	"T2T/proxyServer/structs"
	"errors"
	"fmt"
	"net"
)

func IPAllowBlock(proxy *structs.Proxy, conn net.Conn) error {
	IPList := []config.IPItem{}
	switch config.AllowBlockCfg.AllowBlock.Mode {
	case 1:
		// Block Mode
		IPList = config.AllowBlockCfg.AllowBlock.Block.BlockIPs
	case 2:
		// Allow Mode
		IPList = config.AllowBlockCfg.AllowBlock.Allow.AllowIPs
	default:
		// Disable Mode
		return nil
	}
	for _, ip := range IPList {
		if ip.IP == conn.RemoteAddr().(*net.TCPAddr).IP.String() && (ip.Port == "" || ip.Port == proxy.RemoteAddress) {
			return errors.New(fmt.Sprintf("IP %s is not allowed.", ip.IP))
		}
	}
	return nil
}
