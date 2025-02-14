CREATE TABLE IF NOT EXISTS spot_instance_stats (
    id INTEGER PRIMARY KEY,
    region TEXT NOT NULL,
    operating_system TEXT NOT NULL,
    instance_type TEXT NOT NULL,
    interruption_level INTEGER NOT NULL,
    observed_time TIMESTAMP NOT NULL
);
