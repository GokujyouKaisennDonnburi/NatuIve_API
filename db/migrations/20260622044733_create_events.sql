-- +goose Up
CREATE TABLE events (
    uuid UUID PRIMARY KEY,
    user_id UUID REFERENCES users(uuid),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    location VARCHAR(255),
    event_date TIMESTAMP NOT NULL,
	capacity INTEGER,
	external_url VARCHAR(255),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose Down
-- (ロールバック時はeventsを削除する)
DROP TABLE events;
