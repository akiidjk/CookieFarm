-- name: GetFlagByCode :one
SELECT *
FROM flags
WHERE flag_code = ?
LIMIT 1;

-- name: GetFlagsByTeam :many
SELECT *
FROM flags
WHERE team_id = ?
ORDER BY submit_time DESC
LIMIT ? OFFSET ?;

-- name: GetAllFlags :many
SELECT *
FROM flags
ORDER BY submit_time DESC;

-- name: CountFlags :one
SELECT COUNT(*)
FROM flags
LIMIT 1;

-- name: GetFirstNFlags :many
SELECT *
FROM flags
ORDER BY submit_time DESC
LIMIT ?;

-- name: GetUnsubmittedFlags :many
SELECT *
FROM flags
WHERE status = 0
ORDER BY submit_time ASC
LIMIT ?;

-- name: GetPagedFlags :many
SELECT *
FROM flags
ORDER BY submit_time DESC
LIMIT ? OFFSET ?;

-- name: GetAllFlagCodes :many
SELECT flag_code FROM flags;

-- name: GetFirstNFlagCodes :many
SELECT flag_code FROM flags
LIMIT ?;

-- name: GetUnsubmittedFlagCodes :many
SELECT flag_code FROM flags
WHERE status = 0
LIMIT ?;

-- name: GetFilteredFlags :many
SELECT * FROM flags
WHERE
    (team_id = sqlc.narg('team_id') OR sqlc.narg('team_id') IS NULL)
    AND (status = sqlc.narg('status') OR sqlc.narg('status') IS NULL)
    AND (
        sqlc.narg('search') IS NULL
        OR (
            (sqlc.narg('search_field') = 'flag_code'    AND flag_code    LIKE sqlc.narg('search_like'))
            OR (sqlc.narg('search_field') = 'service_name' AND service_name LIKE sqlc.narg('search_like'))
            OR (sqlc.narg('search_field') = 'exploit_name' AND exploit_name LIKE sqlc.narg('search_like'))
            OR (sqlc.narg('search_field') = 'msg'          AND msg          LIKE sqlc.narg('search_like'))
            OR (sqlc.narg('search_field') = 'all' AND (
                flag_code    LIKE sqlc.narg('search_like')
                OR service_name  LIKE sqlc.narg('search_like')
                OR exploit_name  LIKE sqlc.narg('search_like')
                OR msg           LIKE sqlc.narg('search_like')
                OR CAST(team_id AS TEXT) LIKE sqlc.narg('search_like')
            ))
            OR (sqlc.narg('search_field') IS NULL AND flag_code LIKE sqlc.narg('search_like'))
        )
)
ORDER BY submit_time DESC
LIMIT sqlc.narg('limit') OFFSET sqlc.narg('offset');


-- name: CountFilteredFlags :one
SELECT COUNT(*) FROM flags
WHERE
    (team_id = sqlc.narg('team_id') OR sqlc.narg('team_id') IS NULL)
    AND (status = sqlc.narg('status') OR sqlc.narg('status') is NULL)
    AND (service_name = sqlc.narg('service_name') OR sqlc.narg('service_name') IS NULL)
    AND (
        sqlc.narg('search') IS NULL
        OR (
            (sqlc.narg('search_field') = 'flag_code'    AND flag_code    LIKE sqlc.narg('search_like'))
            OR (sqlc.narg('search_field') = 'service_name' AND service_name LIKE sqlc.narg('search_like'))
            OR (sqlc.narg('search_field') = 'exploit_name' AND exploit_name LIKE sqlc.narg('search_like'))
            OR (sqlc.narg('search_field') = 'msg'          AND msg          LIKE sqlc.narg('search_like'))
            OR (sqlc.narg('search_field') = 'all' AND (
                flag_code    LIKE sqlc.narg('search_like')
                OR service_name  LIKE sqlc.narg('search_like')
                OR exploit_name  LIKE sqlc.narg('search_like')
                OR msg           LIKE sqlc.narg('search_like')
                OR CAST(team_id AS TEXT) LIKE sqlc.narg('search_like')
            ))
            OR (sqlc.narg('search_field') IS NULL AND flag_code LIKE sqlc.narg('search_like'))
    )
);

-- name: GetPagedFlagCodes :many
SELECT flag_code FROM flags
LIMIT ? OFFSET ?;

-- name: AddFlag :exec
INSERT OR IGNORE INTO flags(
	flag_code, service_name, port_service,
	submit_time, response_time, status,
	team_id, msg, username, exploit_name
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateFlagStatusByCode :exec
UPDATE flags
SET
	status = ?,
	msg = ?,
	response_time = ?
WHERE flag_code = ?;

-- name: DeleteFlagByCode :exec
DELETE FROM flags
WHERE flag_code = ?;

-- name: DeleteFlagByTTL :execrows
DELETE FROM flags
WHERE response_time < (CAST(strftime('%s', 'now') AS INTEGER) - ?);

-- name: FlagsStats :many
SELECT
    team_id,
    COUNT(*) AS total_flags,
    SUM(CASE WHEN status = 1 THEN 1 ELSE 0 END) AS accepted_flags,
    SUM(CASE WHEN status = 2 THEN 1 ELSE 0 END) AS denied_flags,
    SUM(CASE WHEN status = 0 THEN 1 ELSE 0 END) AS unsubmitted_flags,
    SUM(CASE WHEN status = 3 THEN 1 ELSE 0 END) AS error_flags,
    SUM(CASE WHEN status = 4 THEN 1 ELSE 0 END) AS not_valid_flags
FROM flags
GROUP BY team_id
ORDER BY team_id;

-- name: FlagsTickStats :many
SELECT
    (submit_time / ?) * ? AS timestamp,
    COUNT(*) AS total,
    SUM(CASE WHEN status = 0 THEN 1 ELSE 0 END) AS queued,
    SUM(CASE WHEN status = 1 THEN 1 ELSE 0 END) AS accepted,
    SUM(CASE WHEN status = 2 THEN 1 ELSE 0 END) AS denied,
    SUM(CASE WHEN status = 3 THEN 1 ELSE 0 END) AS error,
    SUM(CASE WHEN status = 4 THEN 1 ELSE 0 END) AS invalid
FROM flags
WHERE submit_time > 0
GROUP BY timestamp
ORDER BY timestamp;

-- name: FlagsExploitShare :many
SELECT
    exploit_name,
    COUNT(*) AS value
FROM flags
GROUP BY exploit_name
ORDER BY value DESC, exploit_name ASC;

-- EXPLOITS QUERIES

-- name: GetExploitByHash :one
SELECT *
FROM exploits
WHERE hash = ?
LIMIT 1;

-- name: GetExploitsByUsername :many
SELECT *
FROM exploits
WHERE username = ?
ORDER BY submit_time DESC
LIMIT ? OFFSET ?;

-- name: GetAllExploits :many
SELECT *
FROM exploits
ORDER BY submit_time DESC;

-- name: CountExploits :one
SELECT COUNT(*)
FROM exploits
LIMIT 1;

-- name: GetExploitsByName :many
SELECT *
FROM exploits
WHERE name = ?
ORDER BY submit_time DESC;

-- name: CreateExploit :exec

INSERT INTO exploits(name, hash, submit_time, username, version)
VALUES (?, ?, ?, ?, ?);

-- name: DeleteExploitByID :exec
DELETE FROM exploits
WHERE id = ?;
