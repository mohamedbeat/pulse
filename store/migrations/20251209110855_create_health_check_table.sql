-- +goose Up
-- +goose StatementBegin

CREATE TABLE http_check_results (
    -- Basic endpoint info
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    method TEXT DEFAULT 'GET',
    interval_seconds INTEGER DEFAULT 0,
    timeout_seconds INTEGER DEFAULT 0,
    expected_status INTEGER DEFAULT 0,
    
    -- Additional config (optional)
    headers TEXT DEFAULT '{}',
    body TEXT,
    must_match_status INTEGER DEFAULT 1, -- 1 for true, 0 for false
    body_contains TEXT,
    body_regex TEXT,
    
    -- Results
    status TEXT NOT NULL CHECK (status IN ('up', 'degraded', 'down', 'unreachable')),
    status_code INTEGER,
    content_length INTEGER,
    response_headers TEXT,
    response_body TEXT,
    error_message TEXT,
    duration_ms INTEGER NOT NULL,
    
    -- Timestamps
    checked_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_http_check_name ON http_check_results (name);
CREATE INDEX idx_http_check_status ON http_check_results (status);
CREATE INDEX idx_http_check_checked_at ON http_check_results (checked_at DESC);
CREATE INDEX idx_http_check_created_at ON http_check_results (created_at DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS http_check_results;
-- +goose StatementEnd
