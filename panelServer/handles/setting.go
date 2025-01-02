package handles

import (
	"T2T/config"
	"T2T/storages"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func UpdateSetting(ctx *gin.Context) {
	secretPassword := "~nononono$y0ucantsee.meme@"
	type SettingRequestData struct {
		PanelPassword string `json:"panel_password"`
		CaptchaType   uint   `json:"captcha_type"`
	}
	settingRequest := &SettingRequestData{}
	err := ctx.ShouldBindJSON(settingRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}
	if settingRequest.PanelPassword != secretPassword {
		config.Cfg.PanelPassword = settingRequest.PanelPassword
		viper.Set("panel_password", config.Cfg.PanelPassword)
		storages.StorageInstance.DeleteWithPrefix("l_")
	}
	config.Cfg.CaptchaType = settingRequest.CaptchaType
	viper.Set("captcha_type", config.Cfg.CaptchaType)
	err = viper.WriteConfig()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}
