package handles

import (
	"T2T/config"
	"T2T/storages"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wenlng/go-captcha/v2/click"
	"github.com/wenlng/go-captcha/v2/rotate"
	"github.com/wenlng/go-captcha/v2/slide"
	"net/http"
	"time"
)

func Login(ctx *gin.Context) {
	type XYCaptchaData struct {
		X int `json:"x"`
		Y int `json:"y"`
	}
	type loginRequestData struct {
		Password          string          `json:"password"`
		CaptchaID         string          `json:"captcha_id"`
		ClickCaptchaData  []XYCaptchaData `json:"click_captcha_data"`
		SlideCaptchaData  XYCaptchaData   `json:"slide_captcha_data"`
		RotateCaptchaData int             `json:"rotate_captcha_data"`
	}
	var reqData loginRequestData
	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	if config.Cfg.CaptchaType != 0 {
		if reqData.CaptchaID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			return
		}
		switch config.Cfg.CaptchaType {
		case 1, 2:
			if len(reqData.ClickCaptchaData) == 0 {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
				return
			}
			captchaID := "c_" + reqData.CaptchaID
			dotData, exist := storages.StorageInstance.Get(captchaID)
			if !exist {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid captchaID"})
				return
			}
			defer storages.StorageInstance.Delete(captchaID)
			captchaAnswer := dotData.(map[int]*click.Dot)
			if len(captchaAnswer) != len(reqData.ClickCaptchaData) {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Captcha verification failed"})
				return
			}
			for idx, dot := range captchaAnswer {
				if !click.CheckPoint(int64(reqData.ClickCaptchaData[idx].X), int64(reqData.ClickCaptchaData[idx].Y), int64(dot.X), int64(dot.Y), int64(dot.Height), int64(dot.Width), 0) {
					ctx.JSON(http.StatusBadRequest, gin.H{"error": "Captcha verification failed"})
					return
				}
			}
		case 3, 4:
			captchaID := "c_" + reqData.CaptchaID
			slideData, exist := storages.StorageInstance.Get(captchaID)
			if !exist {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid captchaID"})
				return
			}
			defer storages.StorageInstance.Delete(captchaID)
			captchaAnswer := slideData.(*slide.Block)
			if !slide.CheckPoint(int64(reqData.SlideCaptchaData.X), int64(reqData.SlideCaptchaData.Y), int64(captchaAnswer.X), int64(captchaAnswer.Y), 4) {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Captcha verification failed"})
				return
			}
		case 5:
			captchaID := "c_" + reqData.CaptchaID
			rotateData, exist := storages.StorageInstance.Get(captchaID)
			if !exist {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid captchaID"})
				return
			}
			defer storages.StorageInstance.Delete(captchaID)
			captchaAnswer := rotateData.(*rotate.Block)
			if !rotate.CheckAngle(int64(reqData.RotateCaptchaData), int64(captchaAnswer.Angle), 2) {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Captcha verification failed"})
				return
			}
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Captcha verification failed"})
			return
		}
	}
	if reqData.Password != config.Cfg.PanelPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}

	loginSession := uuid.New().String()
	loginSessionKey := "l_" + loginSession
	loginSessionExpire := time.Now().Add(time.Hour * 24 * 7)
	storages.StorageInstance.Set(loginSessionKey, "", &loginSessionExpire)
	ctx.JSON(http.StatusOK, gin.H{"token": loginSession})
}

func Logout(ctx *gin.Context) {
	tokenCookie, _ := ctx.Request.Cookie("token")
	if tokenCookie.Value != "" {
		loginSessionKey := "l_" + tokenCookie.Value
		storages.StorageInstance.Delete(loginSessionKey)
	}
	ctx.JSON(http.StatusOK, gin.H{})
}
