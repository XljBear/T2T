package handles

import (
	"T2T/config"
	"T2T/storages"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wenlng/go-captcha/v2/click"
	"time"
)

func Login(ctx *gin.Context) {
	type captchaData struct {
		X int `json:"x"`
		Y int `json:"y"`
	}
	type loginRequestData struct {
		Password    string        `json:"password"`
		CaptchaID   string        `json:"captcha_id"`
		CaptchaData []captchaData `json:"captcha_data"`
	}
	var reqData loginRequestData
	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}
	if reqData.CaptchaID == "" || len(reqData.CaptchaData) == 0 {
		ctx.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}
	captchaID := "c_" + reqData.CaptchaID
	dotData, exist := storages.StorageInstance.Get(captchaID)
	if !exist {
		ctx.JSON(400, gin.H{"error": "Invalid captchaID"})
		return
	}
	defer storages.StorageInstance.Delete(captchaID)
	captchaAnswer := dotData.(map[int]*click.Dot)
	if len(captchaAnswer) != len(reqData.CaptchaData) {
		ctx.JSON(400, gin.H{"error": "Captcha verification failed"})
		return
	}
	for idx, dot := range captchaAnswer {
		if !click.CheckPoint(int64(reqData.CaptchaData[idx].X), int64(reqData.CaptchaData[idx].Y), int64(dot.X), int64(dot.Y), int64(dot.Height), int64(dot.Width), 0) {
			ctx.JSON(400, gin.H{"error": "Captcha verification failed"})
			return
		}
	}
	if reqData.Password != config.Cfg.PanelPassword {
		ctx.JSON(400, gin.H{"error": "Invalid password"})
		return
	}

	loginSession := uuid.New().String()
	loginSessionKey := "l_" + loginSession
	loginSessionExpire := time.Now().Add(time.Hour * 24 * 7)
	storages.StorageInstance.Set(loginSessionKey, "", &loginSessionExpire)
	ctx.JSON(200, gin.H{"token": loginSession})
}

func Logout(ctx *gin.Context) {
	tokenCookie, _ := ctx.Request.Cookie("token")
	if tokenCookie.Value != "" {
		loginSessionKey := "l_" + tokenCookie.Value
		storages.StorageInstance.Delete(loginSessionKey)
	}
	ctx.JSON(200, gin.H{})
}
