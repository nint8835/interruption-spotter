-- name: InsertStat :exec
INSERT INTO spot_instance_stats (
        region,
        operating_system,
        instance_type,
        interruption_level,
        observed_time
    )
VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP);
