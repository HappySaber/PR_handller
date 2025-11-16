-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pull_requests(
    id VARCHAR(128) PRIMARY KEY,
    title VARCHAR(64) NOT NULL,
    author_id VARCHAR(64) REFERENCES users(id),
    status VARCHAR(16) NOT NULL CHECK (status IN('OPEN','MERGED')),
    reviewer_ids VARCHAR(64)[] DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    merged_at TIMESTAMPTZ
);
ALTER TABLE pull_requests ADD CONSTRAINT reviewer_count_check CHECK(cardinality(reviewer_ids)<=2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pull_requests;
-- +goose StatementEnd
