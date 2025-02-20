package handles

import (
	"T2T/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetInfo(ctx *gin.Context) {
	type responseData struct {
		CaptchaType uint   `json:"captcha_type"`
		DarkMode    bool   `json:"dark_mode"`
		CommitID    string `json:"commit_id"`
		Version     string `json:"version"`
		BuildTime   string `json:"build_time"`
		OS          string `json:"os"`
		Arch        string `json:"arch"`
	}
	data := responseData{
		CaptchaType: config.Cfg.CaptchaType,
		DarkMode:    config.Cfg.DarkMode,
		CommitID:    config.CommitID,
		Version:     config.Version,
		BuildTime:   config.BuildTime,
		OS:          config.OS,
		Arch:        config.Arch,
	}
	ctx.JSON(http.StatusOK, data)
	return
}

func ReloadConfig(ctx *gin.Context) {
	_ = config.ReloadConfig()
	ctx.JSON(http.StatusOK, config.Cfg)
}
