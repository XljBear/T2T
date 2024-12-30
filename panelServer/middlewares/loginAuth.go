package middlewares

import (
	"T2T/panelServer/storages"
	"github.com/gin-gonic/gin"
)

func LoginAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenCookie, err := ctx.Request.Cookie("token")
		if err != nil {
			ctx.JSON(401, gin.H{
				"error": "Unauthorized",
			})
			ctx.Abort()
			return
		}
		if tokenCookie.Value == "" {
			ctx.JSON(401, gin.H{
				"error": "Unauthorized",
			})
			ctx.Abort()
			return
		}
		loginSessionKey := "l_" + tokenCookie.Value
		_, exist := storages.StorageInstance.Get(loginSessionKey)
		if !exist {
			ctx.JSON(401, gin.H{
				"error": "Unauthorized",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
