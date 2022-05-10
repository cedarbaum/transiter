// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: feed_queries.sql

package db

import (
	"context"
	"database/sql"
)

const deleteFeed = `-- name: DeleteFeed :exec
DELETE FROM feed WHERE pk = $1
`

func (q *Queries) DeleteFeed(ctx context.Context, pk int64) error {
	_, err := q.db.Exec(ctx, deleteFeed, pk)
	return err
}

const getFeedForUpdate = `-- name: GetFeedForUpdate :one
SELECT feed.pk, feed.id, feed.system_pk, feed.auto_update_enabled, feed.auto_update_period, feed.config FROM feed
    INNER JOIN feed_update ON feed_update.feed_pk = feed.pk
    WHERE feed_update.pk = $1
`

func (q *Queries) GetFeedForUpdate(ctx context.Context, updatePk int64) (Feed, error) {
	row := q.db.QueryRow(ctx, getFeedForUpdate, updatePk)
	var i Feed
	err := row.Scan(
		&i.Pk,
		&i.ID,
		&i.SystemPk,
		&i.AutoUpdateEnabled,
		&i.AutoUpdatePeriod,
		&i.Config,
	)
	return i, err
}

const getFeedInSystem = `-- name: GetFeedInSystem :one
SELECT feed.pk, feed.id, feed.system_pk, feed.auto_update_enabled, feed.auto_update_period, feed.config FROM feed
    INNER JOIN system ON feed.system_pk = system.pk
    WHERE system.id = $1
    AND feed.id = $2
`

type GetFeedInSystemParams struct {
	SystemID string
	FeedID   string
}

func (q *Queries) GetFeedInSystem(ctx context.Context, arg GetFeedInSystemParams) (Feed, error) {
	row := q.db.QueryRow(ctx, getFeedInSystem, arg.SystemID, arg.FeedID)
	var i Feed
	err := row.Scan(
		&i.Pk,
		&i.ID,
		&i.SystemPk,
		&i.AutoUpdateEnabled,
		&i.AutoUpdatePeriod,
		&i.Config,
	)
	return i, err
}

const insertFeed = `-- name: InsertFeed :exec
INSERT INTO feed
    (id, system_pk, auto_update_enabled, auto_update_period, config)
VALUES
    ($1, $2, $3, 
     $4, $5)
`

type InsertFeedParams struct {
	ID                string
	SystemPk          int64
	AutoUpdateEnabled bool
	AutoUpdatePeriod  sql.NullInt32
	Config            string
}

func (q *Queries) InsertFeed(ctx context.Context, arg InsertFeedParams) error {
	_, err := q.db.Exec(ctx, insertFeed,
		arg.ID,
		arg.SystemPk,
		arg.AutoUpdateEnabled,
		arg.AutoUpdatePeriod,
		arg.Config,
	)
	return err
}

const insertFeedUpdate = `-- name: InsertFeedUpdate :one
INSERT INTO feed_update
    (feed_pk, status)
VALUES
    ($1, $2)
RETURNING pk
`

type InsertFeedUpdateParams struct {
	FeedPk int64
	Status string
}

func (q *Queries) InsertFeedUpdate(ctx context.Context, arg InsertFeedUpdateParams) (int64, error) {
	row := q.db.QueryRow(ctx, insertFeedUpdate, arg.FeedPk, arg.Status)
	var pk int64
	err := row.Scan(&pk)
	return pk, err
}

const listAutoUpdateFeedsForSystem = `-- name: ListAutoUpdateFeedsForSystem :many
SELECT feed.id, feed.auto_update_period
FROM feed
    INNER JOIN system ON system.pk = feed.system_pk
WHERE feed.auto_update_enabled
    AND system.id = $1
`

type ListAutoUpdateFeedsForSystemRow struct {
	ID               string
	AutoUpdatePeriod sql.NullInt32
}

func (q *Queries) ListAutoUpdateFeedsForSystem(ctx context.Context, systemID string) ([]ListAutoUpdateFeedsForSystemRow, error) {
	rows, err := q.db.Query(ctx, listAutoUpdateFeedsForSystem, systemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListAutoUpdateFeedsForSystemRow
	for rows.Next() {
		var i ListAutoUpdateFeedsForSystemRow
		if err := rows.Scan(&i.ID, &i.AutoUpdatePeriod); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listFeedsInSystem = `-- name: ListFeedsInSystem :many
SELECT pk, id, system_pk, auto_update_enabled, auto_update_period, config FROM feed WHERE system_pk = $1 ORDER BY id
`

func (q *Queries) ListFeedsInSystem(ctx context.Context, systemPk int64) ([]Feed, error) {
	rows, err := q.db.Query(ctx, listFeedsInSystem, systemPk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feed
	for rows.Next() {
		var i Feed
		if err := rows.Scan(
			&i.Pk,
			&i.ID,
			&i.SystemPk,
			&i.AutoUpdateEnabled,
			&i.AutoUpdatePeriod,
			&i.Config,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateFeed = `-- name: UpdateFeed :exec
UPDATE feed
SET auto_update_enabled = $1, 
    auto_update_period = $2, 
    config = $3
WHERE pk = $4
`

type UpdateFeedParams struct {
	AutoUpdateEnabled bool
	AutoUpdatePeriod  sql.NullInt32
	Config            string
	FeedPk            int64
}

func (q *Queries) UpdateFeed(ctx context.Context, arg UpdateFeedParams) error {
	_, err := q.db.Exec(ctx, updateFeed,
		arg.AutoUpdateEnabled,
		arg.AutoUpdatePeriod,
		arg.Config,
		arg.FeedPk,
	)
	return err
}
