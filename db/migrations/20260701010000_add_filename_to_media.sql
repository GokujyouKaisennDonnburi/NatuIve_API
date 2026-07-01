-- +goose Up
-- 画像・PDFの元ファイル名を保持する filename カラムを追加する。
-- ダウンロード時のファイル名（Content-Disposition）表示と、UIでの元名表示に使う。
-- 既存行は空文字（未設定）とし、アプリ側は空なら objectkey の basename にフォールバックする。
ALTER TABLE event_images ADD COLUMN filename TEXT NOT NULL DEFAULT '';
ALTER TABLE event_pdfs ADD COLUMN filename TEXT NOT NULL DEFAULT '';
ALTER TABLE report_images ADD COLUMN filename TEXT NOT NULL DEFAULT '';
ALTER TABLE report_pdfs ADD COLUMN filename TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE event_images DROP COLUMN filename;
ALTER TABLE event_pdfs DROP COLUMN filename;
ALTER TABLE report_images DROP COLUMN filename;
ALTER TABLE report_pdfs DROP COLUMN filename;
