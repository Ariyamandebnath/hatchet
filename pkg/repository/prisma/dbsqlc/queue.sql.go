// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: queue.sql

package dbsqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const bulkQueueItems = `-- name: BulkQueueItems :exec
UPDATE
    "QueueItem" qi
SET
    "isQueued" = false
WHERE
    qi."id" = ANY($1::bigint[])
`

func (q *Queries) BulkQueueItems(ctx context.Context, db DBTX, ids []int64) error {
	_, err := db.Exec(ctx, bulkQueueItems, ids)
	return err
}

const createQueueItem = `-- name: CreateQueueItem :exec
INSERT INTO
    "QueueItem" (
        "stepRunId",
        "stepId",
        "actionId",
        "scheduleTimeoutAt",
        "stepTimeout",
        "priority",
        "isQueued",
        "tenantId",
        "queue",
        "sticky",
        "desiredWorkerId"
    )
VALUES
    (
        $1::uuid,
        $2::uuid,
        $3::text,
        $4::timestamp,
        $5::text,
        COALESCE($6::integer, 1),
        true,
        $7::uuid,
        $8,
        $9::"StickyStrategy",
        $10::uuid
    )
`

type CreateQueueItemParams struct {
	StepRunId         pgtype.UUID        `json:"stepRunId"`
	StepId            pgtype.UUID        `json:"stepId"`
	ActionId          pgtype.Text        `json:"actionId"`
	ScheduleTimeoutAt pgtype.Timestamp   `json:"scheduleTimeoutAt"`
	StepTimeout       pgtype.Text        `json:"stepTimeout"`
	Priority          pgtype.Int4        `json:"priority"`
	Tenantid          pgtype.UUID        `json:"tenantid"`
	Queue             string             `json:"queue"`
	Sticky            NullStickyStrategy `json:"sticky"`
	DesiredWorkerId   pgtype.UUID        `json:"desiredWorkerId"`
}

func (q *Queries) CreateQueueItem(ctx context.Context, db DBTX, arg CreateQueueItemParams) error {
	_, err := db.Exec(ctx, createQueueItem,
		arg.StepRunId,
		arg.StepId,
		arg.ActionId,
		arg.ScheduleTimeoutAt,
		arg.StepTimeout,
		arg.Priority,
		arg.Tenantid,
		arg.Queue,
		arg.Sticky,
		arg.DesiredWorkerId,
	)
	return err
}

const listQueues = `-- name: ListQueues :many
SELECT
    id, "tenantId", name
FROM
    "Queue"
WHERE
    "tenantId" = $1::uuid
`

func (q *Queries) ListQueues(ctx context.Context, db DBTX, tenantid pgtype.UUID) ([]*Queue, error) {
	rows, err := db.Query(ctx, listQueues, tenantid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Queue
	for rows.Next() {
		var i Queue
		if err := rows.Scan(&i.ID, &i.TenantId, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const upsertQueue = `-- name: UpsertQueue :exec
INSERT INTO
    "Queue" (
        "tenantId",
        "name"
    )
VALUES
    (
        $1::uuid,
        $2::text
    )
ON CONFLICT ("tenantId", "name") DO NOTHING
`

type UpsertQueueParams struct {
	Tenantid pgtype.UUID `json:"tenantid"`
	Name     string      `json:"name"`
}

func (q *Queries) UpsertQueue(ctx context.Context, db DBTX, arg UpsertQueueParams) error {
	_, err := db.Exec(ctx, upsertQueue, arg.Tenantid, arg.Name)
	return err
}
