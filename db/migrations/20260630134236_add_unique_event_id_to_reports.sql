-- +goose Up
-- 1 イベントにつきレポートは 1 件とする仕様を DB で保証する。
-- 既存データに event_id 重複があると失敗するため、移行前に重複解消が必要。
ALTER TABLE reports ADD CONSTRAINT reports_event_id_unique UNIQUE (event_id);

-- +goose Down
ALTER TABLE reports DROP CONSTRAINT reports_event_id_unique;
