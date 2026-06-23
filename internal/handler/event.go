package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/GokujyouKaisennDonnburi/NatuIve_API/internal/model"
	"github.com/GokujyouKaisennDonnburi/NatuIve_API/internal/service"
)

// EventHandler はイベント系のエンドポイントを担当する。
type EventHandler struct {
	svc *service.EventQueryService
}

// NewEventHandler は EventHandler を生成する。
func NewEventHandler(svc *service.EventQueryService) *EventHandler {
	return &EventHandler{svc: svc}
}

// List godoc
//
//	@Summary		イベント一覧取得
//	@Description	公開イベントを作成日時の降順で返す。認証不要。
//	@Tags			event
//	@Produce		json
//	@Param			limit	query		int	false	"取得件数(default 20, 最大 100)"
//	@Param			offset	query		int	false	"取得開始位置(default 0)"
//	@Success		200		{array}		model.EventSummary
//	@Failure		500		{object}	model.ErrorResponse
//	@Router			/api/v1/events [get]
func (h *EventHandler) List(c *gin.Context) {
	// クエリパラメータを取得する（不正値は service 層で安全側に丸める）。
	limit := queryInt(c, "limit", 0)
	offset := queryInt(c, "offset", 0)

	summaries, err := h.svc.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("internal_error", "イベント一覧の取得に失敗しました"))
		return
	}

	c.JSON(http.StatusOK, summaries)
}

// queryInt はクエリパラメータを int に変換する。変換できない場合は defaultVal を返す。
func queryInt(c *gin.Context, key string, defaultVal int) int {
	raw := c.Query(key)
	if raw == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return defaultVal
	}
	return v
}
