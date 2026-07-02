package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/GokujyouKaisennDonnburi/NatuEve_API/internal/model"
)

// stubEventJoinRepository は EventJoinRepository のテスト用スタブ。
type stubEventJoinRepository struct {
	// ExistsEvent 返却値。
	existsEventResult bool
	existsEventErr    error
	// ExistsMember 返却値。
	existsMemberResult bool
	existsMemberErr    error
	// CountMembers 返却値。
	countMembersResult int
	countMembersErr    error
	// GetCapacity 返却値。
	getCapacityResult int
	getCapacityErr    error
	// Join 返却値（joinCreatedAt は成功時に member.CreatedAt へセットする）。
	joinCreatedAt time.Time
	joinErr       error
	// 呼び出し時に Join へ渡された引数を記録する。
	gotMember *model.EventMember
}

func (s *stubEventJoinRepository) ExistsEvent(_ context.Context, _ uuid.UUID) (bool, error) {
	return s.existsEventResult, s.existsEventErr
}

func (s *stubEventJoinRepository) ExistsMember(_ context.Context, _, _ uuid.UUID) (bool, error) {
	return s.existsMemberResult, s.existsMemberErr
}

func (s *stubEventJoinRepository) CountMembers(_ context.Context, _ uuid.UUID) (int, error) {
	return s.countMembersResult, s.countMembersErr
}

func (s *stubEventJoinRepository) GetCapacity(_ context.Context, _ uuid.UUID) (int, error) {
	return s.getCapacityResult, s.getCapacityErr
}

func (s *stubEventJoinRepository) Join(_ context.Context, member *model.EventMember) error {
	s.gotMember = member
	if s.joinErr != nil {
		return s.joinErr
	}
	member.CreatedAt = s.joinCreatedAt
	return nil
}

// assertNotFoundError はテストヘルパー: err が *NotFoundError であることを確認する。
func assertNotFoundError(t *testing.T, err error) *NotFoundError {
	t.Helper()
	if err == nil {
		t.Fatal("NotFoundError を期待したが nil だった")
	}
	var nfe *NotFoundError
	if !errors.As(err, &nfe) {
		t.Fatalf("*NotFoundError を期待したが %T だった: %v", err, err)
	}
	return nfe
}

// assertConflictError はテストヘルパー: err が *ConflictError であることを確認する。
func assertConflictError(t *testing.T, err error) *ConflictError {
	t.Helper()
	if err == nil {
		t.Fatal("ConflictError を期待したが nil だった")
	}
	var ce *ConflictError
	if !errors.As(err, &ce) {
		t.Fatalf("*ConflictError を期待したが %T だった: %v", err, err)
	}
	return ce
}

// defaultJoinStub は正常系テスト用のデフォルトスタブを返す。
func defaultJoinStub(createdAt time.Time) *stubEventJoinRepository {
	return &stubEventJoinRepository{
		existsEventResult:  true,
		existsMemberResult: false,
		getCapacityResult:  30,
		countMembersResult: 5,
		joinCreatedAt:      createdAt,
	}
}

func TestEventJoinServiceJoin(t *testing.T) {
	// 固定 UUID でテストの再現性を確保する。
	eventID := uuid.MustParse("a1b2c3d4-e5f6-7890-abcd-ef1234567890")
	profileID := uuid.MustParse("b2c3d4e5-f6a7-8901-bcde-f23456789012")
	createdAt := time.Date(2026, 6, 26, 4, 54, 35, 0, time.UTC)

	validReq := model.JoinEventRequest{
		Username:    "山田太郎",
		MailAddress: "yamada@example.com",
	}

	tests := []struct {
		name         string
		stub         *stubEventJoinRepository
		req          model.JoinEventRequest
		wantValErr   bool
		wantNotFound bool
		wantConflict bool
		wantErr      bool
		// 正常系: レスポンスの全フィールドを検証する。
		checkResp func(t *testing.T, resp model.JoinEventResponse)
		// 正常系: repo に渡った EventMember の内容を検証する。
		checkMember func(t *testing.T, m *model.EventMember)
	}{
		// --- 正常系 ---
		{
			name: "正常: レスポンスの全フィールドが正しく返る",
			stub: defaultJoinStub(createdAt),
			req:  validReq,
			checkResp: func(t *testing.T, resp model.JoinEventResponse) {
				t.Helper()
				if resp.EventID != eventID {
					t.Errorf("EventID: got %v, want %v", resp.EventID, eventID)
				}
				if resp.ProfileID != profileID {
					t.Errorf("ProfileID: got %v, want %v", resp.ProfileID, profileID)
				}
				if resp.Username != "山田太郎" {
					t.Errorf("Username: got %q, want %q", resp.Username, "山田太郎")
				}
				if resp.MailAddress != "yamada@example.com" {
					t.Errorf("MailAddress: got %q, want %q", resp.MailAddress, "yamada@example.com")
				}
				if !resp.CreatedAt.Equal(createdAt) {
					t.Errorf("CreatedAt: got %v, want %v", resp.CreatedAt, createdAt)
				}
			},
			checkMember: func(t *testing.T, m *model.EventMember) {
				t.Helper()
				if m == nil {
					t.Fatal("gotMember が nil")
				}
				if m.EventID != eventID {
					t.Errorf("EventMember.EventID: got %v, want %v", m.EventID, eventID)
				}
				if m.ProfileID != profileID {
					t.Errorf("EventMember.ProfileID: got %v, want %v", m.ProfileID, profileID)
				}
				if m.Username != "山田太郎" {
					t.Errorf("EventMember.Username: got %q, want %q", m.Username, "山田太郎")
				}
				if m.MailAddress != "yamada@example.com" {
					t.Errorf("EventMember.MailAddress: got %q, want %q", m.MailAddress, "yamada@example.com")
				}
			},
		},
		{
			name: "正常: username・mailAddress の TrimSpace が反映される",
			stub: defaultJoinStub(createdAt),
			req: model.JoinEventRequest{
				Username:    "  山田太郎  ",
				MailAddress: "  yamada@example.com  ",
			},
			checkMember: func(t *testing.T, m *model.EventMember) {
				t.Helper()
				if m.Username != "山田太郎" {
					t.Errorf("Username trim: got %q, want %q", m.Username, "山田太郎")
				}
				if m.MailAddress != "yamada@example.com" {
					t.Errorf("MailAddress trim: got %q, want %q", m.MailAddress, "yamada@example.com")
				}
			},
		},
		{
			name: "正常: capacity 0 は定員なしのため参加成功",
			stub: func() *stubEventJoinRepository {
				s := defaultJoinStub(createdAt)
				s.getCapacityResult = 0
				s.countMembersResult = 999
				return s
			}(),
			req: validReq,
		},
		// --- バリデーションエラー ---
		{
			name:       "異常: username が空",
			stub:       defaultJoinStub(createdAt),
			req:        model.JoinEventRequest{Username: "", MailAddress: "yamada@example.com"},
			wantValErr: true,
		},
		{
			name:       "異常: username が空白のみ",
			stub:       defaultJoinStub(createdAt),
			req:        model.JoinEventRequest{Username: "   ", MailAddress: "yamada@example.com"},
			wantValErr: true,
		},
		{
			name: "異常: username が 256 文字",
			stub: defaultJoinStub(createdAt),
			req: func() model.JoinEventRequest {
				runes := make([]rune, 256)
				for i := range runes {
					runes[i] = 'あ'
				}
				return model.JoinEventRequest{
					Username:    string(runes),
					MailAddress: "yamada@example.com",
				}
			}(),
			wantValErr: true,
		},
		{
			name:       "異常: mailAddress が空",
			stub:       defaultJoinStub(createdAt),
			req:        model.JoinEventRequest{Username: "山田太郎", MailAddress: ""},
			wantValErr: true,
		},
		{
			name:       "異常: mailAddress の形式が不正",
			stub:       defaultJoinStub(createdAt),
			req:        model.JoinEventRequest{Username: "山田太郎", MailAddress: "not-an-email"},
			wantValErr: true,
		},
		// --- ビジネスチェックエラー ---
		{
			name: "異常: イベントが存在しない（NotFoundError）",
			stub: func() *stubEventJoinRepository {
				s := defaultJoinStub(createdAt)
				s.existsEventResult = false
				return s
			}(),
			req:          validReq,
			wantNotFound: true,
		},
		{
			name: "異常: 既に参加済み（ConflictError）",
			stub: func() *stubEventJoinRepository {
				s := defaultJoinStub(createdAt)
				s.existsMemberResult = true
				return s
			}(),
			req:          validReq,
			wantConflict: true,
		},
		{
			name: "異常: 定員超過（ConflictError）",
			stub: func() *stubEventJoinRepository {
				s := defaultJoinStub(createdAt)
				s.getCapacityResult = 30
				s.countMembersResult = 30
				return s
			}(),
			req:          validReq,
			wantConflict: true,
		},
		// --- リポジトリエラー伝播 ---
		{
			name: "異常: repo.ExistsEvent がエラーを返す",
			stub: func() *stubEventJoinRepository {
				s := defaultJoinStub(createdAt)
				s.existsEventErr = errors.New("db error")
				return s
			}(),
			req:     validReq,
			wantErr: true,
		},
		{
			name: "異常: repo.Join がエラーを返す",
			stub: func() *stubEventJoinRepository {
				s := defaultJoinStub(createdAt)
				s.joinErr = errors.New("db error")
				return s
			}(),
			req:     validReq,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewEventJoinService(tt.stub)

			resp, err := svc.Join(context.Background(), eventID, profileID, tt.req)

			switch {
			case tt.wantValErr:
				_ = assertValidationError(t, err)
				return
			case tt.wantNotFound:
				_ = assertNotFoundError(t, err)
				return
			case tt.wantConflict:
				_ = assertConflictError(t, err)
				return
			case tt.wantErr:
				if err == nil {
					t.Fatal("エラーを期待したが nil だった")
				}
				return
			}

			assertNoErr(t, err)

			if tt.checkResp != nil {
				tt.checkResp(t, resp)
			}
			if tt.checkMember != nil {
				tt.checkMember(t, tt.stub.gotMember)
			}
		})
	}
}
