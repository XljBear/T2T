package middlewares

import (
	"T2T/storages"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoginAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenCookie, err := ctx.Request.Cookie("token")
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			ctx.Abort()
			return
		}
		if tokenCookie.Value == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			ctx.Abort()
			return
		}
		loginSessionKey := "l_" + tokenCookie.Value
		if !storages.StorageInstance.Exists(loginSessionKey) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
