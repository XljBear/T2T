package handles

import (
	"T2T/config"
	"T2T/panelServer/structs"
	"T2T/proxyServer"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"sort"
)

func GetProxyList(ctx *gin.Context) {
	ctx.JSON(200, config.Cfg.Proxy)
}
func CreateProxy(ctx *gin.Context) {
	proxyData := config.ProxyAddressRecord{}
	err := ctx.BindJSON(&proxyData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if proxyData.LocalAddress == "" || proxyData.RemoteAddress == "" || proxyData.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "LocalAddress, RemoteAddress and Name are required"})
		return
	}
	for _, proxy := range config.Cfg.Proxy {
		if proxy.LocalAddress == proxyData.LocalAddress && proxy.RemoteAddress == proxyData.RemoteAddress {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Proxy already exists"})
			return
		}
	}
	proxyData.UUID = uuid.New().String()
	err = config.SaveProxy()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{})
}

func UpdateProxy(ctx *gin.Context) {
	uuidStr := ctx.Param("uuid")
	if uuidStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid uuid"})
		return
	}
	updateData := config.ProxyAddressRecord{}
	err := ctx.BindJSON(&updateData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if updateData.LocalAddress == "" || updateData.RemoteAddress == "" || updateData.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "LocalAddress, RemoteAddress and Name are required"})
		return
	}
	proxy := config.FindProxyByUUID(uuidStr)
	if proxy == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid uuid"})
		return
	}
	proxy.Name = updateData.Name
	proxy.LocalAddress = updateData.LocalAddress
	proxy.RemoteAddress = updateData.RemoteAddress
	proxy.MaxLink = updateData.MaxLink
	proxy.Status = updateData.Status
	err = config.SaveProxy()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{})
}

func DeleteProxy(ctx *gin.Context) {
	deleteUUIDStr := ctx.Param("uuid")
	if deleteUUIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid uuid"})
		return
	}
	success := config.DeleteProxyByUUID(deleteUUIDStr)
	if !success {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid uuid"})
		return
	}
	ctx.JSON(200, gin.H{})
}

func RestartProxyServer(ctx *gin.Context) {
	config.Init()
	proxyServer.ProxyServerInstance.Start()
	ctx.JSON(200, gin.H{})
}

func GetTraffic(ctx *gin.Context) {
	type trafficData struct {
		DownlinkInSecond uint64 `json:"downlink_in_second"`
		DownlinkTotal    uint64 `json:"downlink_total"`
		UplinkInSecond   uint64 `json:"uplink_in_second"`
		UplinkTotal      uint64 `json:"uplink_total"`
		LinkCount        uint   `json:"link_count"`
	}
	var traffic map[string]trafficData
	traffic = make(map[string]trafficData)
	for UUID, proxy := range proxyServer.ProxyServerInstance.ProxyManager {
		traffic[UUID] = trafficData{
			DownlinkInSecond: proxy.Traffic.DownlinkInSecond,
			DownlinkTotal:    proxy.Traffic.DownlinkTotal,
			UplinkInSecond:   proxy.Traffic.UplinkInSecond,
			UplinkTotal:      proxy.Traffic.UplinkTotal,
			LinkCount:        proxy.LinksCount,
		}
	}
	ctx.JSON(200, traffic)
}
func GetLinks(ctx *gin.Context) {

	var links []structs.Link

	uuidStr := ctx.Param("uuid")
	if uuidStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid uuid"})
		return
	}
	if _, ok := proxyServer.ProxyServerInstance.ProxyManager[uuidStr]; !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid uuid"})
		return
	}

	proxy := proxyServer.ProxyServerInstance.ProxyManager[uuidStr]
	links = make([]structs.Link, 0)
	if proxy != nil {
		for linkUUID, link := range proxy.Links.Range {
			linkData := structs.Link{
				UUID:     linkUUID.(string),
				IP:       link.(*proxyServer.Link).RemoteIP,
				LinkTime: link.(*proxyServer.Link).Start,
				Traffic: &structs.TrafficData{
					DownlinkInSecond: link.(*proxyServer.Link).Traffic.DownlinkInSecond,
					DownlinkTotal:    link.(*proxyServer.Link).Traffic.DownlinkTotal,
					UplinkInSecond:   link.(*proxyServer.Link).Traffic.UplinkInSecond,
					UplinkTotal:      link.(*proxyServer.Link).Traffic.UplinkTotal,
				},
			}
			links = append(links, linkData)
		}
	}
	sort.Sort(structs.ByLinkTime(links))
	ctx.JSON(200, links)
}

func KickProxyServer(ctx *gin.Context) {
	proxyUUIDStr := ctx.Param("uuid")
	if proxyUUIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid uuid"})
		return
	}
	if _, ok := proxyServer.ProxyServerInstance.ProxyManager[proxyUUIDStr]; !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid uuid"})
		return
	}
	linkUUIDStr := ctx.Param("link_uuid")
	if linkUUIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid link uuid"})
		return
	}
	if _, ok := proxyServer.ProxyServerInstance.ProxyManager[proxyUUIDStr].Links.Load(linkUUIDStr); !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid link uuid"})
		return
	}
	proxy, ok := proxyServer.ProxyServerInstance.ProxyManager[proxyUUIDStr].Links.Load(linkUUIDStr)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}
	proxy.(*proxyServer.Link).Close()
	ctx.JSON(http.StatusOK, gin.H{})
}
