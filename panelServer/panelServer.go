package panelServer

import (
	"T2T/panelServer/routers"
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/fs"
	"log"
	"net/http"
	"time"
)

//go:embed dist
var FrontendDir embed.FS

func StartPanelServer(panelListenAddress string) bool {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	r := gin.Default()
	frontendFs, _ := fs.Sub(FrontendDir, "dist")
	r.StaticFS("/panel", http.FS(frontendFs))
	r.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(302, "/panel")
	})
	routers.RegisterApiRouter(r)
	log.Println("Panel server starting on " + panelListenAddress + "...")
	var broken chan bool
	broken = make(chan bool, 1)

	go func() {
		err := r.Run(panelListenAddress)
		if err != nil {
			fmt.Printf("Panel server start with error: (%v)\n", err)
		}
		broken <- true
	}()
	select {
	case <-broken:
		return false
	case <-time.After(time.Second * 1):
		return true
	}
}
