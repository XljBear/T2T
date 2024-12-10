package panelServer

import (
	"T2T/config"
	"T2T/proxyServer"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

//go:embed dist
var FrontendDir embed.FS

func StartPanelServer(panelListenAddress string) {

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	r := gin.Default()
	r.Use(cors.Default())
	frontendFs, _ := fs.Sub(FrontendDir, "dist")
	r.StaticFS("/panel", http.FS(frontendFs))
	r.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(302, "/panel")
	})
	r.GET("/api/proxy", func(ctx *gin.Context) {
		ctx.JSON(200, config.Cfg.Proxy)
	})
	r.POST("/api/proxy", func(ctx *gin.Context) {
		proxyData := config.ProxyAddressrRecord{}
		err := ctx.BindJSON(&proxyData)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if proxyData.LocalAddress == "" || proxyData.RemoteAddress == "" || proxyData.Name == "" {
			ctx.JSON(400, gin.H{"error": "LocalAddress, RemoteAddress and Name are required"})
			return
		}
		config.Cfg.Proxy = append(config.Cfg.Proxy, proxyData)
		viper.Set("proxy", config.Cfg.Proxy)
		err = viper.WriteConfig()
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{})
	})
	r.PUT("/api/proxy/:index", func(ctx *gin.Context) {
		updateIndexStr := ctx.Param("index")
		updateIndex, err := strconv.Atoi(updateIndexStr)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid index"})
			return
		}
		updateData := config.ProxyAddressrRecord{}
		err = ctx.BindJSON(&updateData)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if updateData.LocalAddress == "" || updateData.RemoteAddress == "" || updateData.Name == "" {
			ctx.JSON(400, gin.H{"error": "LocalAddress, RemoteAddress and Name are required"})
			return
		}
		if updateIndex < 0 || updateIndex >= len(config.Cfg.Proxy) {
			ctx.JSON(400, gin.H{"error": "Index out of range"})
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
	})
	r.DELETE("/api/proxy/:index", func(ctx *gin.Context) {
		deleteIndexStr := ctx.Param("index")
		deleteIndex, err := strconv.Atoi(deleteIndexStr)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid index"})
			return
		}
		if deleteIndex < 0 || deleteIndex >= len(config.Cfg.Proxy) {
			ctx.JSON(400, gin.H{"error": "Index out of range"})
			return
		}
		config.Cfg.Proxy = append(config.Cfg.Proxy[:deleteIndex], config.Cfg.Proxy[deleteIndex+1:]...)
		viper.Set("proxy", config.Cfg.Proxy)
		err = viper.WriteConfig()
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{})
	})
	r.POST("/api/restart", func(ctx *gin.Context) {
		proxyServer.StartProxyServer()
		ctx.JSON(200, gin.H{})
	})
	fmt.Println("Panel server is running on " + panelListenAddress)
	r.Run(panelListenAddress)
}
