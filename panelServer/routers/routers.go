package routers

import (
	"T2T/panelServer/handles"
	"T2T/panelServer/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterApiRouter(app *gin.Engine) {
	// 登录相关接口
	app.GET("/api/captcha", handles.Captcha)
	app.GET("/api/info", handles.GetInfo)
	app.POST("/api/login", handles.Login)
	app.POST("/api/logout", handles.Logout)

	// 登陆鉴权中间件
	app.Use(middlewares.LoginAuth())
	// 代理服务相关接口
	app.GET("/api/proxy", handles.GetProxyList)
	app.POST("/api/proxy", handles.CreateProxy)
	app.PUT("/api/proxy/:uuid", handles.UpdateProxy)
	app.DELETE("/api/proxy/:uuid", handles.DeleteProxy)
	app.GET("/api/traffic", handles.GetTraffic)
	app.GET("/api/proxy/:uuid/links", handles.GetLinks)
	app.DELETE("/api/proxy/:uuid/links/:link_uuid", handles.KickProxyServer)
	// 系统相关接口
	app.POST("/api/restart", handles.RestartProxyServer)
	// 系统设置相关接口
	app.POST("/api/setting", handles.UpdateSetting)
	app.POST("/api/reload", handles.ReloadConfig)
	// 防火墙相关接口
	app.GET("/api/ipRules", handles.GetIPRules)
	app.POST("/api/ipRules/reload", handles.ReloadAllowBlock)
	app.PUT("/api/ipRules/mode", handles.UpdateRunMode)
	app.DELETE("/api/ipRules/allow/:uuid", handles.DeleteAllowIPRule)
	app.DELETE("/api/ipRules/block/:uuid", handles.DeleteBlockIPRule)
	app.POST("/api/ipRules/allow", handles.CreateAllowIPRule)
	app.POST("/api/ipRules/block", handles.CreateBlockIPRule)
}
