CREATE TABLE IF NOT EXISTS flags (
    flag_code VARCHAR(255) PRIMARY KEY,
    service_name VARCHAR(255) NOT NULL,
    port_service INTEGER NOT NULL,
    submit_time INTEGER,  -- Unix timestamp
    response_time INTEGER, -- Unix timestamp
    msg VARCHAR(255) NOT NULL DEFAULT '',
    status INTEGER NOT NULL DEFAULT 0,
    team_id INTEGER NOT NULL,
    username VARCHAR(255) NOT NULL DEFAULT '',
    exploit_name VARCHAR(255) NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_flags_submit_time
  ON flags(submit_time);

CREATE TABLE IF NOT EXISTS exploits (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    hash VARCHAR(255) NOT NULL,
    submit_time INTEGER,  -- Unix timestamp
    username VARCHAR(255) NOT NULL DEFAULT 'cookie',
    version INTEGER NOT NULL DEFAULT '1'
);
