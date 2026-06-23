-- +goose Up
CREATE TABLE event_costs (
    id UUID PRIMARY KEY,
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    category VARCHAR(255) NOT NULL,
    cost INTEGER NOT NULL
);          

-- +goose Down
-- (ロールバック時は event_costs を削除する)
DROP TABLE event_costs;
