package handles

import (
	"T2T/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func GetIPRules(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, &config.AllowBlockCfg.AllowBlock)
}
func ReloadAllowBlock(ctx *gin.Context) {
	config.ReloadAllowBlock()
	ctx.JSON(http.StatusOK, config.Cfg)
}
func UpdateRunMode(ctx *gin.Context) {
	type requestData struct {
		Mode int `json:"mode"`
	}
	var req requestData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	if req.Mode < 0 || req.Mode > 2 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mode"})
		return
	}
	config.AllowBlockCfg.AllowBlock.Mode = req.Mode
	config.SaveAllowBlock()
	ctx.JSON(http.StatusOK, gin.H{})
}

func DeleteAllowIPRule(ctx *gin.Context) {
	deleteUUIDStr := ctx.Param("uuid")
	if deleteUUIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid uuid"})
		return
	}
	config.AllowBlockCfg.AllowBlock.DeleteAllowIPByUUID(deleteUUIDStr)
	ctx.JSON(http.StatusOK, gin.H{})
}

func DeleteBlockIPRule(ctx *gin.Context) {
	deleteUUIDStr := ctx.Param("uuid")
	if deleteUUIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid uuid"})
		return
	}
	config.AllowBlockCfg.AllowBlock.DeleteBlockIPByUUID(deleteUUIDStr)
	ctx.JSON(http.StatusOK, gin.H{})
}

func CreateIPRule(ctx *gin.Context) {
	type requestData struct {
		IP      string     `json:"ip"`
		Port    []string   `json:"port"`
		EndTime *time.Time `json:"end_time"`
		Reason  string     `json:"reason"`
		Type    int        `json:"type"`
	}
	var req requestData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	if req.IP == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "IP are required"})
		return
	}
	ipItem := config.IPItem{
		UUID:      uuid.New().String(),
		IP:        req.IP,
		Port:      req.Port,
		StartTime: time.Now(),
		EndTime:   req.EndTime,
		Reason:    req.Reason,
	}
	switch req.Type {
	case 0:
		config.AllowBlockCfg.AllowBlock.AddBlockIP(ipItem)
	case 1:
		config.AllowBlockCfg.AllowBlock.AddAllowIP(ipItem)
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}
