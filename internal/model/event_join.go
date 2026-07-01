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
	EventID uuid.UUID `json:"eventId" example:"a1b2c3d4-e5f6-7890-abcd-ef1234567890" validate:"required,max=255"`
	// ProfileID は参加するユーザーのUUID(あれば)
	ProfileID uuid.UUID `json:"profileId" example:"b2c3d4e5-f6a7-8901-bcde-f23456789012" validate:"omitempty"`
	// Usenameは参加するユーザーの表示名(必須)
	Username string `json:"username" example:"山田太郎" validate:"required,max=255"`
	// MailAddressは参加するユーザーのメールアドレス(必須)
	MailAddress string `json:"mailAddress" example:"yamada@example.com" validate:"required,email,max=255"`
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
	EventID     string
	ProfileID   string
	Username    string
	MailAddress string
}
