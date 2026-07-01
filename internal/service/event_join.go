package service

import (
	"context"
	"fmt"
	"time"

	"github.com/GokujyouKaisennDonnburi/NatuEve_API/internal/model"
	"github.com/GokujyouKaisennDonnburi/NatuEve_API/internal/repository"
)


// EventJoinService はイベント参加申込のビジネスロジックを担当する。
type EventJoinService struct {
	repo repository.EventJoinRepository
}


// NewEventJoinService は Service を生成する。
func NewEventJoinService(
	repo repository.EventJoinRepository,
) *EventJoinService {

	return &EventJoinService{
		repo: repo,
	}
}


// Join はイベント参加処理を行う。
func (s *EventJoinService) Join(

	ctx context.Context,

	req model.JoinEventRequest,

	profile_id uuid.UUID,
	username uuid.UUID,
	mail_address uuid.UUID,

) (model.JoinEventResponse, error) {

	// イベント存在確認
	exists, err := s.repo.ExistsEvent(
		ctx,
		req.event_id,
	)

	if err != nil {
		return model.JoinEventResponse{},
			fmt.Errorf("exists event: %w", err)
	}

	if !exists {
		return model.JoinEventResponse{},
			&ValidationError{
				Message: "イベントが存在しません",
			}
	}

	// 重複参加確認
	joined, err := s.repo.ExistsMember(
		ctx,
		req.event_id,
		req.profile_id,
	)

	if err != nil {
		return model.JoinEventResponse{},
			fmt.Errorf("exists member: %w", err)
	}

	if joined {
		return model.JoinEventResponse{},
			&ValidationError{
				Message: "既に参加しています",
			}
	}

	// 定員取得
	capacity, err := s.repo.GetCapacity(
		ctx,
		req.event_id,
	)

	if err != nil {
		return model.JoinEventResponse{},
			fmt.Errorf("get capacity: %w", err)
	}

	// 現在人数取得
	memberCount, err := s.repo.CountMembers(
		ctx,
		req.event_id,
	)

	if err != nil {
		return model.JoinEventResponse{},
			fmt.Errorf("count members: %w", err)
	}

	// 定員チェック
	// capacity == 0 は「定員なし」
	if capacity != 0 && memberCount >= capacity {

		return model.JoinEventResponse{},
			&ValidationError{
				Message: "定員に達しています",
			}
	}

	// Repositoryへ保存
	member := &model.EventMember{

		EventID: req.event_id,

		ProfileID: profile_id,

		Username: username,

		MailAddress: mail_address,

		CreatedAt: time.Now(),
	}

	err = s.repo.Join(
		ctx,
		member,
	)

	if err != nil {

		return model.JoinEventResponse{},
			fmt.Errorf("join event: %w", err)
	}

	// レスポンス作成joinedAt := time.Now()

member.CreatedAt = joinedAt

return model.JoinEventResponse{
    Message:  "イベントへ参加しました",
    JoinedAt: joinedAt,
}, nil
}