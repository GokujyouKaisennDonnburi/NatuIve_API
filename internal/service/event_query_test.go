package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GokujyouKaisennDonnburi/NatuIve_API/internal/model"
)

// stubEventRepository は EventRepository のテスト用スタブ。
type stubEventRepository struct {
	// 呼び出し時に渡された引数を記録する。
	gotLimit  int
	gotOffset int
	// 返却値。
	results []model.EventSummary
	err     error
}

func (s *stubEventRepository) ListSummaries(_ context.Context, limit, offset int) ([]model.EventSummary, error) {
	s.gotLimit = limit
	s.gotOffset = offset
	return s.results, s.err
}

// makeHelper はテストヘルパー共通処理を担う。
func makeHelper(t *testing.T) {
	t.Helper()
}

func TestEventQueryServiceList_Normalization(t *testing.T) {
	t.Helper()

	// ダミーのサマリーデータ（正規化検証には内容不問）。
	dummyResults := []model.EventSummary{
		{ID: "id-1", Title: "テストイベント", EventDate: time.Now(), CreatedAt: time.Now()},
	}

	tests := []struct {
		name        string
		inputLimit  int
		inputOffset int
		wantLimit   int
		wantOffset  int
		repoErr     error
		wantErr     bool
	}{
		{
			name:        "正常: limit/offset がデフォルト値(0)の場合は default20/0 に正規化",
			inputLimit:  0,
			inputOffset: 0,
			wantLimit:   20,
			wantOffset:  0,
		},
		{
			name:        "正常: limit が負値の場合は default20 に正規化",
			inputLimit:  -5,
			inputOffset: 0,
			wantLimit:   20,
			wantOffset:  0,
		},
		{
			name:        "正常: limit が 100 超過の場合は 100 に丸める",
			inputLimit:  200,
			inputOffset: 0,
			wantLimit:   100,
			wantOffset:  0,
		},
		{
			name:        "正常: limit が最大値ちょうど(100)はそのまま",
			inputLimit:  100,
			inputOffset: 0,
			wantLimit:   100,
			wantOffset:  0,
		},
		{
			name:        "正常: limit が有効範囲内はそのまま",
			inputLimit:  50,
			inputOffset: 10,
			wantLimit:   50,
			wantOffset:  10,
		},
		{
			name:        "正常: offset が負値の場合は 0 に丸める",
			inputLimit:  20,
			inputOffset: -1,
			wantLimit:   20,
			wantOffset:  0,
		},
		{
			name:        "異常: repository のエラーが伝播する",
			inputLimit:  20,
			inputOffset: 0,
			wantLimit:   20,
			wantOffset:  0,
			repoErr:     errors.New("db error"),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			makeHelper(t)

			stub := &stubEventRepository{results: dummyResults, err: tt.repoErr}
			svc := NewEventQueryService(stub)

			got, err := svc.List(context.Background(), tt.inputLimit, tt.inputOffset)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("エラーを期待したが nil だった")
				}
				return
			}
			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}

			// 正規化後の limit / offset が repository に渡されているか確認。
			if stub.gotLimit != tt.wantLimit {
				t.Errorf("limit: got %d, want %d", stub.gotLimit, tt.wantLimit)
			}
			if stub.gotOffset != tt.wantOffset {
				t.Errorf("offset: got %d, want %d", stub.gotOffset, tt.wantOffset)
			}

			// 結果が repository の返却値と一致するか確認。
			if len(got) != len(dummyResults) {
				t.Errorf("件数: got %d, want %d", len(got), len(dummyResults))
			}
		})
	}
}
