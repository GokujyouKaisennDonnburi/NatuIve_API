package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/GokujyouKaisennDonnburi/NatuEve_API/internal/middleware"
	"github.com/GokujyouKaisennDonnburi/NatuEve_API/internal/model"
	"github.com/GokujyouKaisennDonnburi/NatuEve_API/internal/service"
)

// Join はイベント参加申込API。
//
//	@Summary		イベント参加
//	@Description	ログイン中のユーザーがイベントへ参加する。
//	@Tags			event
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		model.JoinEventRequest	true	"参加申込"
//	@Success		200		{object}	model.JoinEventResponse
//	@Failure		400		{object}	model.ValidationErrorResponse
//	@Failure		401		{object}	model.UnauthorizedErrorResponse
//	@Failure		500		{object}	model.InternalErrorResponse
//	@Router			/api/v1/events/join [post]
func (h *EventHandler) Join(c *gin.Context) {

	// JWT認証情報取得
	authUser, ok := middleware.AuthFromContext(c)

	if !ok {

		c.JSON(
			http.StatusUnauthorized,
			model.NewErrorResponse(
				"unauthorized",
				"認証が必要です",
			),
		)

		return
	}

	// JSON受け取り
	var req model.JoinEventRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(
			http.StatusBadRequest,
			model.NewErrorResponse(
				"invalid_request",
				"リクエストが不正です",
			),
		)

		return
	}

	// Service呼び出し
	resp, err := h.joinSvc.Join(

		c.Request.Context(),

		req,

		authUser.id,

		authUser.display_name,

		authUser.mail_address,
	)

	if err != nil {

		var ve *service.ValidationError

		if errors.As(err, &ve) {

			c.JSON(
				http.StatusBadRequest,
				model.NewErrorResponse(
					"invalid_request",
					ve.Message,
				),
			)

			return
		}

		c.JSON(
			http.StatusInternalServerError,
			model.NewErrorResponse(
				"internal_error",
				"参加申込に失敗しました",
			),
		)

		return
	}

	// 成功
	c.JSON(
		http.StatusOK,
		resp,
	)
}