CREATE TABLE IF NOT EXISTS flags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    flag_code VARCHAR(255) NOT NULL UNIQUE,
    service_name VARCHAR(255) NOT NULL,
    port_service INTEGER NOT NULL,
    submit_time INTEGER,  -- Unix timestamp
    response_time INTEGER, -- Unix timestamp
    msg VARCHAR(255) NOT NULL DEFAULT '',
    status INTEGER NOT NULL DEFAULT 0,
    team_id INTEGER NOT NULL,
    username VARCHAR(255) NOT NULL DEFAULT '',
    exploit_name VARCHAR(255) NOT NULL DEFAULT '',
    deleted_at DATETIME NULL
);

CREATE INDEX IF NOT EXISTS idx_flags_deleted_at ON flags(deleted_at, submit_time DESC, id DESC);
CREATE INDEX IF NOT EXISTS idx_flags_submit_time ON flags(submit_time);
CREATE INDEX IF NOT EXISTS idx_flags_team_status ON flags(team_id, status);
CREATE INDEX IF NOT EXISTS idx_flags_submit_time_status ON flags(submit_time, status);
CREATE INDEX IF NOT EXISTS idx_flags_exploit_name ON flags(exploit_name);

CREATE TABLE IF NOT EXISTS exploits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    hash VARCHAR(255) NOT NULL,
    submit_time INTEGER NOT NULL,  -- Unix timestamp
    username VARCHAR(255) NOT NULL DEFAULT 'cookie',
    version INTEGER NOT NULL DEFAULT '1'
);
