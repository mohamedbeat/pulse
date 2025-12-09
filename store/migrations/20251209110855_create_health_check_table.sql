-- +goose Up
-- +goose StatementBegin


CREATE TABLE http_check_results (
    -- Basic endpoint info
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    method VARCHAR(10) DEFAULT 'GET',
    interval_seconds INTEGER DEFAULT 0,
    timeout_seconds INTEGER DEFAULT 0,
    expected_status INTEGER DEFAULT 0,
    
    -- Additional config (optional)
    headers JSONB DEFAULT '{}'::jsonb,
    body TEXT,
    must_match_status BOOLEAN DEFAULT true,
    body_contains TEXT,
    body_regex TEXT,
    
    -- Results
    status VARCHAR(20) NOT NULL CHECK (status IN ('up', 'degraded', 'down', 'unreachable')),
    status_code INTEGER,
    content_length INTEGER,
    response_headers JSONB,
    response_body TEXT,
    error_message TEXT,
    duration_ms INTEGER NOT NULL,
    
    -- Timing metrics (optional - can be added later)
    -- dns_duration_ms INTEGER,
    -- tls_duration_ms INTEGER,
    -- connect_duration_ms INTEGER,
    -- first_byte_duration_ms INTEGER,
    -- download_duration_ms INTEGER,
    -- ssl_valid BOOLEAN,
    -- ssl_expiry_days INTEGER,
    
    -- Timestamps
    checked_at TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
   -- INDEX idx_http_check_name (name),
   --  INDEX idx_http_check_status (status),
   --  INDEX idx_http_check_checked_at (checked_at DESC),
   --  INDEX idx_http_check_created_at (created_at DESC),   
    );

-- Indexes for performance
CREATE INDEX idx_http_check_name ON http_check_results (name);
CREATE INDEX idx_http_check_status ON http_check_results (status);
CREATE INDEX idx_http_check_checked_at ON http_check_results (checked_at DESC);
CREATE INDEX idx_http_check_created_at ON http_check_results (created_at DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
    Drop table http_check_results 
-- +goose StatementEnd
