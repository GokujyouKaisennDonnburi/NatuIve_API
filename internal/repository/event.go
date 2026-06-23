package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/GokujyouKaisennDonnburi/NatuIve_API/internal/model"
)

// EventRepository は events テーブルへのアクセスを抽象化する。
type EventRepository interface {
	// ListSummaries は作成日時の降順でイベントサマリーを取得する。
	// 同一作成日時のレコードは uuid 昇順で安定ソートする。
	ListSummaries(ctx context.Context, limit, offset int) ([]model.EventSummary, error)
}

// eventPostgres は EventRepository の PostgreSQL 実装。
type eventPostgres struct {
	db *sql.DB
}

// NewEventRepository は *sql.DB を使う EventRepository を生成する。
func NewEventRepository(db *sql.DB) EventRepository {
	return &eventPostgres{db: db}
}

// ListSummaries は一覧表示に必要なカラムのみ SELECT する。
// description / external_url / capacity / updated_at は取得しない。
func (r *eventPostgres) ListSummaries(ctx context.Context, limit, offset int) ([]model.EventSummary, error) {
	const query = `
		SELECT id, title, event_date, location, profile_id, created_at
		FROM events
		ORDER BY created_at DESC, id
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list event summaries: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var summaries []model.EventSummary
	for rows.Next() {
		var s model.EventSummary
		var (
			location  sql.NullString
			profileID sql.NullString
		)
		if err := rows.Scan(
			&s.ID,
			&s.Title,
			&s.EventDate,
			&location,
			&profileID,
			&s.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan event summary: %w", err)
		}
		s.Location = location.String
		s.ProfileID = profileID.String
		summaries = append(summaries, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate event summaries: %w", err)
	}

	// レコードが 0 件でも nil ではなく空スライスを返す。
	if summaries == nil {
		summaries = []model.EventSummary{}
	}
	return summaries, nil
}
