package handles

import (
	"T2T/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetConfig(ctx *gin.Context) {
	type responseData struct {
		CaptchaType uint `json:"captcha_type"`
	}
	data := responseData{
		CaptchaType: config.Cfg.CaptchaType,
	}
	ctx.JSON(http.StatusOK, data)
	return
}
