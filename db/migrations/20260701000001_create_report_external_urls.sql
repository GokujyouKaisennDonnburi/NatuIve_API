-- +goose Up
CREATE TABLE report_external_urls(
    id UUID PRIMARY KEY,
    report_id UUID REFERENCES reports(id),
    external_url TEXT NOT NULL
);

-- +goose Down
DROP TABLE report_external_urls;
