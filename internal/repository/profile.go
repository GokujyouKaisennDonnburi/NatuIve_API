package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/GokujyouKaisennDonnburi/NatuEve_API/internal/model"
)

// ErrProfileNotFound はプロフィールが存在しないことを表すセンチネルエラー。
var ErrProfileNotFound = errors.New("profile not found")

// ProfileRepository は profiles テーブルへのアクセスを抽象化する。
type ProfileRepository interface {
	// GetByID は ID でプロフィールを取得する。無ければ ErrProfileNotFound。
	GetByID(ctx context.Context, id string) (*model.Profile, error)
	// Upsert はログイン時の get-or-create 用。存在しなければ作成し、存在すれば
	// email/updated_at のみ更新する（display_name/avatar_url/description は初回のみ
	// 投入し、以後は上書きしない = ユーザー編集を保持する）。実行後、p には DB 上の
	// 最新の状態が反映される。
	Upsert(ctx context.Context, p *model.Profile) error
	// Update はユーザー自身によるプロフィール編集用。display_name/description を更新する。
	// 対象が存在しなければ ErrProfileNotFound。実行後、p に DB 上の最新状態が反映される。
	Update(ctx context.Context, p *model.Profile) error
}

// profilePostgres は ProfileRepository の PostgreSQL 実装。
type profilePostgres struct {
	db *sql.DB
}

// NewProfileRepository は *sql.DB を使う ProfileRepository を生成する。
func NewProfileRepository(db *sql.DB) ProfileRepository {
	return &profilePostgres{db: db}
}

// GetByID は ID でプロフィールを取得する。
func (r *profilePostgres) GetByID(ctx context.Context, id string) (*model.Profile, error) {
	const query = `
		SELECT id, email, display_name, avatar_url, description, created_at, updated_at
		FROM profiles
		WHERE id = $1`

	var (
		p           model.Profile
		displayName sql.NullString
		avatarURL   sql.NullString
		description sql.NullString
	)
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Email, &displayName, &avatarURL, &description, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrProfileNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get profile: %w", err)
	}
	p.DisplayName = displayName.String
	p.AvatarURL = avatarURL.String
	p.Description = description.String
	return &p, nil
}

// Upsert はログイン時の get-or-create。存在しなければ JWT 由来の初期値で作成し、
// 存在すれば email/updated_at のみ更新する。
//
// display_name/avatar_url/description は DO UPDATE SET に含めないため、2 回目以降の
// ログイン（GET /me）では上書きされず、ユーザーが編集した値がそのまま保持される。
// RETURNING で DB 上の最新値を読み戻すため、p には編集後の値が反映される。
func (r *profilePostgres) Upsert(ctx context.Context, p *model.Profile) error {
	const query = `
		INSERT INTO profiles (id, email, display_name, avatar_url, description)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			email      = EXCLUDED.email,
			updated_at = now()
		RETURNING id, email, display_name, avatar_url, description, created_at, updated_at`

	// display_name/avatar_url は nullable なため NullString で受ける。
	// description は NOT NULL DEFAULT '' のため、空文字を NULL 化せずそのまま渡す
	// （nullString で NULL 化すると NOT NULL 制約違反。DEFAULT は値未指定時のみ適用）。
	var (
		displayName sql.NullString
		avatarURL   sql.NullString
		description sql.NullString
	)
	err := r.db.QueryRowContext(ctx, query,
		p.ID, p.Email, nullString(p.DisplayName), nullString(p.AvatarURL), p.Description,
	).Scan(&p.ID, &p.Email, &displayName, &avatarURL, &description, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("upsert profile: %w", err)
	}
	p.DisplayName = displayName.String
	p.AvatarURL = avatarURL.String
	p.Description = description.String
	return nil
}

// Update はユーザー自身によるプロフィール編集。display_name/description を更新する。
//
// avatar_url は Supabase 由来の値を使うためこの経路では更新しない（将来必要になれば
// 専用の更新手段を設ける）。対象行が無ければ ErrProfileNotFound を返す。
func (r *profilePostgres) Update(ctx context.Context, p *model.Profile) error {
	const query = `
		UPDATE profiles
		SET display_name = $2, description = $3, updated_at = now()
		WHERE id = $1
		RETURNING id, email, display_name, avatar_url, description, created_at, updated_at`

	var (
		displayName sql.NullString
		avatarURL   sql.NullString
		description sql.NullString
	)
	err := r.db.QueryRowContext(ctx, query,
		p.ID, nullString(p.DisplayName), p.Description,
	).Scan(&p.ID, &p.Email, &displayName, &avatarURL, &description, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrProfileNotFound
	}
	if err != nil {
		return fmt.Errorf("update profile: %w", err)
	}
	p.DisplayName = displayName.String
	p.AvatarURL = avatarURL.String
	p.Description = description.String
	return nil
}

// nullString は空文字を NULL として扱うための変換を行う。
func nullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}
