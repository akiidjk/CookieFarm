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

CREATE INDEX IF NOT EXISTS idx_flags_deleted_at ON flags(deleted_at, submit_time DESC, id DESC); -- Used by: GetAllFlags, GetFirstNFlags, CountFlags, GetAllFlagCodes, GetFirstNFlagCodes, GetFilteredFlags (cursor), FlagsTickStats
CREATE INDEX IF NOT EXISTS idx_flags_exploit_name ON flags(exploit_name); -- Used by: FlagsExploitShare
CREATE INDEX IF NOT EXISTS idx_flags_unsubmitted ON flags(deleted_at, status, submit_time ASC); -- Used by: GetUnsubmittedFlags, GetUnsubmittedFlagCodes
CREATE INDEX IF NOT EXISTS idx_flags_deleted_response ON flags(deleted_at, response_time); -- Used by: DeleteFlagByTTL
CREATE INDEX IF NOT EXISTS idx_flags_stats ON flags(deleted_at, team_id, status); -- Used by: FlagsStats
CREATE INDEX IF NOT EXISTS idx_flags_service_name ON flags(service_name, deleted_at); -- Used by: GetFilteredFlags, CountFilteredFlags (when service_name filter is active)

CREATE TABLE IF NOT EXISTS exploits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    hash VARCHAR(255) NOT NULL,
    submit_time INTEGER NOT NULL,  -- Unix timestamp
    username VARCHAR(255) NOT NULL DEFAULT 'cookie',
    version INTEGER NOT NULL DEFAULT '1'
);

CREATE INDEX IF NOT EXISTS idx_exploits_hash ON exploits(hash); -- Used by: GetExploitByHash
CREATE INDEX IF NOT EXISTS idx_exploits_username ON exploits(username, id, submit_time DESC); -- Used by: GetExploitsByUsername
CREATE INDEX IF NOT EXISTS idx_exploits_submit_time ON exploits(submit_time DESC); -- Used by: GetAllExploits
CREATE INDEX IF NOT EXISTS idx_exploits_name ON exploits(name, submit_time DESC); -- Used by: GetExploitsByName
