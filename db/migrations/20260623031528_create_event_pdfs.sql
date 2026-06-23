-- +goose Up
CREATE TABLE event_pdfs (
    id UUID PRIMARY KEY,
    event_id UUID REFERENCES events(id),
    pdf_objectkey VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE event_pdfs;
