package handles

import (
	"T2T/config"
	"T2T/storages"
	"github.com/gin-gonic/gin"
	"github.com/golang/freetype/truetype"
	"github.com/google/uuid"
	"github.com/wenlng/go-captcha-assets/bindata/chars"
	"github.com/wenlng/go-captcha-assets/resources/fonts/fzshengsksjw"
	"github.com/wenlng/go-captcha-assets/resources/images"
	"github.com/wenlng/go-captcha-assets/resources/shapes"
	"github.com/wenlng/go-captcha-assets/resources/tiles"
	"github.com/wenlng/go-captcha/v2/base/option"
	"github.com/wenlng/go-captcha/v2/click"
	"github.com/wenlng/go-captcha/v2/rotate"
	"github.com/wenlng/go-captcha/v2/slide"
	"log"
	"net/http"
	"time"
)

func Captcha(ctx *gin.Context) {
	type CaptchaData struct {
		CaptchaID   string `json:"captcha_id"`
		Captcha     string `json:"captcha"`
		Thumb       string `json:"thumb"`
		ThumbX      int    `json:"thumb_x"`
		ThumbY      int    `json:"thumb_y"`
		ThumbWidth  int    `json:"thumb_width"`
		ThumbHeight int    `json:"thumb_height"`
		Angle       int    `json:"angle"`
	}
	capt := CaptchaData{}
	switch config.Cfg.CaptchaType {
	case 1:
		mBase64, tBase64, dotData, err := GenerateTextCaptcha()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		capt = CaptchaData{
			CaptchaID: uuid.New().String(),
			Captcha:   mBase64,
			Thumb:     tBase64,
		}
		expiredTime := time.Now().Add(time.Minute * 3)
		storages.StorageInstance.Set("c_"+capt.CaptchaID, dotData, &expiredTime)
	case 2:
		mBase64, tBase64, dotData, err := GenerateShapeCaptcha()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		capt = CaptchaData{
			CaptchaID: uuid.New().String(),
			Captcha:   mBase64,
			Thumb:     tBase64,
		}
		expiredTime := time.Now().Add(time.Minute * 3)
		storages.StorageInstance.Set("c_"+capt.CaptchaID, dotData, &expiredTime)
	case 3:
		mBase64, tBase64, blockData, err := GenerateSlideCaptcha()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		capt = CaptchaData{
			CaptchaID:   uuid.New().String(),
			Captcha:     mBase64,
			Thumb:       tBase64,
			ThumbHeight: blockData.Height,
			ThumbWidth:  blockData.Width,
			ThumbX:      0,
			ThumbY:      blockData.Y,
		}
		expiredTime := time.Now().Add(time.Minute * 3)
		storages.StorageInstance.Set("c_"+capt.CaptchaID, blockData, &expiredTime)
	case 4:
		mBase64, tBase64, blockData, err := GenerateSlideCaptcha()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		capt = CaptchaData{
			CaptchaID:   uuid.New().String(),
			Captcha:     mBase64,
			Thumb:       tBase64,
			ThumbHeight: blockData.Height,
			ThumbWidth:  blockData.Width,
			ThumbX:      0,
			ThumbY:      0,
		}
		expiredTime := time.Now().Add(time.Minute * 3)
		storages.StorageInstance.Set("c_"+capt.CaptchaID, blockData, &expiredTime)
	case 5:
		mBase64, tBase64, blockData, err := GenerateRotateCaptcha()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		capt = CaptchaData{
			CaptchaID: uuid.New().String(),
			Captcha:   mBase64,
			Thumb:     tBase64,
			Angle:     blockData.Angle,
		}
		expiredTime := time.Now().Add(time.Minute * 3)
		storages.StorageInstance.Set("c_"+capt.CaptchaID, blockData, &expiredTime)
	}

	ctx.JSON(http.StatusOK, capt)
}

func GenerateShapeCaptcha() (mBase64 string, tBase64 string, dotData map[int]*click.Dot, err error) {
	var shapeCapt click.Captcha
	builder := click.NewBuilder()
	shapeMaps, err := shapes.GetShapes()
	if err != nil {
		return
	}
	imgs, err := images.GetImages()
	if err != nil {
		return
	}
	builder.SetResources(
		click.WithShapes(shapeMaps),
		click.WithBackgrounds(imgs),
	)
	shapeCapt = builder.MakeWithShape()
	captData, err := shapeCapt.Generate()
	if err != nil {
		return
	}
	dotData = captData.GetData()
	if dotData == nil {
		return
	}
	mBase64, err = captData.GetMasterImage().ToBase64()
	if err != nil {
		return
	}
	tBase64, err = captData.GetThumbImage().ToBase64()
	return
}

func GenerateTextCaptcha() (mBase64 string, tBase64 string, dotData map[int]*click.Dot, err error) {
	var textCapt click.Captcha
	builder := click.NewBuilder()
	fonts, err := fzshengsksjw.GetFont()
	if err != nil {
		log.Fatalln(err)
	}

	imgs, err := images.GetImages()
	if err != nil {
		return
	}
	builder.SetResources(
		click.WithChars(chars.GetChineseChars()),
		click.WithFonts([]*truetype.Font{fonts}),
		click.WithBackgrounds(imgs),
	)
	textCapt = builder.Make()
	captData, err := textCapt.Generate()
	if err != nil {
		return
	}
	dotData = captData.GetData()
	if dotData == nil {
		return
	}
	mBase64, err = captData.GetMasterImage().ToBase64()
	if err != nil {
		return
	}
	tBase64, err = captData.GetThumbImage().ToBase64()
	return
}

func GenerateSlideCaptcha() (mBase64 string, tBase64 string, blockData *slide.Block, err error) {
	builder := slide.NewBuilder(
		slide.WithGenGraphNumber(2),
		slide.WithEnableGraphVerticalRandom(true),
	)

	imgs, err := images.GetImages()
	if err != nil {
		return
	}
	graphs, err := tiles.GetTiles()
	if err != nil {
		return
	}

	var newGraphs = make([]*slide.GraphImage, 0, len(graphs))
	for i := 0; i < len(graphs); i++ {
		graph := graphs[i]
		newGraphs = append(newGraphs, &slide.GraphImage{
			OverlayImage: graph.OverlayImage,
			MaskImage:    graph.MaskImage,
			ShadowImage:  graph.ShadowImage,
		})
	}

	builder.SetResources(
		slide.WithGraphImages(newGraphs),
		slide.WithBackgrounds(imgs),
	)

	slideCapt := builder.Make()
	captData, err := slideCapt.Generate()
	if err != nil {
		return
	}
	mBase64, err = captData.GetMasterImage().ToBase64()
	if err != nil {
		return
	}
	tBase64, err = captData.GetTileImage().ToBase64()
	if err != nil {
		return
	}
	blockData = captData.GetData()
	return
}

func GenerateRotateCaptcha() (mBase64 string, tBase64 string, blockData *rotate.Block, err error) {
	builder := rotate.NewBuilder(rotate.WithRangeAnglePos([]option.RangeVal{
		{Min: 20, Max: 330},
	}))

	imgs, err := images.GetImages()
	if err != nil {
		return
	}

	builder.SetResources(
		rotate.WithImages(imgs),
	)

	rotateCapt := builder.Make()
	captData, err := rotateCapt.Generate()
	if err != nil {
		return
	}

	blockData = captData.GetData()
	if blockData == nil {
		return
	}
	mBase64, err = captData.GetMasterImage().ToBase64()
	if err != nil {
		return
	}
	tBase64, err = captData.GetThumbImage().ToBase64()
	if err != nil {
		return
	}
	return
}
