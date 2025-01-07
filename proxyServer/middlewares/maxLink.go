package middlewares

import (
	"T2T/proxyServer/structs"
	"errors"
	"fmt"
	"net"
)

func MaxLinks(proxy *structs.Proxy, conn net.Conn) error {
	if proxy.MaxLink > 0 && proxy.LinksCount >= proxy.MaxLink {
		return errors.New(fmt.Sprintf("[%s]Max link reached, rejecting connection.", proxy.Name))
	}
	return nil
}
