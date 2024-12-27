package routers

import (
	"T2T/proxyServer/handles"
	"github.com/gin-gonic/gin"
)

func RegisterApiRouter(app *gin.Engine) {
	app.GET("/api/proxy", handles.GetProxyList)
	app.POST("/api/proxy", handles.CreateProxy)
	app.PUT("/api/proxy/:uuid", handles.UpdateProxy)
	app.DELETE("/api/proxy/:uuid", handles.DeleteProxy)
	app.POST("/api/restart", handles.RestartProxyServer)
	app.GET("/api/traffic", handles.GetTraffic)
	app.GET("/api/proxy/:uuid/links", handles.GetLinks)
	app.DELETE("/api/proxy/:uuid/links/:link_uuid", handles.KickProxyServer)
}
