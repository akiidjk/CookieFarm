CREATE TABLE IF NOT EXISTS flags (
    flag_code VARCHAR(255) PRIMARY KEY,
    service_name VARCHAR(255) NOT NULL,
    service_port INTEGER NOT NULL,
    submit_time INTEGER,  -- Unix timestamp
    response_time INTEGER, -- Unix timestamp
    status VARCHAR(255) NOT NULL,
    team_id INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_flags_submit_time
  ON flags(submit_time);
