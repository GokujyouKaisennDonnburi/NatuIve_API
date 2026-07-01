package model

import (
	"time"
	"github.com/google/uuid"
)

// JoinEventRequest はイベント参加申込APIのリクエストDTO。
//
// フロントから送られてくるJSONを受け取るための構造体。
type JoinEventRequest struct {

	// EventID は参加したいイベントのUUID
	event_id uuid.UUID `json:"eventId" validate:"required,uuid"`
}

// JoinEventResponse は参加申込完了時に返すレスポンス。
type JoinEventResponse struct {

	// Message は処理結果
	Message string `json:"message"`

	// JoinedAt は参加日時
	JoinedAt time.Time `json:"joinedAt"`
}

// EventMember は event_members テーブルと対応するモデル。
// Repository層でINSERT・SELECTする際に使用する。
type EventMember struct {

	// ID は event_members のUUID
	id uuid.UUID `db:"id"`

	// event_id は参加イベント
	event_id uuid.UUID `db:"event_id"`

	// profile_id は参加ユーザー
	profile_id uuid.UUID `db:"profile_id"`

	// username は参加時の表示名
	username string `db:"username"`

	// mail_address は参加時のメールアドレス
	mail_address string `db:"mail_address"`
}
