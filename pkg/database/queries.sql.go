// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package database

import (
	"context"
	"strings"
	"time"
)

const getCurrentInterruptionLevels = `-- name: GetCurrentInterruptionLevels :many
SELECT region,
    operating_system,
    instance_type,
    MAX(observed_time),
    interruption_level,
    interruption_level_label
FROM spot_instance_stats
GROUP BY region,
    operating_system,
    instance_type
`

type GetCurrentInterruptionLevelsRow struct {
	Region                 string
	OperatingSystem        string
	InstanceType           string
	Max                    interface{}
	InterruptionLevel      int64
	InterruptionLevelLabel string
}

func (q *Queries) GetCurrentInterruptionLevels(ctx context.Context) ([]GetCurrentInterruptionLevelsRow, error) {
	rows, err := q.db.QueryContext(ctx, getCurrentInterruptionLevels)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCurrentInterruptionLevelsRow
	for rows.Next() {
		var i GetCurrentInterruptionLevelsRow
		if err := rows.Scan(
			&i.Region,
			&i.OperatingSystem,
			&i.InstanceType,
			&i.Max,
			&i.InterruptionLevel,
			&i.InterruptionLevelLabel,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getInterruptionChanges = `-- name: GetInterruptionChanges :many
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
WHERE region IN (/*SLICE:regions*/?)
    AND instance_type IN (/*SLICE:instance_types*/?)
    AND operating_system IN (/*SLICE:operating_systems*/?)
ORDER BY observed_time DESC
`

type GetInterruptionChangesParams struct {
	Regions          []string
	InstanceTypes    []string
	OperatingSystems []string
}

type GetInterruptionChangesRow struct {
	ID                         int64
	Region                     string
	InstanceType               string
	OperatingSystem            string
	InterruptionLevel          int64
	InterruptionLevelLabel     string
	LastInterruptionLevel      interface{}
	LastInterruptionLevelLabel interface{}
	ObservedTime               time.Time
}

func (q *Queries) GetInterruptionChanges(ctx context.Context, arg GetInterruptionChangesParams) ([]GetInterruptionChangesRow, error) {
	query := getInterruptionChanges
	var queryParams []interface{}
	if len(arg.Regions) > 0 {
		for _, v := range arg.Regions {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:regions*/?", strings.Repeat(",?", len(arg.Regions))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:regions*/?", "NULL", 1)
	}
	if len(arg.InstanceTypes) > 0 {
		for _, v := range arg.InstanceTypes {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:instance_types*/?", strings.Repeat(",?", len(arg.InstanceTypes))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:instance_types*/?", "NULL", 1)
	}
	if len(arg.OperatingSystems) > 0 {
		for _, v := range arg.OperatingSystems {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:operating_systems*/?", strings.Repeat(",?", len(arg.OperatingSystems))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:operating_systems*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetInterruptionChangesRow
	for rows.Next() {
		var i GetInterruptionChangesRow
		if err := rows.Scan(
			&i.ID,
			&i.Region,
			&i.InstanceType,
			&i.OperatingSystem,
			&i.InterruptionLevel,
			&i.InterruptionLevelLabel,
			&i.LastInterruptionLevel,
			&i.LastInterruptionLevelLabel,
			&i.ObservedTime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertStat = `-- name: InsertStat :exec
INSERT INTO spot_instance_stats (
        region,
        operating_system,
        instance_type,
        interruption_level,
        interruption_level_label,
        observed_time
    )
VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
`

type InsertStatParams struct {
	Region                 string
	OperatingSystem        string
	InstanceType           string
	InterruptionLevel      int64
	InterruptionLevelLabel string
}

func (q *Queries) InsertStat(ctx context.Context, arg InsertStatParams) error {
	_, err := q.db.ExecContext(ctx, insertStat,
		arg.Region,
		arg.OperatingSystem,
		arg.InstanceType,
		arg.InterruptionLevel,
		arg.InterruptionLevelLabel,
	)
	return err
}
