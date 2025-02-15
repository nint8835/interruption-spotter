-- name: InsertStat :exec
INSERT INTO spot_instance_stats (
        region,
        operating_system,
        instance_type,
        interruption_level,
        interruption_level_label,
        observed_time
    )
VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP);

-- name: GetCurrentInterruptionLevels :many
SELECT region,
    operating_system,
    instance_type,
    MAX(observed_time),
    interruption_level,
    interruption_level_label
FROM spot_instance_stats
GROUP BY region,
    operating_system,
    instance_type;

-- name: GetInterruptionChanges :many
SELECT id,
    region,
    instance_type,
    operating_system,
    interruption_level,
    interruption_level_label,
    LAG(interruption_level) OVER (
        PARTITION BY region,
        operating_system,
        instance_type
        ORDER BY observed_time
    ) AS last_interruption_level,
    LAG(interruption_level_label) OVER (
        PARTITION BY region,
        operating_system,
        instance_type
        ORDER BY observed_time
    ) AS last_interruption_level_label,
    observed_time
FROM spot_instance_stats
WHERE region IN (sqlc.slice('regions'))
    AND instance_type IN (sqlc.slice('instance_types'))
    AND operating_system IN (sqlc.slice('operating_systems'))
ORDER BY observed_time DESC
