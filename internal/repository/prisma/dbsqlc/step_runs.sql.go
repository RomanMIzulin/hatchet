// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: step_runs.sql

package dbsqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const archiveStepRunResultFromStepRun = `-- name: ArchiveStepRunResultFromStepRun :one
WITH step_run_data AS (
    SELECT
        "id" AS step_run_id,
        "createdAt",
        "updatedAt",
        "deletedAt",
        "order",
        "input",
        "output",
        "error",
        "startedAt",
        "finishedAt",
        "timeoutAt",
        "cancelledAt",
        "cancelledReason",
        "cancelledError"
    FROM "StepRun"
    WHERE "id" = $2::uuid AND "tenantId" = $3::uuid
)
INSERT INTO "StepRunResultArchive" (
    "id",
    "createdAt",
    "updatedAt",
    "deletedAt",
    "stepRunId",
    "input",
    "output",
    "error",
    "startedAt",
    "finishedAt",
    "timeoutAt",
    "cancelledAt",
    "cancelledReason",
    "cancelledError"
)
SELECT
    COALESCE($1::uuid, gen_random_uuid()),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    step_run_data."deletedAt",
    step_run_data.step_run_id,
    step_run_data."input",
    step_run_data."output",
    step_run_data."error",
    step_run_data."startedAt",
    step_run_data."finishedAt",
    step_run_data."timeoutAt",
    step_run_data."cancelledAt",
    step_run_data."cancelledReason",
    step_run_data."cancelledError"
FROM step_run_data
RETURNING id, "createdAt", "updatedAt", "deletedAt", "stepRunId", "order", input, output, error, "startedAt", "finishedAt", "timeoutAt", "cancelledAt", "cancelledReason", "cancelledError"
`

type ArchiveStepRunResultFromStepRunParams struct {
	ID        pgtype.UUID `json:"id"`
	Steprunid pgtype.UUID `json:"steprunid"`
	Tenantid  pgtype.UUID `json:"tenantid"`
}

func (q *Queries) ArchiveStepRunResultFromStepRun(ctx context.Context, db DBTX, arg ArchiveStepRunResultFromStepRunParams) (*StepRunResultArchive, error) {
	row := db.QueryRow(ctx, archiveStepRunResultFromStepRun, arg.ID, arg.Steprunid, arg.Tenantid)
	var i StepRunResultArchive
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.StepRunId,
		&i.Order,
		&i.Input,
		&i.Output,
		&i.Error,
		&i.StartedAt,
		&i.FinishedAt,
		&i.TimeoutAt,
		&i.CancelledAt,
		&i.CancelledReason,
		&i.CancelledError,
	)
	return &i, err
}

const assignStepRunToWorker = `-- name: AssignStepRunToWorker :one
WITH valid_workers AS (
    SELECT
        w."id", w."dispatcherId", COALESCE(ws."slots", 100) AS "slots"
    FROM
        "Worker" w
    LEFT JOIN 
        "WorkerSemaphore" ws ON w."id" = ws."workerId"
    WHERE
        w."tenantId" = $1::uuid
        AND w."dispatcherId" IS NOT NULL
        AND w."lastHeartbeatAt" > NOW() - INTERVAL '5 seconds'
        AND w."id" IN (
            SELECT "_ActionToWorker"."B"
            FROM "_ActionToWorker"
            INNER JOIN "Action" ON "Action"."id" = "_ActionToWorker"."A"
            WHERE "Action"."tenantId" = $1 AND "Action"."actionId" = $2::text
        )
        AND (
            ws."workerId" IS NULL OR
            ws."slots" > 0
        )
    ORDER BY ws."slots" DESC NULLS FIRST, RANDOM()
), total_slots AS (
    SELECT
        COALESCE(SUM(vw."slots"), 0) AS "totalSlots"
    FROM
        valid_workers vw
), selected_worker AS (
    SELECT
        id, "dispatcherId", slots
    FROM
        valid_workers vw
    LIMIT 1
),
step_run AS (
    SELECT
        "id", "workerId"
    FROM
        "StepRun"
    WHERE
        "id" = $3::uuid AND
        "tenantId" = $1::uuid AND
        "status" = 'PENDING_ASSIGNMENT' AND
        EXISTS (SELECT 1 FROM selected_worker)
    FOR UPDATE
),
update_step_run AS (
    UPDATE
        "StepRun"
    SET
        "status" = 'ASSIGNED',
        "workerId" = (
            SELECT "id"
            FROM selected_worker
            LIMIT 1
        ),
        "tickerId" = NULL,
        "updatedAt" = CURRENT_TIMESTAMP,
        "timeoutAt" = CASE
            WHEN $4::text IS NOT NULL THEN
                CURRENT_TIMESTAMP + convert_duration_to_interval($4::text)
            ELSE CURRENT_TIMESTAMP + INTERVAL '5 minutes'
        END
    WHERE
        "id" = $3::uuid AND
        "tenantId" = $1::uuid AND
        "status" = 'PENDING_ASSIGNMENT' AND
        EXISTS (SELECT 1 FROM selected_worker)
    RETURNING 
        "StepRun"."id", "StepRun"."workerId", 
        (SELECT "dispatcherId" FROM selected_worker) AS "dispatcherId"
)
SELECT ts."totalSlots"::int, usr."id", usr."workerId", usr."dispatcherId"
FROM total_slots ts
LEFT JOIN update_step_run usr ON true
`

type AssignStepRunToWorkerParams struct {
	Tenantid    pgtype.UUID `json:"tenantid"`
	Actionid    string      `json:"actionid"`
	Steprunid   pgtype.UUID `json:"steprunid"`
	StepTimeout pgtype.Text `json:"stepTimeout"`
}

type AssignStepRunToWorkerRow struct {
	TsTotalSlots int32       `json:"ts_totalSlots"`
	ID           pgtype.UUID `json:"id"`
	WorkerId     pgtype.UUID `json:"workerId"`
	DispatcherId pgtype.UUID `json:"dispatcherId"`
}

func (q *Queries) AssignStepRunToWorker(ctx context.Context, db DBTX, arg AssignStepRunToWorkerParams) (*AssignStepRunToWorkerRow, error) {
	row := db.QueryRow(ctx, assignStepRunToWorker,
		arg.Tenantid,
		arg.Actionid,
		arg.Steprunid,
		arg.StepTimeout,
	)
	var i AssignStepRunToWorkerRow
	err := row.Scan(
		&i.TsTotalSlots,
		&i.ID,
		&i.WorkerId,
		&i.DispatcherId,
	)
	return &i, err
}

const countStepRunEvents = `-- name: CountStepRunEvents :one
SELECT
    count(*) OVER() AS total
FROM
    "StepRunEvent"
WHERE
    "stepRunId" = $1::uuid
`

func (q *Queries) CountStepRunEvents(ctx context.Context, db DBTX, steprunid pgtype.UUID) (int64, error) {
	row := db.QueryRow(ctx, countStepRunEvents, steprunid)
	var total int64
	err := row.Scan(&total)
	return total, err
}

const createStepRunEvent = `-- name: CreateStepRunEvent :exec
WITH input_values AS (
    SELECT
        CURRENT_TIMESTAMP AS "timeFirstSeen",
        CURRENT_TIMESTAMP AS "timeLastSeen",
        $1::uuid AS "stepRunId",
        $2::"StepRunEventReason" AS "reason",
        $3::"StepRunEventSeverity" AS "severity",
        $4::text AS "message",
        1 AS "count",
        $5::jsonb AS "data"
),
updated AS (
    UPDATE "StepRunEvent"
    SET
        "timeLastSeen" = CURRENT_TIMESTAMP,
        "message" = input_values."message",
        "count" = "StepRunEvent"."count" + 1,
        "data" = input_values."data"
    FROM input_values
    WHERE
        "StepRunEvent"."stepRunId" = input_values."stepRunId"
        AND "StepRunEvent"."reason" = input_values."reason"
        AND "StepRunEvent"."severity" = input_values."severity"
        AND "StepRunEvent"."id" = (
            SELECT "id"
            FROM "StepRunEvent"
            WHERE "stepRunId" = input_values."stepRunId"
            ORDER BY "id" DESC
            LIMIT 1
        )
    RETURNING "StepRunEvent".id, "StepRunEvent"."timeFirstSeen", "StepRunEvent"."timeLastSeen", "StepRunEvent"."stepRunId", "StepRunEvent".reason, "StepRunEvent".severity, "StepRunEvent".message, "StepRunEvent".count, "StepRunEvent".data
)
INSERT INTO "StepRunEvent" (
    "timeFirstSeen",
    "timeLastSeen",
    "stepRunId",
    "reason",
    "severity",
    "message",
    "count",
    "data"
)
SELECT
    "timeFirstSeen",
    "timeLastSeen",
    "stepRunId",
    "reason",
    "severity",
    "message",
    "count",
    "data"
FROM input_values
WHERE NOT EXISTS (
    SELECT 1 FROM updated WHERE "stepRunId" = input_values."stepRunId"
)
`

type CreateStepRunEventParams struct {
	Steprunid pgtype.UUID          `json:"steprunid"`
	Reason    StepRunEventReason   `json:"reason"`
	Severity  StepRunEventSeverity `json:"severity"`
	Message   string               `json:"message"`
	Data      []byte               `json:"data"`
}

func (q *Queries) CreateStepRunEvent(ctx context.Context, db DBTX, arg CreateStepRunEventParams) error {
	_, err := db.Exec(ctx, createStepRunEvent,
		arg.Steprunid,
		arg.Reason,
		arg.Severity,
		arg.Message,
		arg.Data,
	)
	return err
}

const getStepRun = `-- name: GetStepRun :one
SELECT
    "StepRun".id, "StepRun"."createdAt", "StepRun"."updatedAt", "StepRun"."deletedAt", "StepRun"."tenantId", "StepRun"."jobRunId", "StepRun"."stepId", "StepRun"."order", "StepRun"."workerId", "StepRun"."tickerId", "StepRun".status, "StepRun".input, "StepRun".output, "StepRun"."requeueAfter", "StepRun"."scheduleTimeoutAt", "StepRun".error, "StepRun"."startedAt", "StepRun"."finishedAt", "StepRun"."timeoutAt", "StepRun"."cancelledAt", "StepRun"."cancelledReason", "StepRun"."cancelledError", "StepRun"."inputSchema", "StepRun"."callerFiles", "StepRun"."gitRepoBranch", "StepRun"."retryCount"
FROM
    "StepRun"
WHERE
    "id" = $1::uuid AND
    "tenantId" = $2::uuid
`

type GetStepRunParams struct {
	ID       pgtype.UUID `json:"id"`
	Tenantid pgtype.UUID `json:"tenantid"`
}

func (q *Queries) GetStepRun(ctx context.Context, db DBTX, arg GetStepRunParams) (*StepRun, error) {
	row := db.QueryRow(ctx, getStepRun, arg.ID, arg.Tenantid)
	var i StepRun
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.JobRunId,
		&i.StepId,
		&i.Order,
		&i.WorkerId,
		&i.TickerId,
		&i.Status,
		&i.Input,
		&i.Output,
		&i.RequeueAfter,
		&i.ScheduleTimeoutAt,
		&i.Error,
		&i.StartedAt,
		&i.FinishedAt,
		&i.TimeoutAt,
		&i.CancelledAt,
		&i.CancelledReason,
		&i.CancelledError,
		&i.InputSchema,
		&i.CallerFiles,
		&i.GitRepoBranch,
		&i.RetryCount,
	)
	return &i, err
}

const getStepRunForEngine = `-- name: GetStepRunForEngine :many
SELECT
    DISTINCT ON (sr."id")
    sr.id, sr."createdAt", sr."updatedAt", sr."deletedAt", sr."tenantId", sr."jobRunId", sr."stepId", sr."order", sr."workerId", sr."tickerId", sr.status, sr.input, sr.output, sr."requeueAfter", sr."scheduleTimeoutAt", sr.error, sr."startedAt", sr."finishedAt", sr."timeoutAt", sr."cancelledAt", sr."cancelledReason", sr."cancelledError", sr."inputSchema", sr."callerFiles", sr."gitRepoBranch", sr."retryCount",
    jrld."data" AS "jobRunLookupData",
    -- TODO: everything below this line is cacheable and should be moved to a separate query
    jr."id" AS "jobRunId",
    wr."id" AS "workflowRunId",
    s."id" AS "stepId",
    s."retries" AS "stepRetries",
    s."timeout" AS "stepTimeout",
    s."scheduleTimeout" AS "stepScheduleTimeout",
    s."readableId" AS "stepReadableId",
    s."customUserData" AS "stepCustomUserData",
    j."name" AS "jobName",
    j."id" AS "jobId",
    j."kind" AS "jobKind",
    wv."id" AS "workflowVersionId",
    w."name" AS "workflowName",
    w."id" AS "workflowId",
    a."actionId" AS "actionId"
FROM
    "StepRun" sr
JOIN
    "Step" s ON sr."stepId" = s."id"
JOIN
    "Action" a ON s."actionId" = a."actionId" AND s."tenantId" = a."tenantId"
JOIN
    "JobRun" jr ON sr."jobRunId" = jr."id"
JOIN
    "JobRunLookupData" jrld ON jr."id" = jrld."jobRunId"
JOIN
    "Job" j ON jr."jobId" = j."id"
JOIN 
    "WorkflowRun" wr ON jr."workflowRunId" = wr."id"
JOIN
    "WorkflowVersion" wv ON wr."workflowVersionId" = wv."id"
JOIN
    "Workflow" w ON wv."workflowId" = w."id"
WHERE
    sr."id" = ANY($1::uuid[]) AND
    (
        $2::uuid IS NULL OR
        sr."tenantId" = $2::uuid
    )
`

type GetStepRunForEngineParams struct {
	Ids      []pgtype.UUID `json:"ids"`
	TenantId pgtype.UUID   `json:"tenantId"`
}

type GetStepRunForEngineRow struct {
	StepRun             StepRun     `json:"step_run"`
	JobRunLookupData    []byte      `json:"jobRunLookupData"`
	JobRunId            pgtype.UUID `json:"jobRunId"`
	WorkflowRunId       pgtype.UUID `json:"workflowRunId"`
	StepId              pgtype.UUID `json:"stepId"`
	StepRetries         int32       `json:"stepRetries"`
	StepTimeout         pgtype.Text `json:"stepTimeout"`
	StepScheduleTimeout string      `json:"stepScheduleTimeout"`
	StepReadableId      pgtype.Text `json:"stepReadableId"`
	StepCustomUserData  []byte      `json:"stepCustomUserData"`
	JobName             string      `json:"jobName"`
	JobId               pgtype.UUID `json:"jobId"`
	JobKind             JobKind     `json:"jobKind"`
	WorkflowVersionId   pgtype.UUID `json:"workflowVersionId"`
	WorkflowName        string      `json:"workflowName"`
	WorkflowId          pgtype.UUID `json:"workflowId"`
	ActionId            string      `json:"actionId"`
}

func (q *Queries) GetStepRunForEngine(ctx context.Context, db DBTX, arg GetStepRunForEngineParams) ([]*GetStepRunForEngineRow, error) {
	rows, err := db.Query(ctx, getStepRunForEngine, arg.Ids, arg.TenantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetStepRunForEngineRow
	for rows.Next() {
		var i GetStepRunForEngineRow
		if err := rows.Scan(
			&i.StepRun.ID,
			&i.StepRun.CreatedAt,
			&i.StepRun.UpdatedAt,
			&i.StepRun.DeletedAt,
			&i.StepRun.TenantId,
			&i.StepRun.JobRunId,
			&i.StepRun.StepId,
			&i.StepRun.Order,
			&i.StepRun.WorkerId,
			&i.StepRun.TickerId,
			&i.StepRun.Status,
			&i.StepRun.Input,
			&i.StepRun.Output,
			&i.StepRun.RequeueAfter,
			&i.StepRun.ScheduleTimeoutAt,
			&i.StepRun.Error,
			&i.StepRun.StartedAt,
			&i.StepRun.FinishedAt,
			&i.StepRun.TimeoutAt,
			&i.StepRun.CancelledAt,
			&i.StepRun.CancelledReason,
			&i.StepRun.CancelledError,
			&i.StepRun.InputSchema,
			&i.StepRun.CallerFiles,
			&i.StepRun.GitRepoBranch,
			&i.StepRun.RetryCount,
			&i.JobRunLookupData,
			&i.JobRunId,
			&i.WorkflowRunId,
			&i.StepId,
			&i.StepRetries,
			&i.StepTimeout,
			&i.StepScheduleTimeout,
			&i.StepReadableId,
			&i.StepCustomUserData,
			&i.JobName,
			&i.JobId,
			&i.JobKind,
			&i.WorkflowVersionId,
			&i.WorkflowName,
			&i.WorkflowId,
			&i.ActionId,
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

const getTotalSlots = `-- name: GetTotalSlots :many
WITH valid_workers AS (
    SELECT
        w."id", w."dispatcherId", COALESCE(ws."slots", 100) AS "slots"
    FROM
        "Worker" w
    LEFT JOIN 
        "WorkerSemaphore" ws ON w."id" = ws."workerId"
    WHERE
        w."tenantId" = $1::uuid
        AND w."dispatcherId" IS NOT NULL
        AND w."lastHeartbeatAt" > NOW() - INTERVAL '5 seconds'
        AND w."id" IN (
            SELECT "_ActionToWorker"."B"
            FROM "_ActionToWorker"
            INNER JOIN "Action" ON "Action"."id" = "_ActionToWorker"."A"
            WHERE "Action"."tenantId" = $1 AND "Action"."actionId" = $2::text
        )
        AND (
            ws."workerId" IS NULL OR
            ws."slots" > 0
        )
    ORDER BY ws."slots" DESC NULLS FIRST, RANDOM()
),
total_slots AS (
    SELECT
        COALESCE(SUM(vw."slots"), 0) AS "totalSlots"
    FROM
        valid_workers vw
)
SELECT ts."totalSlots"::int, vw."id", vw."slots"
FROM valid_workers vw
LEFT JOIN total_slots ts ON true
`

type GetTotalSlotsParams struct {
	Tenantid pgtype.UUID `json:"tenantid"`
	Actionid string      `json:"actionid"`
}

type GetTotalSlotsRow struct {
	TsTotalSlots int32       `json:"ts_totalSlots"`
	ID           pgtype.UUID `json:"id"`
	Slots        int32       `json:"slots"`
}

func (q *Queries) GetTotalSlots(ctx context.Context, db DBTX, arg GetTotalSlotsParams) ([]*GetTotalSlotsRow, error) {
	rows, err := db.Query(ctx, getTotalSlots, arg.Tenantid, arg.Actionid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetTotalSlotsRow
	for rows.Next() {
		var i GetTotalSlotsRow
		if err := rows.Scan(&i.TsTotalSlots, &i.ID, &i.Slots); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listStartableStepRuns = `-- name: ListStartableStepRuns :many
WITH job_run AS (
    SELECT "status"
    FROM "JobRun"
    WHERE "id" = $1::uuid
)
SELECT 
    DISTINCT ON (child_run."id")
    child_run."id" AS "id"
FROM 
    "StepRun" AS child_run
LEFT JOIN 
    "_StepRunOrder" AS step_run_order ON step_run_order."B" = child_run."id"
JOIN
    job_run ON true
WHERE 
    child_run."jobRunId" = $1::uuid
    AND child_run."status" = 'PENDING'
    AND job_run."status" = 'RUNNING'
    -- case on whether parentStepRunId is null
    AND (
        ($2::uuid IS NULL AND step_run_order."A" IS NULL) OR 
        (
            step_run_order."A" = $2::uuid
            AND NOT EXISTS (
                SELECT 1
                FROM "_StepRunOrder" AS parent_order
                JOIN "StepRun" AS parent_run ON parent_order."A" = parent_run."id"
                WHERE 
                    parent_order."B" = child_run."id"
                    AND parent_run."status" != 'SUCCEEDED'
            )
        )
    )
`

type ListStartableStepRunsParams struct {
	Jobrunid        pgtype.UUID `json:"jobrunid"`
	ParentStepRunId pgtype.UUID `json:"parentStepRunId"`
}

func (q *Queries) ListStartableStepRuns(ctx context.Context, db DBTX, arg ListStartableStepRunsParams) ([]pgtype.UUID, error) {
	rows, err := db.Query(ctx, listStartableStepRuns, arg.Jobrunid, arg.ParentStepRunId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []pgtype.UUID
	for rows.Next() {
		var id pgtype.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listStepRunEvents = `-- name: ListStepRunEvents :many
SELECT
    id, "timeFirstSeen", "timeLastSeen", "stepRunId", reason, severity, message, count, data
FROM
    "StepRunEvent"
WHERE
    "stepRunId" = $1::uuid
ORDER BY
    "id" DESC
OFFSET
    COALESCE($2, 0)
LIMIT
    COALESCE($3, 50)
`

type ListStepRunEventsParams struct {
	Steprunid pgtype.UUID `json:"steprunid"`
	Offset    interface{} `json:"offset"`
	Limit     interface{} `json:"limit"`
}

func (q *Queries) ListStepRunEvents(ctx context.Context, db DBTX, arg ListStepRunEventsParams) ([]*StepRunEvent, error) {
	rows, err := db.Query(ctx, listStepRunEvents, arg.Steprunid, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*StepRunEvent
	for rows.Next() {
		var i StepRunEvent
		if err := rows.Scan(
			&i.ID,
			&i.TimeFirstSeen,
			&i.TimeLastSeen,
			&i.StepRunId,
			&i.Reason,
			&i.Severity,
			&i.Message,
			&i.Count,
			&i.Data,
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

const listStepRuns = `-- name: ListStepRuns :many
SELECT
    DISTINCT ON ("StepRun"."id")
    "StepRun"."id"
FROM
    "StepRun"
JOIN
    "JobRun" ON "StepRun"."jobRunId" = "JobRun"."id"
WHERE
    (
        $1::uuid IS NULL OR
        "StepRun"."tenantId" = $1::uuid
    )
    AND (
        $2::"StepRunStatus" IS NULL OR
        "StepRun"."status" = $2::"StepRunStatus"
    )
    AND (
        $3::uuid[] IS NULL OR
        "JobRun"."workflowRunId" = ANY($3::uuid[])
    )
    AND (
        $4::uuid IS NULL OR
        "StepRun"."jobRunId" = $4::uuid
    )
    AND (
        $5::uuid IS NULL OR
        "StepRun"."tickerId" = $5::uuid
    )
`

type ListStepRunsParams struct {
	TenantId       pgtype.UUID       `json:"tenantId"`
	Status         NullStepRunStatus `json:"status"`
	WorkflowRunIds []pgtype.UUID     `json:"workflowRunIds"`
	JobRunId       pgtype.UUID       `json:"jobRunId"`
	TickerId       pgtype.UUID       `json:"tickerId"`
}

func (q *Queries) ListStepRuns(ctx context.Context, db DBTX, arg ListStepRunsParams) ([]pgtype.UUID, error) {
	rows, err := db.Query(ctx, listStepRuns,
		arg.TenantId,
		arg.Status,
		arg.WorkflowRunIds,
		arg.JobRunId,
		arg.TickerId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []pgtype.UUID
	for rows.Next() {
		var id pgtype.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listStepRunsToReassign = `-- name: ListStepRunsToReassign :many
WITH valid_workers AS (
    SELECT
        DISTINCT ON (w."id")
        w."id",
        COALESCE(SUM(ws."slots"), 100) AS "slots"
    FROM
        "Worker" w
    LEFT JOIN
        "WorkerSemaphore" ws ON ws."workerId" = w."id"
    WHERE
        w."tenantId" = $1::uuid
        AND w."lastHeartbeatAt" > NOW() - INTERVAL '5 seconds'
    GROUP BY
        w."id"
),
total_max_runs AS (
    SELECT
        SUM("slots") AS "totalMaxRuns"
    FROM
        valid_workers
),
limit_max_runs AS (
    SELECT
        GREATEST("totalMaxRuns", 100) AS "limitMaxRuns"
    FROM
        total_max_runs
),
step_runs AS (
    SELECT
        sr.id, sr."createdAt", sr."updatedAt", sr."deletedAt", sr."tenantId", sr."jobRunId", sr."stepId", sr."order", sr."workerId", sr."tickerId", sr.status, sr.input, sr.output, sr."requeueAfter", sr."scheduleTimeoutAt", sr.error, sr."startedAt", sr."finishedAt", sr."timeoutAt", sr."cancelledAt", sr."cancelledReason", sr."cancelledError", sr."inputSchema", sr."callerFiles", sr."gitRepoBranch", sr."retryCount"
    FROM
        "StepRun" sr
    LEFT JOIN
        "Worker" w ON sr."workerId" = w."id"
    JOIN
        "JobRun" jr ON sr."jobRunId" = jr."id"
    JOIN
        "Step" s ON sr."stepId" = s."id"
    WHERE
        sr."tenantId" = $1::uuid
        AND ((
            sr."status" = 'RUNNING'
            AND w."lastHeartbeatAt" < NOW() - INTERVAL '30 seconds'
            AND s."retries" > sr."retryCount"
        ) OR (
            sr."status" = 'ASSIGNED'
            AND w."lastHeartbeatAt" < NOW() - INTERVAL '30 seconds'
        ))
        AND jr."status" = 'RUNNING'
        AND sr."input" IS NOT NULL
        -- Step run cannot have a failed parent
        AND NOT EXISTS (
            SELECT 1
            FROM "_StepRunOrder" AS order_table
            JOIN "StepRun" AS prev_sr ON order_table."A" = prev_sr."id"
            WHERE 
                order_table."B" = sr."id"
                AND prev_sr."status" != 'SUCCEEDED'
        )
    ORDER BY
        sr."createdAt" ASC
    LIMIT
        (SELECT "limitMaxRuns" FROM limit_max_runs)
),
locked_step_runs AS (
    SELECT
        sr."id", sr."status", sr."workerId"
    FROM
        step_runs sr
    FOR UPDATE SKIP LOCKED
)
UPDATE
    "StepRun"
SET
    "status" = 'PENDING_ASSIGNMENT',
    -- requeue after now plus 4 seconds
    "requeueAfter" = CURRENT_TIMESTAMP + INTERVAL '4 seconds',
    "updatedAt" = CURRENT_TIMESTAMP
FROM 
    locked_step_runs
WHERE
    "StepRun"."id" = locked_step_runs."id"
RETURNING "StepRun"."id"
`

// Count the total number of slots across all workers
func (q *Queries) ListStepRunsToReassign(ctx context.Context, db DBTX, tenantid pgtype.UUID) ([]pgtype.UUID, error) {
	rows, err := db.Query(ctx, listStepRunsToReassign, tenantid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []pgtype.UUID
	for rows.Next() {
		var id pgtype.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listStepRunsToRequeue = `-- name: ListStepRunsToRequeue :many
WITH valid_workers AS (
    SELECT
        DISTINCT ON (w."id")
        w."id",
        COALESCE(SUM(ws."slots"), 100) AS "slots"
    FROM
        "Worker" w
    LEFT JOIN
        "WorkerSemaphore" ws ON w."id" = ws."workerId"
    WHERE
        w."tenantId" = $1::uuid
        AND w."lastHeartbeatAt" > NOW() - INTERVAL '5 seconds'
    GROUP BY
        w."id"
),
total_max_runs AS (
    SELECT
        -- if maxRuns is null, then we assume the worker can run 100 step runs
        SUM("slots") AS "totalMaxRuns"
    FROM
        valid_workers
),
step_runs AS (
    SELECT
        sr.id, sr."createdAt", sr."updatedAt", sr."deletedAt", sr."tenantId", sr."jobRunId", sr."stepId", sr."order", sr."workerId", sr."tickerId", sr.status, sr.input, sr.output, sr."requeueAfter", sr."scheduleTimeoutAt", sr.error, sr."startedAt", sr."finishedAt", sr."timeoutAt", sr."cancelledAt", sr."cancelledReason", sr."cancelledError", sr."inputSchema", sr."callerFiles", sr."gitRepoBranch", sr."retryCount"
    FROM
        "StepRun" sr
    LEFT JOIN
        "Worker" w ON sr."workerId" = w."id"
    JOIN
        "JobRun" jr ON sr."jobRunId" = jr."id"
    WHERE
        sr."tenantId" = $1::uuid
        AND sr."requeueAfter" < NOW()
        AND (sr."status" = 'PENDING' OR sr."status" = 'PENDING_ASSIGNMENT')
        AND jr."status" = 'RUNNING'
        AND sr."input" IS NOT NULL
        AND NOT EXISTS (
            SELECT 1
            FROM "_StepRunOrder" AS order_table
            JOIN "StepRun" AS prev_sr ON order_table."A" = prev_sr."id"
            WHERE 
                order_table."B" = sr."id"
                AND prev_sr."status" != 'SUCCEEDED'
        )
    ORDER BY
        sr."createdAt" ASC
    LIMIT
        COALESCE((SELECT "totalMaxRuns" FROM total_max_runs), 100)
),
locked_step_runs AS (
    SELECT
        sr."id", sr."status", sr."workerId"
    FROM
        step_runs sr
    FOR UPDATE SKIP LOCKED
)
UPDATE
    "StepRun"
SET
    "status" = 'PENDING_ASSIGNMENT',
    -- requeue after now plus 4 seconds
    "requeueAfter" = CURRENT_TIMESTAMP + INTERVAL '4 seconds',
    "updatedAt" = CURRENT_TIMESTAMP
FROM 
    locked_step_runs
WHERE
    "StepRun"."id" = locked_step_runs."id"
RETURNING "StepRun"."id"
`

// Count the total number of maxRuns - runningStepRuns across all workers
func (q *Queries) ListStepRunsToRequeue(ctx context.Context, db DBTX, tenantid pgtype.UUID) ([]pgtype.UUID, error) {
	rows, err := db.Query(ctx, listStepRunsToRequeue, tenantid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []pgtype.UUID
	for rows.Next() {
		var id pgtype.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const resolveLaterStepRuns = `-- name: ResolveLaterStepRuns :many
WITH RECURSIVE currStepRun AS (
  SELECT id, "createdAt", "updatedAt", "deletedAt", "tenantId", "jobRunId", "stepId", "order", "workerId", "tickerId", status, input, output, "requeueAfter", "scheduleTimeoutAt", error, "startedAt", "finishedAt", "timeoutAt", "cancelledAt", "cancelledReason", "cancelledError", "inputSchema", "callerFiles", "gitRepoBranch", "retryCount"
  FROM "StepRun"
  WHERE
    "id" = $2::uuid AND
    "tenantId" = $1::uuid
), childStepRuns AS (
  SELECT sr."id", sr."status"
  FROM "StepRun" sr
  JOIN "_StepRunOrder" sro ON sr."id" = sro."B"
  WHERE sro."A" = (SELECT "id" FROM currStepRun)
  
  UNION ALL
  
  SELECT sr."id", sr."status"
  FROM "StepRun" sr
  JOIN "_StepRunOrder" sro ON sr."id" = sro."B"
  JOIN childStepRuns csr ON sro."A" = csr."id"
)
UPDATE
    "StepRun" as sr
SET  "status" = CASE
    -- When the step is in a final state, it cannot be updated
    WHEN sr."status" IN ('SUCCEEDED', 'FAILED', 'CANCELLED') THEN sr."status"
    -- When the given step run has failed or been cancelled, then all child step runs are cancelled
    WHEN (SELECT "status" FROM currStepRun) IN ('FAILED', 'CANCELLED') THEN 'CANCELLED'
    ELSE sr."status"
    END,
    -- When the previous step run timed out, the cancelled reason is set
    "cancelledReason" = CASE
    -- When the step is in a final state, it cannot be updated
    WHEN sr."status" IN ('SUCCEEDED', 'FAILED', 'CANCELLED') THEN sr."cancelledReason"
    WHEN (SELECT "status" FROM currStepRun) = 'CANCELLED' AND (SELECT "cancelledReason" FROM currStepRun) = 'TIMED_OUT'::text THEN 'PREVIOUS_STEP_TIMED_OUT'
    WHEN (SELECT "status" FROM currStepRun) = 'FAILED' THEN 'PREVIOUS_STEP_FAILED'
    WHEN (SELECT "status" FROM currStepRun) = 'CANCELLED' THEN 'PREVIOUS_STEP_CANCELLED'
    ELSE NULL
    END
FROM
    childStepRuns csr
WHERE
    sr."id" = csr."id" AND
    sr."tenantId" = $1::uuid
RETURNING sr.id, sr."createdAt", sr."updatedAt", sr."deletedAt", sr."tenantId", sr."jobRunId", sr."stepId", sr."order", sr."workerId", sr."tickerId", sr.status, sr.input, sr.output, sr."requeueAfter", sr."scheduleTimeoutAt", sr.error, sr."startedAt", sr."finishedAt", sr."timeoutAt", sr."cancelledAt", sr."cancelledReason", sr."cancelledError", sr."inputSchema", sr."callerFiles", sr."gitRepoBranch", sr."retryCount"
`

type ResolveLaterStepRunsParams struct {
	Tenantid  pgtype.UUID `json:"tenantid"`
	Steprunid pgtype.UUID `json:"steprunid"`
}

func (q *Queries) ResolveLaterStepRuns(ctx context.Context, db DBTX, arg ResolveLaterStepRunsParams) ([]*StepRun, error) {
	rows, err := db.Query(ctx, resolveLaterStepRuns, arg.Tenantid, arg.Steprunid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*StepRun
	for rows.Next() {
		var i StepRun
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.TenantId,
			&i.JobRunId,
			&i.StepId,
			&i.Order,
			&i.WorkerId,
			&i.TickerId,
			&i.Status,
			&i.Input,
			&i.Output,
			&i.RequeueAfter,
			&i.ScheduleTimeoutAt,
			&i.Error,
			&i.StartedAt,
			&i.FinishedAt,
			&i.TimeoutAt,
			&i.CancelledAt,
			&i.CancelledReason,
			&i.CancelledError,
			&i.InputSchema,
			&i.CallerFiles,
			&i.GitRepoBranch,
			&i.RetryCount,
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

const unlinkStepRunFromWorker = `-- name: UnlinkStepRunFromWorker :one
UPDATE
    "StepRun"
SET
    "workerId" = NULL
WHERE
    "id" = $1::uuid AND
    "tenantId" = $2::uuid
RETURNING id, "createdAt", "updatedAt", "deletedAt", "tenantId", "jobRunId", "stepId", "order", "workerId", "tickerId", status, input, output, "requeueAfter", "scheduleTimeoutAt", error, "startedAt", "finishedAt", "timeoutAt", "cancelledAt", "cancelledReason", "cancelledError", "inputSchema", "callerFiles", "gitRepoBranch", "retryCount"
`

type UnlinkStepRunFromWorkerParams struct {
	Steprunid pgtype.UUID `json:"steprunid"`
	Tenantid  pgtype.UUID `json:"tenantid"`
}

func (q *Queries) UnlinkStepRunFromWorker(ctx context.Context, db DBTX, arg UnlinkStepRunFromWorkerParams) (*StepRun, error) {
	row := db.QueryRow(ctx, unlinkStepRunFromWorker, arg.Steprunid, arg.Tenantid)
	var i StepRun
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.JobRunId,
		&i.StepId,
		&i.Order,
		&i.WorkerId,
		&i.TickerId,
		&i.Status,
		&i.Input,
		&i.Output,
		&i.RequeueAfter,
		&i.ScheduleTimeoutAt,
		&i.Error,
		&i.StartedAt,
		&i.FinishedAt,
		&i.TimeoutAt,
		&i.CancelledAt,
		&i.CancelledReason,
		&i.CancelledError,
		&i.InputSchema,
		&i.CallerFiles,
		&i.GitRepoBranch,
		&i.RetryCount,
	)
	return &i, err
}

const updateStepRateLimits = `-- name: UpdateStepRateLimits :many
WITH step_rate_limits AS (
    SELECT
        rl."units" AS "units",
        rl."rateLimitKey" AS "rateLimitKey"
    FROM
        "StepRateLimit" rl
    WHERE
        rl."stepId" = $1::uuid AND
        rl."tenantId" = $2::uuid
), locked_rate_limits AS (
    SELECT
        srl."tenantId", srl.key, srl."limitValue", srl.value, srl."window", srl."lastRefill",
        step_rate_limits."units"
    FROM
        step_rate_limits
    JOIN
        "RateLimit" srl ON srl."key" = step_rate_limits."rateLimitKey" AND srl."tenantId" = $2::uuid
    FOR UPDATE
)
UPDATE
    "RateLimit" srl
SET
    "value" = get_refill_value(srl) - lrl."units",
    "lastRefill" = CASE
        WHEN NOW() - srl."lastRefill" >= srl."window"::INTERVAL THEN
            CURRENT_TIMESTAMP
        ELSE
            srl."lastRefill"
    END
FROM
    locked_rate_limits lrl
WHERE
    srl."tenantId" = lrl."tenantId" AND
    srl."key" = lrl."key"
RETURNING srl."tenantId", srl.key, srl."limitValue", srl.value, srl."window", srl."lastRefill"
`

type UpdateStepRateLimitsParams struct {
	Stepid   pgtype.UUID `json:"stepid"`
	Tenantid pgtype.UUID `json:"tenantid"`
}

func (q *Queries) UpdateStepRateLimits(ctx context.Context, db DBTX, arg UpdateStepRateLimitsParams) ([]*RateLimit, error) {
	rows, err := db.Query(ctx, updateStepRateLimits, arg.Stepid, arg.Tenantid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*RateLimit
	for rows.Next() {
		var i RateLimit
		if err := rows.Scan(
			&i.TenantId,
			&i.Key,
			&i.LimitValue,
			&i.Value,
			&i.Window,
			&i.LastRefill,
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

const updateStepRun = `-- name: UpdateStepRun :one
UPDATE
    "StepRun"
SET
    "requeueAfter" = COALESCE($1::timestamp, "requeueAfter"),
    "scheduleTimeoutAt" = COALESCE($2::timestamp, "scheduleTimeoutAt"),
    "startedAt" = COALESCE($3::timestamp, "startedAt"),
    "finishedAt" = CASE
        -- if this is a rerun, we clear the finishedAt
        WHEN $4::boolean THEN NULL
        ELSE  COALESCE($5::timestamp, "finishedAt")
    END,
    "status" = CASE 
        -- if this is a rerun, we permit status updates
        WHEN $4::boolean THEN COALESCE($6, "status")
        -- Final states are final, cannot be updated
        WHEN "status" IN ('SUCCEEDED', 'FAILED', 'CANCELLED') THEN "status"
        ELSE COALESCE($6, "status")
    END,
    "input" = COALESCE($7::jsonb, "input"),
    "output" = CASE
        -- if this is a rerun, we clear the output
        WHEN $4::boolean THEN NULL
        ELSE COALESCE($8::jsonb, "output")
    END,
    "error" = CASE
        -- if this is a rerun, we clear the error
        WHEN $4::boolean THEN NULL
        ELSE COALESCE($9::text, "error")
    END,
    "cancelledAt" = CASE
        -- if this is a rerun, we clear the cancelledAt
        WHEN $4::boolean THEN NULL
        ELSE COALESCE($10::timestamp, "cancelledAt")
    END,
    "cancelledReason" = CASE
        -- if this is a rerun, we clear the cancelledReason
        WHEN $4::boolean THEN NULL
        ELSE COALESCE($11::text, "cancelledReason")
    END,
    "retryCount" = COALESCE($12::int, "retryCount")
WHERE 
  "id" = $13::uuid AND
  "tenantId" = $14::uuid
RETURNING "StepRun".id, "StepRun"."createdAt", "StepRun"."updatedAt", "StepRun"."deletedAt", "StepRun"."tenantId", "StepRun"."jobRunId", "StepRun"."stepId", "StepRun"."order", "StepRun"."workerId", "StepRun"."tickerId", "StepRun".status, "StepRun".input, "StepRun".output, "StepRun"."requeueAfter", "StepRun"."scheduleTimeoutAt", "StepRun".error, "StepRun"."startedAt", "StepRun"."finishedAt", "StepRun"."timeoutAt", "StepRun"."cancelledAt", "StepRun"."cancelledReason", "StepRun"."cancelledError", "StepRun"."inputSchema", "StepRun"."callerFiles", "StepRun"."gitRepoBranch", "StepRun"."retryCount"
`

type UpdateStepRunParams struct {
	RequeueAfter      pgtype.Timestamp  `json:"requeueAfter"`
	ScheduleTimeoutAt pgtype.Timestamp  `json:"scheduleTimeoutAt"`
	StartedAt         pgtype.Timestamp  `json:"startedAt"`
	Rerun             pgtype.Bool       `json:"rerun"`
	FinishedAt        pgtype.Timestamp  `json:"finishedAt"`
	Status            NullStepRunStatus `json:"status"`
	Input             []byte            `json:"input"`
	Output            []byte            `json:"output"`
	Error             pgtype.Text       `json:"error"`
	CancelledAt       pgtype.Timestamp  `json:"cancelledAt"`
	CancelledReason   pgtype.Text       `json:"cancelledReason"`
	RetryCount        pgtype.Int4       `json:"retryCount"`
	ID                pgtype.UUID       `json:"id"`
	Tenantid          pgtype.UUID       `json:"tenantid"`
}

func (q *Queries) UpdateStepRun(ctx context.Context, db DBTX, arg UpdateStepRunParams) (*StepRun, error) {
	row := db.QueryRow(ctx, updateStepRun,
		arg.RequeueAfter,
		arg.ScheduleTimeoutAt,
		arg.StartedAt,
		arg.Rerun,
		arg.FinishedAt,
		arg.Status,
		arg.Input,
		arg.Output,
		arg.Error,
		arg.CancelledAt,
		arg.CancelledReason,
		arg.RetryCount,
		arg.ID,
		arg.Tenantid,
	)
	var i StepRun
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.JobRunId,
		&i.StepId,
		&i.Order,
		&i.WorkerId,
		&i.TickerId,
		&i.Status,
		&i.Input,
		&i.Output,
		&i.RequeueAfter,
		&i.ScheduleTimeoutAt,
		&i.Error,
		&i.StartedAt,
		&i.FinishedAt,
		&i.TimeoutAt,
		&i.CancelledAt,
		&i.CancelledReason,
		&i.CancelledError,
		&i.InputSchema,
		&i.CallerFiles,
		&i.GitRepoBranch,
		&i.RetryCount,
	)
	return &i, err
}

const updateStepRunInputSchema = `-- name: UpdateStepRunInputSchema :one
UPDATE
    "StepRun" sr
SET
    "inputSchema" = coalesce($1::jsonb, '{}'),
    "updatedAt" = CURRENT_TIMESTAMP
WHERE
    sr."tenantId" = $2::uuid AND
    sr."id" = $3::uuid
RETURNING "inputSchema"
`

type UpdateStepRunInputSchemaParams struct {
	InputSchema []byte      `json:"inputSchema"`
	Tenantid    pgtype.UUID `json:"tenantid"`
	Steprunid   pgtype.UUID `json:"steprunid"`
}

func (q *Queries) UpdateStepRunInputSchema(ctx context.Context, db DBTX, arg UpdateStepRunInputSchemaParams) ([]byte, error) {
	row := db.QueryRow(ctx, updateStepRunInputSchema, arg.InputSchema, arg.Tenantid, arg.Steprunid)
	var inputSchema []byte
	err := row.Scan(&inputSchema)
	return inputSchema, err
}

const updateStepRunOverridesData = `-- name: UpdateStepRunOverridesData :one
UPDATE
    "StepRun" AS sr
SET 
    "updatedAt" = CURRENT_TIMESTAMP,
    "input" = jsonb_set("input", $1::text[], $2::jsonb, true),
    "callerFiles" = jsonb_set("callerFiles", $3::text[], to_jsonb($4::text), true)
WHERE
    sr."tenantId" = $5::uuid AND
    sr."id" = $6::uuid
RETURNING "input"
`

type UpdateStepRunOverridesDataParams struct {
	Fieldpath    []string    `json:"fieldpath"`
	Jsondata     []byte      `json:"jsondata"`
	Overrideskey []string    `json:"overrideskey"`
	Callerfile   string      `json:"callerfile"`
	Tenantid     pgtype.UUID `json:"tenantid"`
	Steprunid    pgtype.UUID `json:"steprunid"`
}

func (q *Queries) UpdateStepRunOverridesData(ctx context.Context, db DBTX, arg UpdateStepRunOverridesDataParams) ([]byte, error) {
	row := db.QueryRow(ctx, updateStepRunOverridesData,
		arg.Fieldpath,
		arg.Jsondata,
		arg.Overrideskey,
		arg.Callerfile,
		arg.Tenantid,
		arg.Steprunid,
	)
	var input []byte
	err := row.Scan(&input)
	return input, err
}

const updateWorkerSemaphore = `-- name: UpdateWorkerSemaphore :one
WITH step_run AS (
    SELECT
        "id", "workerId"
    FROM
        "StepRun"
    WHERE
        "id" = $2::uuid AND
        "tenantId" = $3::uuid
), worker AS (
    SELECT
        "id",
        "maxRuns"
    FROM
        "Worker"
    WHERE
        "id" = (SELECT "workerId" FROM step_run)
)
UPDATE
    "WorkerSemaphore" ws
SET
    -- This shouldn't happen, but we set guardrails to prevent negative slots or slots over
    -- the worker's maxRuns
    "slots" = CASE 
        WHEN (ws."slots" + $1::int) < 0 THEN 0
        WHEN (ws."slots" + $1::int) > COALESCE(worker."maxRuns", 100) THEN COALESCE(worker."maxRuns", 100)
        ELSE (ws."slots" + $1::int)
    END
FROM
    worker
WHERE
    ws."workerId" = worker."id"
RETURNING ws."workerId", ws.slots
`

type UpdateWorkerSemaphoreParams struct {
	Inc       int32       `json:"inc"`
	Steprunid pgtype.UUID `json:"steprunid"`
	Tenantid  pgtype.UUID `json:"tenantid"`
}

func (q *Queries) UpdateWorkerSemaphore(ctx context.Context, db DBTX, arg UpdateWorkerSemaphoreParams) (*WorkerSemaphore, error) {
	row := db.QueryRow(ctx, updateWorkerSemaphore, arg.Inc, arg.Steprunid, arg.Tenantid)
	var i WorkerSemaphore
	err := row.Scan(&i.WorkerId, &i.Slots)
	return &i, err
}
