-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pull_requests(
    id SERIAL PRIMARY KEY,
    title VARCHAR(64) NOT NULL,
    author_id INTEGER REFERENCES users(id)
    status VARCHAR(16) NOT NULL CHECK (status IN('OPEN','MERGED')),
    reviewer_ids INTEGER[] DEFAULT '{}',
    created_at TIMESTAMPZ DEFAULT NOW(),
    merged_at TIMESTAMPZ
);
ALTER TABLE pull_requests ADD CONSTRAINT reviewer_count_check CHECK(cardinality(reviewer_ids)<=2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pull_requests;
-- +goose StatementEnd
