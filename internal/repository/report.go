package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/GokujyouKaisennDonnburi/NatuEve_API/internal/model"
)

// ReportRepository は reports テーブルへのアクセスを抽象化する。
type ReportRepository interface {
	// Create はレポートを関連テーブルとともにトランザクション内で一括登録する。
	Create(ctx context.Context, r *model.NewReport) (model.CreateReportResponse, error)
}

// ReportPostgres は ReportRepository の PostgreSQL 実装。
type reportPostgres struct {
	db *sql.DB
}

// Createはレポート関連テーブル(画像, PDF)を含めて一括登録する。
func (r *reportPostgres) Create(ctx context.Context, report *model.NewReport) (model.CreateReportResponse, error) {
	// トランザクション開始
	tx, err := r.db.BeginTx(ctx, nil)

	// トランザクション開始に失敗した場合はエラーを返す
	if err != nil {
		return model.CreateReportResponse{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// reportsテーブルにインサートし、IDと作成日時を取得する
	const insertReport = `
	INSERT INTO reports (event_id, content)
	VALUES ($1, $2)
	RETURNING id, created_at
	`

	// レスポンス用の変数を宣言
	var resp model.CreateReportResponse
	err = tx.QueryRowContext(ctx, insertReport,
		report.EventID,
		report.Content).Scan(&resp.ReportID, &resp.CreatedAt)
	if err != nil {
		return model.CreateReportResponse{}, fmt.Errorf("insert report: %w", err)
	}

	// var reportID string
	// var createdAt string
	// if err := tx.QueryRowContext(ctx, insertReport, report.EventID, report.Content).Scan(&reportID, &createdAt); err != nil {
	// 	return model.CreateReportResponse{}, fmt.Errorf("insert report: %w", err)
	// }

	// 画像オブジェクトキーが存在する場合、関連テーブルにインサートする
	if len(report.ImageObjectKeys) > 0 {
		const insertImage = `
		INSERT INTO report_images (report_id, object_key)
		VALUES ($1, $2)
		`
		for _, objectKey := range report.ImageObjectKeys {
			if _, err := tx.ExecContext(ctx, insertImage, resp.ReportID, objectKey); err != nil {
				return model.CreateReportResponse{}, fmt.Errorf("insert report image: %w", err)
			}
		}
	}

	// PDFオブジェクトキーが存在する場合、関連テーブルにインサートする
	if len(report.PdfObjectKeys) > 0 {
		const insertPDF = `
		INSERT INTO report_pdfs (report_id, object_key)
		VALUES ($1, $2)
		`
		for _, objectKey := range report.PdfObjectKeys {
			if _, err := tx.ExecContext(ctx, insertPDF, resp.ReportID, objectKey); err != nil {
				return model.CreateReportResponse{}, fmt.Errorf("insert report pdf: %w", err)
			}
		}
	}

	// トランザクションをコミットする
	if err := tx.Commit(); err != nil {
		return model.CreateReportResponse{}, fmt.Errorf("commit transaction: %w", err)
	}

	return resp, nil
}
