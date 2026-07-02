package model

import (
"time"  
"github.com/google/uuid"
)

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
// EventID は参加したイベントのUUID
EventID uuid.UUID `json:"eventId" example:"a1b2c3d4-e5f6-7890-abcd-ef1234567890"`
// ProfileID は参加するユーザーのUUID(あれば)
ProfileID uuid.UUID `json:"profileId" example:"b2c3d4e5-f6a7-8901-bcde-f23456789012"`
// Usenameは参加するユーザーの表示名(必須)
Username string `json:"username" example:"ヒラコウキ"`
// MailAddressは参加するユーザーのメールアドレス(必須)
MailAddress string `json:"mailAddress" example:"hirakouki41@gmail.com"`
// CreatedAt は参加申込日時
CreatedAt time.Time `json:"createdAt" example:"2023-01-01T12:00:00Z"`
}

// EventMember は event_members テーブルと対応するモデル。
// Repository層でINSERT・SELECTする際に使用する。
type NewEvent struct {

EventID     uuid.UUID
ProfileID   uuid.UUID
Username    string
MailAddress string
}
