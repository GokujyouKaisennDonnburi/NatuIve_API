package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler はヘルスチェック系のエンドポイントを担当する。
type HealthHandler struct{}

// NewHealthHandler は HealthHandler を生成する。
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check godoc
//
//	@Summary		ヘルスチェック
//	@Description	サーバーが正常に稼働しているか確認する
//	@Tags			system
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
