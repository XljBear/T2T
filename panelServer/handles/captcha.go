package handles

import (
	"T2T/storages"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wenlng/go-captcha-assets/resources/images"
	"github.com/wenlng/go-captcha-assets/resources/shapes"
	"github.com/wenlng/go-captcha/v2/click"
	"time"
)

var shapeCapt click.Captcha

func Captcha(ctx *gin.Context) {
	type CaptchaData struct {
		CaptchaID string `json:"captcha_id"`
		Captcha   string `json:"captcha"`
		Thumb     string `json:"thumb"`
	}
	builder := click.NewBuilder()
	shapeMaps, err := shapes.GetShapes()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	imgs, err := images.GetImages()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	builder.SetResources(
		click.WithShapes(shapeMaps),
		click.WithBackgrounds(imgs),
	)
	shapeCapt = builder.MakeWithShape()
	captData, err := shapeCapt.Generate()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	dotData := captData.GetData()
	if dotData == nil {
		ctx.JSON(500, gin.H{"error": "Captcha generate failed"})
		return
	}

	var mBase64, tBase64 string
	mBase64, err = captData.GetMasterImage().ToBase64()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	tBase64, err = captData.GetThumbImage().ToBase64()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	capt := CaptchaData{
		CaptchaID: uuid.New().String(),
		Captcha:   mBase64,
		Thumb:     tBase64,
	}
	expiredTime := time.Now().Add(time.Minute * 3)
	storages.StorageInstance.Set("c_"+capt.CaptchaID, dotData, &expiredTime)
	ctx.JSON(200, capt)
}
