-- name: InsertStat :exec
INSERT INTO spot_instance_stats (
        region,
        operating_system,
        instance_type,
        interruption_level,
        observed_time
    )
VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP);

-- name: GetCurrentInterruptionLevel :one
SELECT interruption_level
FROM spot_instance_stats
WHERE region = ?
    AND operating_system = ?
    AND instance_type = ?
ORDER BY observed_time DESC
LIMIT 1;
