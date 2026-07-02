package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/GokujyouKaisennDonnburi/NatuEve_API/internal/middleware"
	"github.com/GokujyouKaisennDonnburi/NatuEve_API/internal/model"
	"github.com/GokujyouKaisennDonnburi/NatuEve_API/internal/service"
)

// Join はイベント参加申込 API。
//
//	@Summary		イベント参加
//	@Description	ログイン中のユーザーがイベントへ参加する。
//	@Tags			event
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string					true	"イベントID"
//	@Param			body	body		model.JoinEventRequest	true	"参加申込"
//	@Success		201		{object}	model.JoinEventResponse
//	@Failure		400		{object}	model.ValidationErrorResponse
//	@Failure		401		{object}	model.UnauthorizedErrorResponse
//	@Failure		404		{object}	model.NotFoundErrorResponse
//	@Failure		409		{object}	model.ConflictErrorResponse
//	@Failure		500		{object}	model.InternalErrorResponse
//	@Router			/api/v1/events/{id}/join [post]
func (h *EventHandler) Join(c *gin.Context) {

	// JWT認証情報取得
	authUser, ok := middleware.AuthFromContext(c)
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			model.NewErrorResponse("unauthorized", "認証が必要です"),
		)
		return
	}

	// パスパラメータからイベントID取得
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.NewErrorResponse("invalid_request", "イベントIDが不正です"),
		)
		return
	}

	// JWTからプロフィールID取得
	profileID, err := uuid.Parse(authUser.ID)
	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			model.NewErrorResponse("unauthorized", "ユーザーIDが不正です"),
		)
		return
	}

	// JSON受け取り
	var req model.JoinEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.NewErrorResponse("invalid_request", "リクエストが不正です"),
		)
		return
	}

	// Service呼び出し
	resp, err := h.joinSvc.Join(c.Request.Context(), eventID, profileID, req)
	if err != nil {
		var ve *service.ValidationError
		if errors.As(err, &ve) {
			c.JSON(
				http.StatusBadRequest,
				model.NewErrorResponse("invalid_request", ve.Message),
			)
			return
		}

		var nfe *service.NotFoundError
		if errors.As(err, &nfe) {
			c.JSON(
				http.StatusNotFound,
				model.NewErrorResponse("not_found", nfe.Message),
			)
			return
		}

		var ce *service.ConflictError
		if errors.As(err, &ce) {
			c.JSON(
				http.StatusConflict,
				model.NewErrorResponse("conflict", ce.Message),
			)
			return
		}

		c.JSON(
			http.StatusInternalServerError,
			model.NewErrorResponse("internal_error", "参加申込に失敗しました"),
		)
		return
	}

	// 成功
	c.JSON(http.StatusCreated, resp)
}
