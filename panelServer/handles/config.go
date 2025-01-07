package handles

import (
	"T2T/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetConfig(ctx *gin.Context) {
	type responseData struct {
		CaptchaType uint `json:"captcha_type"`
		DarkMode    bool `json:"dark_mode"`
	}
	data := responseData{
		CaptchaType: config.Cfg.CaptchaType,
		DarkMode:    config.Cfg.DarkMode,
	}
	ctx.JSON(http.StatusOK, data)
	return
}

func ReloadConfig(ctx *gin.Context) {
	config.ReloadConfig()
	ctx.JSON(http.StatusOK, config.Cfg)
}
