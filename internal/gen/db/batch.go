// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: batch.go

package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const markTripStopTimesPast = `-- name: MarkTripStopTimesPast :batchexec
UPDATE trip_stop_time
SET
    past = TRUE
WHERE
    trip_pk = $1
    AND stop_sequence < $2
`

type MarkTripStopTimesPastBatchResults struct {
	br     pgx.BatchResults
	tot    int
	closed bool
}

type MarkTripStopTimesPastParams struct {
	TripPk              int64
	CurrentStopSequence int32
}

func (q *Queries) MarkTripStopTimesPast(ctx context.Context, arg []MarkTripStopTimesPastParams) *MarkTripStopTimesPastBatchResults {
	batch := &pgx.Batch{}
	for _, a := range arg {
		vals := []interface{}{
			a.TripPk,
			a.CurrentStopSequence,
		}
		batch.Queue(markTripStopTimesPast, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &MarkTripStopTimesPastBatchResults{br, len(arg), false}
}

func (b *MarkTripStopTimesPastBatchResults) Exec(f func(int, error)) {
	defer b.br.Close()
	for t := 0; t < b.tot; t++ {
		if b.closed {
			if f != nil {
				f(t, errors.New("batch already closed"))
			}
			continue
		}
		_, err := b.br.Exec()
		if f != nil {
			f(t, err)
		}
	}
}

func (b *MarkTripStopTimesPastBatchResults) Close() error {
	b.closed = true
	return b.br.Close()
}

const updateTrip = `-- name: UpdateTrip :batchexec
UPDATE trip SET
    feed_pk = $1,
    direction_id = $2,
    started_at = $3,
    gtfs_hash = $4
WHERE pk = $5
`

type UpdateTripBatchResults struct {
	br     pgx.BatchResults
	tot    int
	closed bool
}

type UpdateTripParams struct {
	FeedPk      int64
	DirectionID pgtype.Bool
	StartedAt   pgtype.Timestamptz
	GtfsHash    string
	Pk          int64
}

func (q *Queries) UpdateTrip(ctx context.Context, arg []UpdateTripParams) *UpdateTripBatchResults {
	batch := &pgx.Batch{}
	for _, a := range arg {
		vals := []interface{}{
			a.FeedPk,
			a.DirectionID,
			a.StartedAt,
			a.GtfsHash,
			a.Pk,
		}
		batch.Queue(updateTrip, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &UpdateTripBatchResults{br, len(arg), false}
}

func (b *UpdateTripBatchResults) Exec(f func(int, error)) {
	defer b.br.Close()
	for t := 0; t < b.tot; t++ {
		if b.closed {
			if f != nil {
				f(t, errors.New("batch already closed"))
			}
			continue
		}
		_, err := b.br.Exec()
		if f != nil {
			f(t, err)
		}
	}
}

func (b *UpdateTripBatchResults) Close() error {
	b.closed = true
	return b.br.Close()
}
