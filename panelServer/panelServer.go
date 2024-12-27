package panelServer

import (
	"T2T/proxyServer/routers"
	"embed"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io"
	"io/fs"
	"net/http"
)

//go:embed dist
var FrontendDir embed.FS

func StartPanelServer(panelListenAddress string) {

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	r := gin.Default()
	r.Use(cors.Default())
	frontendFs, _ := fs.Sub(FrontendDir, "dist")
	r.StaticFS("/panel", http.FS(frontendFs))
	r.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(302, "/panel")
	})
	routers.RegisterApiRouter(r)
	fmt.Println("Panel server is running on " + panelListenAddress)
	r.Run(panelListenAddress)
}
