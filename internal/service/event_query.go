package service

import (
	"context"

	"github.com/GokujyouKaisennDonnburi/NatuIve_API/internal/model"
	"github.com/GokujyouKaisennDonnburi/NatuIve_API/internal/repository"
)

const (
	// defaultLimit はページネーションのデフォルト取得件数。
	defaultLimit = 20
	// maxLimit はページネーションで許容する最大取得件数。
	maxLimit = 100
)

// EventQueryService はイベント参照系のビジネスロジックを提供する。
//
// CQRS の Query 側として位置づけ、書き込み系とは分離する。
type EventQueryService struct {
	repo repository.EventRepository
}

// NewEventQueryService は EventQueryService を生成する。
func NewEventQueryService(repo repository.EventRepository) *EventQueryService {
	return &EventQueryService{repo: repo}
}

// List は limit / offset を正規化してからイベントサマリー一覧を返す。
//
// 正規化ルール:
//   - limit が 0 以下 → defaultLimit(20)
//   - limit が maxLimit(100) 超過 → maxLimit(100)
//   - offset が負値 → 0
func (s *EventQueryService) List(ctx context.Context, limit, offset int) ([]model.EventSummary, error) {
	limit = normalizeLimit(limit)
	offset = normalizeOffset(offset)
	return s.repo.ListSummaries(ctx, limit, offset)
}

// normalizeLimit は limit を有効範囲(1〜maxLimit)に丸める。
func normalizeLimit(limit int) int {
	if limit <= 0 {
		return defaultLimit
	}
	if limit > maxLimit {
		return maxLimit
	}
	return limit
}

// normalizeOffset は offset の負値を 0 に丸める。
func normalizeOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}
