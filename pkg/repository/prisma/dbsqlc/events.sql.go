// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: events.sql

package dbsqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const clearEventPayloadData = `-- name: ClearEventPayloadData :one
WITH for_delete AS (
    SELECT
        e1."id" as "id"
    FROM "Event" e1
    WHERE
        e1."tenantId" = $1::uuid AND
        e1."deletedAt" IS NOT NULL -- TODO change this for all clear queries
        AND e1."data" IS NOT NULL
    LIMIT $2 + 1
    FOR UPDATE SKIP LOCKED
), expired_with_limit AS (
    SELECT
        for_delete."id" as "id"
    FROM for_delete
    LIMIT $2
),
has_more AS (
    SELECT
        CASE
            WHEN COUNT(*) > $2 THEN TRUE
            ELSE FALSE
        END as has_more
    FROM for_delete
)
UPDATE
    "Event"
SET
    "data" = NULL
WHERE
    "id" IN (SELECT "id" FROM expired_with_limit)
RETURNING
    (SELECT has_more FROM has_more) as has_more
`

type ClearEventPayloadDataParams struct {
	Tenantid pgtype.UUID `json:"tenantid"`
	Limit    interface{} `json:"limit"`
}

func (q *Queries) ClearEventPayloadData(ctx context.Context, db DBTX, arg ClearEventPayloadDataParams) (bool, error) {
	row := db.QueryRow(ctx, clearEventPayloadData, arg.Tenantid, arg.Limit)
	var has_more bool
	err := row.Scan(&has_more)
	return has_more, err
}

const countEvents = `-- name: CountEvents :one
WITH events AS (
    SELECT
        events."id"
    FROM
        "Event" as events
    LEFT JOIN
        "WorkflowRunTriggeredBy" as runTriggers ON events."id" = runTriggers."eventId"
    LEFT JOIN
        "WorkflowRun" as runs ON runTriggers."parentId" = runs."id"
    LEFT JOIN
        "WorkflowVersion" as workflowVersion ON workflowVersion."id" = runs."workflowVersionId"
    LEFT JOIN
        "Workflow" as workflow ON workflowVersion."workflowId" = workflow."id"
    WHERE
        events."tenantId" = $1 AND
        events."deletedAt" IS NULL AND
        (
            $2::text[] IS NULL OR
            events."key" = ANY($2::text[])
        ) AND
        (
            $3::jsonb IS NULL OR
            events."additionalMetadata" @> $3::jsonb
        ) AND
        (
            ($4::text[])::uuid[] IS NULL OR
            (workflow."id" = ANY($4::text[]::uuid[]))
        ) AND
        (
            $5::text IS NULL OR
            workflow.name like concat('%', $5::text, '%') OR
            jsonb_path_exists(events."data", cast(concat('$.** ? (@.type() == "string" && @ like_regex "', $5::text, '")') as jsonpath))
        ) AND
        (
            $6::text[] IS NULL OR
            "status" = ANY(cast($6::text[] as "WorkflowRunStatus"[]))
        )
    ORDER BY
        case when $7 = 'createdAt ASC' THEN events."createdAt" END ASC ,
        case when $7 = 'createdAt DESC' then events."createdAt" END DESC
    LIMIT 10000
)
SELECT
    count(events) AS total
FROM
    events
`

type CountEventsParams struct {
	TenantId           pgtype.UUID `json:"tenantId"`
	Keys               []string    `json:"keys"`
	AdditionalMetadata []byte      `json:"additionalMetadata"`
	Workflows          []string    `json:"workflows"`
	Search             pgtype.Text `json:"search"`
	Statuses           []string    `json:"statuses"`
	Orderby            interface{} `json:"orderby"`
}

func (q *Queries) CountEvents(ctx context.Context, db DBTX, arg CountEventsParams) (int64, error) {
	row := db.QueryRow(ctx, countEvents,
		arg.TenantId,
		arg.Keys,
		arg.AdditionalMetadata,
		arg.Workflows,
		arg.Search,
		arg.Statuses,
		arg.Orderby,
	)
	var total int64
	err := row.Scan(&total)
	return total, err
}

const createEvent = `-- name: CreateEvent :one
INSERT INTO "Event" (
    "id",
    "createdAt",
    "updatedAt",
    "deletedAt",
    "key",
    "tenantId",
    "replayedFromId",
    "data",
    "additionalMetadata"
) VALUES (
    $1::uuid,
    coalesce($2::timestamp, CURRENT_TIMESTAMP),
    coalesce($3::timestamp, CURRENT_TIMESTAMP),
    $4::timestamp,
    $5::text,
    $6::uuid,
    $7::uuid,
    $8::jsonb,
    $9::jsonb
) RETURNING id, "createdAt", "updatedAt", "deletedAt", key, "tenantId", "replayedFromId", data, "additionalMetadata"
`

type CreateEventParams struct {
	ID                 pgtype.UUID      `json:"id"`
	CreatedAt          pgtype.Timestamp `json:"createdAt"`
	UpdatedAt          pgtype.Timestamp `json:"updatedAt"`
	Deletedat          pgtype.Timestamp `json:"deletedat"`
	Key                string           `json:"key"`
	Tenantid           pgtype.UUID      `json:"tenantid"`
	ReplayedFromId     pgtype.UUID      `json:"replayedFromId"`
	Data               []byte           `json:"data"`
	Additionalmetadata []byte           `json:"additionalmetadata"`
}

func (q *Queries) CreateEvent(ctx context.Context, db DBTX, arg CreateEventParams) (*Event, error) {
	row := db.QueryRow(ctx, createEvent,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Deletedat,
		arg.Key,
		arg.Tenantid,
		arg.ReplayedFromId,
		arg.Data,
		arg.Additionalmetadata,
	)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Key,
		&i.TenantId,
		&i.ReplayedFromId,
		&i.Data,
		&i.AdditionalMetadata,
	)
	return &i, err
}

const getEventForEngine = `-- name: GetEventForEngine :one
SELECT
    id, "createdAt", "updatedAt", "deletedAt", key, "tenantId", "replayedFromId", data, "additionalMetadata"
FROM
    "Event"
WHERE
    "deletedAt" IS NULL AND
    "id" = $1::uuid
`

func (q *Queries) GetEventForEngine(ctx context.Context, db DBTX, id pgtype.UUID) (*Event, error) {
	row := db.QueryRow(ctx, getEventForEngine, id)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Key,
		&i.TenantId,
		&i.ReplayedFromId,
		&i.Data,
		&i.AdditionalMetadata,
	)
	return &i, err
}

const getEventsForRange = `-- name: GetEventsForRange :many
SELECT
    date_trunc('hour', "createdAt") AS event_hour,
    COUNT(*) AS event_count
FROM
    "Event"
WHERE
    events."deletedAt" IS NOT NULL AND
    "createdAt" >= NOW() - INTERVAL '1 week'
GROUP BY
    event_hour
ORDER BY
    event_hour
`

type GetEventsForRangeRow struct {
	EventHour  pgtype.Interval `json:"event_hour"`
	EventCount int64           `json:"event_count"`
}

func (q *Queries) GetEventsForRange(ctx context.Context, db DBTX) ([]*GetEventsForRangeRow, error) {
	rows, err := db.Query(ctx, getEventsForRange)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetEventsForRangeRow
	for rows.Next() {
		var i GetEventsForRangeRow
		if err := rows.Scan(&i.EventHour, &i.EventCount); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listEvents = `-- name: ListEvents :many
WITH filtered_events AS (
    SELECT
        events."id"
    FROM
        "Event" as events
    LEFT JOIN
        "WorkflowRunTriggeredBy" as runTriggers ON events."id" = runTriggers."eventId"
    LEFT JOIN
        "WorkflowRun" as runs ON runTriggers."parentId" = runs."id"
    LEFT JOIN
        "WorkflowVersion" as workflowVersion ON workflowVersion."id" = runs."workflowVersionId"
    LEFT JOIN
        "Workflow" as workflow ON workflowVersion."workflowId" = workflow."id"
    WHERE
        events."tenantId" = $1 AND
        events."deletedAt" IS NULL AND
        (
            $2::text[] IS NULL OR
            events."key" = ANY($2::text[])
        ) AND
            (
            $3::jsonb IS NULL OR
            events."additionalMetadata" @> $3::jsonb
        ) AND
        (
            ($4::text[])::uuid[] IS NULL OR
            (workflow."id" = ANY($4::text[]::uuid[]))
        ) AND
        (
            $5::text IS NULL OR
            workflow.name like concat('%', $5::text, '%') OR
            jsonb_path_exists(events."data", cast(concat('$.** ? (@.type() == "string" && @ like_regex "', $5::text, '")') as jsonpath))
        ) AND
        (
            $6::text[] IS NULL OR
            "status" = ANY(cast($6::text[] as "WorkflowRunStatus"[]))
        )
    ORDER BY
        case when $7 = 'createdAt ASC' THEN events."createdAt" END ASC ,
        case when $7 = 'createdAt DESC' then events."createdAt" END DESC
    OFFSET
        COALESCE($8, 0)
    LIMIT
        COALESCE($9, 50)
)
SELECT
    events.id, events."createdAt", events."updatedAt", events."deletedAt", events.key, events."tenantId", events."replayedFromId", events.data, events."additionalMetadata",
    sum(case when runs."status" = 'PENDING' then 1 else 0 end) AS pendingRuns,
    sum(case when runs."status" = 'QUEUED' then 1 else 0 end) AS queuedRuns,
    sum(case when runs."status" = 'RUNNING' then 1 else 0 end) AS runningRuns,
    sum(case when runs."status" = 'SUCCEEDED' then 1 else 0 end) AS succeededRuns,
    sum(case when runs."status" = 'FAILED' then 1 else 0 end) AS failedRuns
FROM
    filtered_events
JOIN
    "Event" as events ON events."id" = filtered_events."id"
LEFT JOIN
    "WorkflowRunTriggeredBy" as runTriggers ON events."id" = runTriggers."eventId"
LEFT JOIN
    "WorkflowRun" as runs ON runTriggers."parentId" = runs."id"
GROUP BY
    events."id", events."createdAt"
`

type ListEventsParams struct {
	TenantId           pgtype.UUID `json:"tenantId"`
	Keys               []string    `json:"keys"`
	AdditionalMetadata []byte      `json:"additionalMetadata"`
	Workflows          []string    `json:"workflows"`
	Search             pgtype.Text `json:"search"`
	Statuses           []string    `json:"statuses"`
	Orderby            interface{} `json:"orderby"`
	Offset             interface{} `json:"offset"`
	Limit              interface{} `json:"limit"`
}

type ListEventsRow struct {
	Event         Event `json:"event"`
	Pendingruns   int64 `json:"pendingruns"`
	Queuedruns    int64 `json:"queuedruns"`
	Runningruns   int64 `json:"runningruns"`
	Succeededruns int64 `json:"succeededruns"`
	Failedruns    int64 `json:"failedruns"`
}

func (q *Queries) ListEvents(ctx context.Context, db DBTX, arg ListEventsParams) ([]*ListEventsRow, error) {
	rows, err := db.Query(ctx, listEvents,
		arg.TenantId,
		arg.Keys,
		arg.AdditionalMetadata,
		arg.Workflows,
		arg.Search,
		arg.Statuses,
		arg.Orderby,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListEventsRow
	for rows.Next() {
		var i ListEventsRow
		if err := rows.Scan(
			&i.Event.ID,
			&i.Event.CreatedAt,
			&i.Event.UpdatedAt,
			&i.Event.DeletedAt,
			&i.Event.Key,
			&i.Event.TenantId,
			&i.Event.ReplayedFromId,
			&i.Event.Data,
			&i.Event.AdditionalMetadata,
			&i.Pendingruns,
			&i.Queuedruns,
			&i.Runningruns,
			&i.Succeededruns,
			&i.Failedruns,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listEventsByIDs = `-- name: ListEventsByIDs :many
SELECT
    id, "createdAt", "updatedAt", "deletedAt", key, "tenantId", "replayedFromId", data, "additionalMetadata"
FROM
    "Event" as events
WHERE
    events."deletedAt" IS NULL AND
    "tenantId" = $1::uuid AND
    "id" = ANY ($2::uuid[])
`

type ListEventsByIDsParams struct {
	Tenantid pgtype.UUID   `json:"tenantid"`
	Ids      []pgtype.UUID `json:"ids"`
}

func (q *Queries) ListEventsByIDs(ctx context.Context, db DBTX, arg ListEventsByIDsParams) ([]*Event, error) {
	rows, err := db.Query(ctx, listEventsByIDs, arg.Tenantid, arg.Ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.Key,
			&i.TenantId,
			&i.ReplayedFromId,
			&i.Data,
			&i.AdditionalMetadata,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const softDeleteExpiredEvents = `-- name: SoftDeleteExpiredEvents :one
WITH for_delete AS (
    SELECT
        "id"
    FROM "Event" e
    WHERE
        e."tenantId" = $1::uuid AND
        e."createdAt" < $2::timestamp AND
        e."deletedAt" IS NULL
    ORDER BY e."createdAt" ASC
    LIMIT $3 +1
    FOR UPDATE SKIP LOCKED
),expired_with_limit AS (
    SELECT
        for_delete."id" as "id"
    FROM for_delete
    LIMIT $3
), has_more AS (
    SELECT
        CASE
            WHEN COUNT(*) > $3 THEN TRUE
            ELSE FALSE
        END as has_more
    FROM for_delete
)
UPDATE
    "Event"
SET
    "deletedAt" = CURRENT_TIMESTAMP
WHERE
    "id" IN (SELECT "id" FROM expired_with_limit)
RETURNING
    (SELECT has_more FROM has_more) as has_more
`

type SoftDeleteExpiredEventsParams struct {
	Tenantid      pgtype.UUID      `json:"tenantid"`
	Createdbefore pgtype.Timestamp `json:"createdbefore"`
	Limit         interface{}      `json:"limit"`
}

func (q *Queries) SoftDeleteExpiredEvents(ctx context.Context, db DBTX, arg SoftDeleteExpiredEventsParams) (bool, error) {
	row := db.QueryRow(ctx, softDeleteExpiredEvents, arg.Tenantid, arg.Createdbefore, arg.Limit)
	var has_more bool
	err := row.Scan(&has_more)
	return has_more, err
}
