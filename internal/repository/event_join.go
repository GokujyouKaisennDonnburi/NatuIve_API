package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/GokujyouKaisennDonnburi/NatuEve_API/internal/model"
)

// EventJoinRepository はイベント参加申込用Repositoryのインターフェース。
// Service層はこのInterfaceだけを知っていればよく、
// 実際のDB実装(PostgreSQLなど)には依存しない。
type EventJoinRepository interface {

	// ExistsEvent はイベントが存在するか確認する。
	ExistsEvent(ctx context.Context, event_id uuid.UUID) (bool, error)

	// ExistsMember はユーザーが既に参加済みか確認する。
	ExistsMember(ctx context.Context, event_id, profile_id uuid.UUID) (bool, error)

	// CountMembers は現在の参加人数を取得する。
	CountMembers(ctx context.Context, event_id uuid.UUID) (int, error)

	// GetCapacity はイベントの定員を取得する。
	GetCapacity(ctx context.Context, event_id uuid.UUID) (int, error)

	// Join はイベント参加を登録する。
	Join(ctx context.Context, member *model.EventMember) error
}

// eventJoinPostgres は PostgreSQL実装。
type eventJoinPostgres struct {
	db *sql.DB
}

// NewEventJoinRepository はRepositoryを生成する。
func NewEventJoinRepository(db *sql.DB) EventJoinRepository {
	return &eventJoinPostgres{
		db: db,
	}
}

// イベント存在確認
func (r *eventJoinPostgres) ExistsEvent(
	ctx context.Context,
	event_id uuid.UUID,
) (bool, error) {

	const query = `
	SELECT EXISTS(
		SELECT 1
		FROM events
		WHERE id = $1
	)
	`

	var exists bool

	err := r.db.QueryRowContext(
		ctx,
		query,
		event_id,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("exists event: %w", err)
	}

	return exists, nil
}

// 参加済み確認
func (r *eventJoinPostgres) ExistsMember(
	ctx context.Context,
	event_id uuid.UUID,
	profile_id uuid.UUID,
) (bool, error) {

	const query = `
	SELECT EXISTS(
		SELECT 1
		FROM event_members
		WHERE event_id = $1
		AND profile_id = $2
	)
	`

	var exists bool

	err := r.db.QueryRowContext(
		ctx,
		query,
		event_id,
		profile_id,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("exists member: %w", err)
	}

	return exists, nil
}

// 現在参加人数取得
func (r *eventJoinPostgres) CountMembers(
	ctx context.Context,
	event_id uuid.UUID,
) (int, error) {

	const query = `
	SELECT COUNT(*)
	FROM event_members
	WHERE event_id = $1
	`

	var count int

	err := r.db.QueryRowContext(
		ctx,
		query,
		event_id,
	).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("count members: %w", err)
	}

	return count, nil
}

// 定員取得
func (r *eventJoinPostgres) GetCapacity(
	ctx context.Context,
	event_id uuid.UUID,
) (int, error) {

	const query = `
	SELECT capacity
	FROM events
	WHERE id = $1
	`

	var capacity sql.NullInt32

	err := r.db.QueryRowContext(
		ctx,
		query,
		event_id,
	).Scan(&capacity)

	if err != nil {
		return 0, fmt.Errorf("get capacity: %w", err)
	}

	// NULLなら定員なし
	if !capacity.Valid {
		return 0, nil
	}

	return int(capacity.Int32), nil
}

// 参加登録
func (r *eventJoinPostgres) Join(
	ctx context.Context,
	member *model.EventMember,
) error {

	const query = `
	INSERT INTO event_members(
		id,
		event_id,
		profile_id,
		username,
		mail_address
	)
	VALUES(
		gen_random_uuid(),
		$1,
		$2,
		$3,
		$4
	)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		member.event_id,
		member.profile_id,
		member.username,
		member.mail_address,
	)

	if err != nil {
		return fmt.Errorf("join event: %w", err)
	}

	return nil
}