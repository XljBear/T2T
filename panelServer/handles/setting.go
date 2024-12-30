package handles

import (
	"T2T/config"
	"T2T/storages"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func UpdateSetting(ctx *gin.Context) {
	secretPassword := "~nononono$y0ucantsee.meme@"
	type SettingRequestData struct {
		PanelPassword string `json:"panel_password"`
	}
	settingRequest := &SettingRequestData{}
	err := ctx.ShouldBindJSON(settingRequest)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Invalid request data",
		})
		return
	}
	if settingRequest.PanelPassword == secretPassword {
		ctx.JSON(200, gin.H{})
		return
	}
	config.Cfg.PanelPassword = settingRequest.PanelPassword
	viper.Set("panel_password", config.Cfg.PanelPassword)
	err = viper.WriteConfig()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	storages.StorageInstance.DeleteWithPrefix("l_")
	ctx.JSON(200, gin.H{})
}
