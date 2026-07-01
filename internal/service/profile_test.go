package service

import (
	"context"
	"errors"
	"testing"

	"github.com/GokujyouKaisennDonnburi/NatuEve_API/internal/model"
)

// stubProfileRepository は ProfileRepository のテスト用スタブ。
type stubProfileRepository struct {
	upsertCalled bool
	updateCalled bool
	gotProfile   *model.Profile
	upsertErr    error
	updateErr    error

	getProfile *model.Profile
	getErr     error
}

func (s *stubProfileRepository) GetByID(_ context.Context, _ string) (*model.Profile, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	return s.getProfile, nil
}

func (s *stubProfileRepository) Upsert(_ context.Context, p *model.Profile) error {
	s.upsertCalled = true
	s.gotProfile = p
	return s.upsertErr
}

func (s *stubProfileRepository) Update(_ context.Context, p *model.Profile) error {
	s.updateCalled = true
	s.gotProfile = p
	return s.updateErr
}

func TestProfileServiceGetOrCreate(t *testing.T) {
	user := AuthenticatedUser{
		ID:          "d290f1ee-6c54-4b01-90e6-d701748f0851",
		Email:       "user@example.com",
		DisplayName: "なちゅいべ太郎",
		AvatarURL:   "https://example.com/a.png",
		Description: "イベントを楽しむのが好きです。",
	}

	tests := []struct {
		name      string
		upsertErr error
		wantErr   bool
	}{
		{name: "正常: Upsert が呼ばれプロフィールを返す", upsertErr: nil, wantErr: false},
		{name: "異常: Upsert のエラーが伝播する", upsertErr: errors.New("db error"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &stubProfileRepository{upsertErr: tt.upsertErr}
			svc := NewProfileService(repo)

			got, err := svc.GetOrCreate(context.Background(), user)

			if !repo.upsertCalled {
				t.Fatalf("Upsert が呼ばれていない")
			}
			if tt.wantErr {
				if err == nil {
					t.Fatalf("エラーを期待したが nil だった")
				}
				return
			}
			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}
			if got.ID != user.ID || got.Email != user.Email ||
				got.DisplayName != user.DisplayName || got.AvatarURL != user.AvatarURL || got.Description != user.Description {
				t.Errorf("返り値が入力と一致しない: got %+v, want %+v", got, user)
			}
		})
	}
}

func TestProfileServiceUpdateMyProfile(t *testing.T) {
	base := &model.Profile{
		ID:          "user-1",
		Email:       "user@example.com",
		DisplayName: "old name",
		Description: "old desc",
	}

	tests := []struct {
		name      string
		req       model.UpdateProfileRequest
		getErr    error
		updateErr error
		wantErr   bool
	}{
		{
			name: "正常: 全項目更新",
			req: model.UpdateProfileRequest{
				DisplayName: "new name",
				Description: "new desc",
			},
		},
		{
			name: "正常: 部分更新",
			req: model.UpdateProfileRequest{
				DisplayName: "new name",
			},
		},
		{
			name:    "異常: GetByIDエラー",
			getErr:  errors.New("db error"),
			wantErr: true,
		},
		{
			name:      "異常: Updateエラー",
			updateErr: errors.New("db error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &stubProfileRepository{
				getProfile: base,
				getErr:     tt.getErr,
				updateErr:  tt.updateErr,
			}

			svc := NewProfileService(repo)

			got, err := svc.UpdateMyProfile(context.Background(), "user-1", tt.req)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("エラーが発生することを期待しましたが、nil でした")
				}
				if got != nil {
					t.Fatalf("エラー時は結果は nil の想定です")
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラーが発生しました: %v", err)
			}

			if !repo.updateCalled {
				t.Fatalf("Update が呼び出されていません")
			}

			if got == nil {
				t.Fatalf("結果が nil でした")
			}

			// 更新確認
			if got.DisplayName != tt.req.DisplayName && tt.req.DisplayName != "" {
				t.Errorf("DisplayName が一致しません")
			}
			if got.Description != tt.req.Description && tt.req.Description != "" {
				t.Errorf("Description が一致しません")
			}
		})
	}
}
