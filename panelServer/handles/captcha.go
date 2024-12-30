package handles

import (
	"T2T/panelServer/storages"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wenlng/go-captcha-assets/resources/images"
	"github.com/wenlng/go-captcha-assets/resources/shapes"
	"github.com/wenlng/go-captcha/v2/click"
	"log"
	"time"
)

var textCapt click.Captcha

func Captcha(ctx *gin.Context) {
	type CaptchaData struct {
		CaptchaID string `json:"captcha_id"`
		Captcha   string `json:"captcha"`
		Thumb     string `json:"thumb"`
	}
	builder := click.NewBuilder()
	shapeMaps, err := shapes.GetShapes()
	if err != nil {
		log.Fatalln(err)
	}
	imgs, err := images.GetImages()
	if err != nil {
		log.Fatalln(err)
	}
	builder.SetResources(
		click.WithShapes(shapeMaps),
		click.WithBackgrounds(imgs),
	)
	textCapt = builder.MakeWithShape()
	captData, err := textCapt.Generate()
	if err != nil {
		log.Fatalln(err)
	}
	dotData := captData.GetData()
	if dotData == nil {
		log.Fatalln(">>>>> generate err")
	}

	var mBase64, tBase64 string
	mBase64, err = captData.GetMasterImage().ToBase64()
	if err != nil {
		fmt.Println(err)
	}
	tBase64, err = captData.GetThumbImage().ToBase64()
	if err != nil {
		fmt.Println(err)
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
