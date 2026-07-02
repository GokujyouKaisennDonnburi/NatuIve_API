-- +goose Up
ALTER TABLE event_members
    ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD CONSTRAINT event_members_event_id_profile_id_key UNIQUE (event_id, profile_id);

-- +goose Down
ALTER TABLE event_members
    DROP CONSTRAINT event_members_event_id_profile_id_key,
    DROP COLUMN created_at;
