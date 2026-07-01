package service

import (
	"context"
	"errors"
	"strings"

	"github.com/GokujyouKaisennDonnburi/NatuEve_API/internal/model"
	"github.com/GokujyouKaisennDonnburi/NatuEve_API/internal/repository"
)

// maxProfileFieldLen はプロフィール文字列フィールドの最大長（rune 単位）。
const maxProfileFieldLen = 255

// AuthenticatedUser は認証済みユーザーの入力 DTO。
//
// middleware への逆依存を避けるため service 側に定義する。
// handler が middleware の AuthUser からこの型へ詰め替える。
type AuthenticatedUser struct {
	ID          string
	Email       string
	DisplayName string
	AvatarURL   string
	Description string
}

// ProfileService はプロフィールに関するビジネスロジックを提供する。
type ProfileService struct {
	repo repository.ProfileRepository
}

// NewProfileService は ProfileService を生成する。
func NewProfileService(repo repository.ProfileRepository) *ProfileService {
	return &ProfileService{repo: repo}
}

// GetOrCreate は認証ユーザーのプロフィールを取得する。存在しなければ作成する。
//
// 認証情報(メール・表示名・アバター)で常に upsert するため、
// Supabase 側のプロフィール変更も次回アクセス時に反映される。
func (s *ProfileService) GetOrCreate(ctx context.Context, u AuthenticatedUser) (*model.Profile, error) {
	p := &model.Profile{
		ID:          u.ID,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		AvatarURL:   u.AvatarURL,
		Description: u.Description,
	}
	if err := s.repo.Upsert(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// ErrProfileNotFound はプロフィールが存在しない場合に返されるエラー。
var ErrProfileNotFound = errors.New("profile not found")

// GetByID はユーザーIDをもとにプロフィールを取得する。
func (s *ProfileService) GetByID(ctx context.Context, id string) (*model.Profile, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrProfileNotFound) {
			return nil, ErrProfileNotFound
		}
		return nil, err
	}
	return p, nil
}

// UpdateMyProfile は本人のプロフィールを部分更新する。
//
// req の各フィールドは nil なら変更せず、非 nil ならその値へ設定する。
// これにより description は空文字へのリセットが可能。表示名は指定時に空不可。
// 検証エラーは *ValidationError で返す（handler 層で 400 に変換）。
func (s *ProfileService) UpdateMyProfile(ctx context.Context, userID string, req model.UpdateProfileRequest) (*model.Profile, error) {
	// まず現在のプロフィール取得
	p, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 部分更新（nil=未指定なので触らない）。値は trim して検証する。
	if req.DisplayName != nil {
		name := strings.TrimSpace(*req.DisplayName)
		if name == "" {
			return nil, &ValidationError{Message: "表示名は空にできません"}
		}
		if len([]rune(name)) > maxProfileFieldLen {
			return nil, &ValidationError{Message: "表示名は255文字以内で入力してください"}
		}
		p.DisplayName = name
	}
	if req.Description != nil {
		desc := strings.TrimSpace(*req.Description)
		if len([]rune(desc)) > maxProfileFieldLen {
			return nil, &ValidationError{Message: "自己紹介は255文字以内で入力してください"}
		}
		p.Description = desc
	}

	// DB更新は編集専用の Update を使う（Upsert は get-or-create 用で編集値を上書きしない）。
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}
