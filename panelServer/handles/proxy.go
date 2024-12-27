package handles

import (
	"T2T/config"
	"T2T/proxyServer"
	"T2T/proxyServer/structs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"sort"
)

func GetProxyList(ctx *gin.Context) {
	ctx.JSON(200, config.Cfg.Proxy)
}
func CreateProxy(ctx *gin.Context) {
	proxyData := config.ProxyAddressRecord{}
	err := ctx.BindJSON(&proxyData)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if proxyData.LocalAddress == "" || proxyData.RemoteAddress == "" || proxyData.Name == "" {
		ctx.JSON(400, gin.H{"error": "LocalAddress, RemoteAddress and Name are required"})
		return
	}
	for _, proxy := range config.Cfg.Proxy {
		if proxy.LocalAddress == proxyData.LocalAddress && proxy.RemoteAddress == proxyData.RemoteAddress {
			ctx.JSON(400, gin.H{"error": "Proxy already exists"})
			return
		}
	}
	proxyData.UUID = uuid.New().String()
	config.Cfg.Proxy = append(config.Cfg.Proxy, proxyData)
	viper.Set("proxy", config.Cfg.Proxy)
	err = viper.WriteConfig()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{})
}

func UpdateProxy(ctx *gin.Context) {
	uuidStr := ctx.Param("uuid")
	if uuidStr == "" {
		ctx.JSON(400, gin.H{"error": "Invalid uuid"})
		return
	}
	updateData := config.ProxyAddressRecord{}
	err := ctx.BindJSON(&updateData)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if updateData.LocalAddress == "" || updateData.RemoteAddress == "" || updateData.Name == "" {
		ctx.JSON(400, gin.H{"error": "LocalAddress, RemoteAddress and Name are required"})
		return
	}
	updateIndex := findProxyIndex(uuidStr)
	if updateIndex < 0 || updateIndex >= len(config.Cfg.Proxy) {
		ctx.JSON(400, gin.H{"error": "Invalid uuid"})
		return
	}
	config.Cfg.Proxy[updateIndex] = updateData
	viper.Set("proxy", config.Cfg.Proxy)
	err = viper.WriteConfig()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{})
}

func DeleteProxy(ctx *gin.Context) {
	deleteUUIDStr := ctx.Param("uuid")
	if deleteUUIDStr == "" {
		ctx.JSON(400, gin.H{"error": "Invalid uuid"})
		return
	}
	deleteIndex := findProxyIndex(deleteUUIDStr)
	if deleteIndex < 0 || deleteIndex >= len(config.Cfg.Proxy) {
		ctx.JSON(400, gin.H{"error": "Invalid uuid"})
		return
	}
	config.Cfg.Proxy = append(config.Cfg.Proxy[:deleteIndex], config.Cfg.Proxy[deleteIndex+1:]...)
	viper.Set("proxy", config.Cfg.Proxy)
	err := viper.WriteConfig()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{})
}

func findProxyIndex(uuid string) int {
	for i, proxy := range config.Cfg.Proxy {
		if proxy.UUID == uuid {
			return i
		}
	}
	return -1
}

func RestartProxyServer(ctx *gin.Context) {
	config.Init()
	proxyServer.StartProxyServer()
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
	for UUID, proxy := range proxyServer.ProxyManager {
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
		ctx.JSON(400, gin.H{"error": "Invalid uuid"})
		return
	}
	if _, ok := proxyServer.ProxyManager[uuidStr]; !ok {
		ctx.JSON(400, gin.H{"error": "Invalid uuid"})
	}

	proxy := proxyServer.ProxyManager[uuidStr]
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
		ctx.JSON(400, gin.H{"error": "Invalid uuid"})
		return
	}
	if _, ok := proxyServer.ProxyManager[proxyUUIDStr]; !ok {
		ctx.JSON(400, gin.H{"error": "Invalid uuid"})
		return
	}
	linkUUIDStr := ctx.Param("link_uuid")
	if linkUUIDStr == "" {
		ctx.JSON(400, gin.H{"error": "Invalid link uuid"})
		return
	}
	if _, ok := proxyServer.ProxyManager[proxyUUIDStr].Links.Load(linkUUIDStr); !ok {
		ctx.JSON(400, gin.H{"error": "Invalid link uuid"})
		return
	}
	proxy, ok := proxyServer.ProxyManager[proxyUUIDStr].Links.Load(linkUUIDStr)
	if !ok {
		ctx.JSON(500, gin.H{"error": "Internal error"})
		return
	}
	proxy.(*proxyServer.Link).Close()
	ctx.JSON(200, gin.H{})
}
