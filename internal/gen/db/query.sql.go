// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: query.sql

package db

import (
	"context"
	"database/sql"

	"github.com/jackc/pgtype"
)

const calculatePeriodicityForRoute = `-- name: CalculatePeriodicityForRoute :one
WITH route_stop_pks AS (
  SELECT DISTINCT trip_stop_time.stop_pk stop_pk FROM trip_stop_time
    INNER JOIN trip ON trip.pk = trip_stop_time.trip_pk
  WHERE trip.route_pk = $1
    AND NOT trip_stop_time.past
    AND trip_stop_time.arrival_time IS NOT NULL
), diffs AS (
  SELECT EXTRACT(epoch FROM MAX(trip_stop_time.arrival_time) - MIN(trip_stop_time.arrival_time)) diff, COUNT(*) n
  FROM trip_stop_time
    INNER JOIN route_stop_pks ON route_stop_pks.stop_pk = trip_stop_time.stop_pk
  GROUP BY trip_stop_time.stop_pk
  HAVING COUNT(*) > 1
)
SELECT coalesce(AVG(diff / (n-1)), -1) FROM diffs
`

func (q *Queries) CalculatePeriodicityForRoute(ctx context.Context, routePk int64) (interface{}, error) {
	row := q.db.QueryRow(ctx, calculatePeriodicityForRoute, routePk)
	var coalesce interface{}
	err := row.Scan(&coalesce)
	return coalesce, err
}

const countAgenciesInSystem = `-- name: CountAgenciesInSystem :one
SELECT COUNT(*) FROM agency WHERE system_pk = $1
`

func (q *Queries) CountAgenciesInSystem(ctx context.Context, systemPk int64) (int64, error) {
	row := q.db.QueryRow(ctx, countAgenciesInSystem, systemPk)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countFeedsInSystem = `-- name: CountFeedsInSystem :one
SELECT COUNT(*) FROM feed WHERE system_pk = $1
`

func (q *Queries) CountFeedsInSystem(ctx context.Context, systemPk int64) (int64, error) {
	row := q.db.QueryRow(ctx, countFeedsInSystem, systemPk)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countRoutesInSystem = `-- name: CountRoutesInSystem :one
SELECT COUNT(*) FROM route WHERE system_pk = $1
`

func (q *Queries) CountRoutesInSystem(ctx context.Context, systemPk int64) (int64, error) {
	row := q.db.QueryRow(ctx, countRoutesInSystem, systemPk)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countStopsInSystem = `-- name: CountStopsInSystem :one
SELECT COUNT(*) FROM stop WHERE system_pk = $1
`

func (q *Queries) CountStopsInSystem(ctx context.Context, systemPk int64) (int64, error) {
	row := q.db.QueryRow(ctx, countStopsInSystem, systemPk)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countSystems = `-- name: CountSystems :one
SELECT COUNT(*) FROM system
`

func (q *Queries) CountSystems(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countSystems)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countTransfersInSystem = `-- name: CountTransfersInSystem :one
SELECT COUNT(*) FROM transfer WHERE system_pk = $1
`

func (q *Queries) CountTransfersInSystem(ctx context.Context, systemPk sql.NullInt64) (int64, error) {
	row := q.db.QueryRow(ctx, countTransfersInSystem, systemPk)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getAgencyInSystem = `-- name: GetAgencyInSystem :one
SELECT agency.pk, agency.id, agency.system_pk, agency.source_pk, agency.name, agency.url, agency.timezone, agency.language, agency.phone, agency.fare_url, agency.email FROM agency
    INNER JOIN system ON agency.system_pk = system.pk
WHERE system.id = $1
    AND agency.id = $2
`

type GetAgencyInSystemParams struct {
	SystemID string
	AgencyID string
}

func (q *Queries) GetAgencyInSystem(ctx context.Context, arg GetAgencyInSystemParams) (Agency, error) {
	row := q.db.QueryRow(ctx, getAgencyInSystem, arg.SystemID, arg.AgencyID)
	var i Agency
	err := row.Scan(
		&i.Pk,
		&i.ID,
		&i.SystemPk,
		&i.SourcePk,
		&i.Name,
		&i.Url,
		&i.Timezone,
		&i.Language,
		&i.Phone,
		&i.FareUrl,
		&i.Email,
	)
	return i, err
}

const getLastStopsForTrips = `-- name: GetLastStopsForTrips :many
WITH last_stop_sequence AS (
  SELECT trip_pk, MAX(stop_sequence) as stop_sequence
    FROM trip_stop_time
    WHERE trip_pk = ANY($1::bigint[])
    GROUP BY trip_pk
)
SELECT lss.trip_pk, stop.id, stop.name
  FROM last_stop_sequence lss
  INNER JOIN trip_stop_time
    ON lss.trip_pk = trip_stop_time.trip_pk 
    AND lss.stop_sequence = trip_stop_time.stop_sequence
  INNER JOIN stop
    ON trip_stop_time.stop_pk = stop.pk
`

type GetLastStopsForTripsRow struct {
	TripPk int64
	ID     string
	Name   sql.NullString
}

func (q *Queries) GetLastStopsForTrips(ctx context.Context, tripPks []int64) ([]GetLastStopsForTripsRow, error) {
	rows, err := q.db.Query(ctx, getLastStopsForTrips, tripPks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetLastStopsForTripsRow
	for rows.Next() {
		var i GetLastStopsForTripsRow
		if err := rows.Scan(&i.TripPk, &i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRouteInSystem = `-- name: GetRouteInSystem :one
SELECT route.pk, route.id, route.system_pk, route.source_pk, route.color, route.text_color, route.short_name, route.long_name, route.description, route.url, route.sort_order, route.type, route.agency_pk, route.continuous_drop_off, route.continuous_pickup, agency.id agency_id, agency.name agency_name FROM route
    INNER JOIN system ON route.system_pk = system.pk
    INNER JOIN agency ON route.agency_pk = agency.pk
    WHERE system.id = $1
    AND route.id = $2
`

type GetRouteInSystemParams struct {
	SystemID string
	RouteID  string
}

type GetRouteInSystemRow struct {
	Pk                int64
	ID                string
	SystemPk          int64
	SourcePk          int64
	Color             string
	TextColor         string
	ShortName         sql.NullString
	LongName          sql.NullString
	Description       sql.NullString
	Url               sql.NullString
	SortOrder         sql.NullInt32
	Type              string
	AgencyPk          int64
	ContinuousDropOff string
	ContinuousPickup  string
	AgencyID          string
	AgencyName        string
}

func (q *Queries) GetRouteInSystem(ctx context.Context, arg GetRouteInSystemParams) (GetRouteInSystemRow, error) {
	row := q.db.QueryRow(ctx, getRouteInSystem, arg.SystemID, arg.RouteID)
	var i GetRouteInSystemRow
	err := row.Scan(
		&i.Pk,
		&i.ID,
		&i.SystemPk,
		&i.SourcePk,
		&i.Color,
		&i.TextColor,
		&i.ShortName,
		&i.LongName,
		&i.Description,
		&i.Url,
		&i.SortOrder,
		&i.Type,
		&i.AgencyPk,
		&i.ContinuousDropOff,
		&i.ContinuousPickup,
		&i.AgencyID,
		&i.AgencyName,
	)
	return i, err
}

const getStopInSystem = `-- name: GetStopInSystem :one
SELECT stop.pk, stop.id, system_pk, source_pk, parent_stop_pk, stop.name, longitude, latitude, url, code, description, platform_code, stop.timezone, type, wheelchair_boarding, zone_id, system.pk, system.id, system.name, system.timezone, status FROM stop
    INNER JOIN system ON stop.system_pk = system.pk
    WHERE system.id = $1
    AND stop.id = $2
`

type GetStopInSystemParams struct {
	SystemID string
	StopID   string
}

type GetStopInSystemRow struct {
	Pk                 int64
	ID                 string
	SystemPk           int64
	SourcePk           int64
	ParentStopPk       sql.NullInt64
	Name               sql.NullString
	Longitude          pgtype.Numeric
	Latitude           pgtype.Numeric
	Url                sql.NullString
	Code               sql.NullString
	Description        sql.NullString
	PlatformCode       sql.NullString
	Timezone           sql.NullString
	Type               string
	WheelchairBoarding string
	ZoneID             sql.NullString
	Pk_2               int64
	ID_2               string
	Name_2             string
	Timezone_2         sql.NullString
	Status             string
}

func (q *Queries) GetStopInSystem(ctx context.Context, arg GetStopInSystemParams) (GetStopInSystemRow, error) {
	row := q.db.QueryRow(ctx, getStopInSystem, arg.SystemID, arg.StopID)
	var i GetStopInSystemRow
	err := row.Scan(
		&i.Pk,
		&i.ID,
		&i.SystemPk,
		&i.SourcePk,
		&i.ParentStopPk,
		&i.Name,
		&i.Longitude,
		&i.Latitude,
		&i.Url,
		&i.Code,
		&i.Description,
		&i.PlatformCode,
		&i.Timezone,
		&i.Type,
		&i.WheelchairBoarding,
		&i.ZoneID,
		&i.Pk_2,
		&i.ID_2,
		&i.Name_2,
		&i.Timezone_2,
		&i.Status,
	)
	return i, err
}

const getSystem = `-- name: GetSystem :one
SELECT pk, id, name, timezone, status FROM system
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetSystem(ctx context.Context, id string) (System, error) {
	row := q.db.QueryRow(ctx, getSystem, id)
	var i System
	err := row.Scan(
		&i.Pk,
		&i.ID,
		&i.Name,
		&i.Timezone,
		&i.Status,
	)
	return i, err
}

const getTrip = `-- name: GetTrip :one
SELECT trip.pk, trip.id, trip.route_pk, trip.source_pk, trip.direction_id, trip.started_at, vehicle.id AS vehicle_id, route.id route_id, route.color route_color FROM trip
    INNER JOIN route ON route.pk = trip.route_pk
    INNER JOIN system ON system.pk = route.system_pk
    LEFT JOIN vehicle ON vehicle.trip_pk = trip.pk
WHERE trip.id = $1
    AND route.id = $2
    AND system.id = $3
`

type GetTripParams struct {
	TripID   string
	RouteID  string
	SystemID string
}

type GetTripRow struct {
	Pk          int64
	ID          string
	RoutePk     int64
	SourcePk    int64
	DirectionID sql.NullBool
	StartedAt   sql.NullTime
	VehicleID   sql.NullString
	RouteID     string
	RouteColor  string
}

func (q *Queries) GetTrip(ctx context.Context, arg GetTripParams) (GetTripRow, error) {
	row := q.db.QueryRow(ctx, getTrip, arg.TripID, arg.RouteID, arg.SystemID)
	var i GetTripRow
	err := row.Scan(
		&i.Pk,
		&i.ID,
		&i.RoutePk,
		&i.SourcePk,
		&i.DirectionID,
		&i.StartedAt,
		&i.VehicleID,
		&i.RouteID,
		&i.RouteColor,
	)
	return i, err
}

const listActiveAlertsForRoutes = `-- name: ListActiveAlertsForRoutes :many
SELECT route.pk route_pk, alert.pk, alert.id, alert.source_pk, alert.system_pk, alert.cause, alert.effect, alert.header, alert.description, alert.url, alert.hash, alert_active_period.starts_at, alert_active_period.ends_at
FROM route
    INNER JOIN alert_route ON route.pk = alert_route.route_pk
    INNER JOIN alert ON alert_route.alert_pk = alert.pk
    INNER JOIN alert_active_period ON alert_active_period.alert_pk = alert.pk
WHERE route.pk = ANY($1::bigint[])
    AND (
        alert_active_period.starts_at < $2
        OR alert_active_period.starts_at IS NULL
    )
    AND (
        alert_active_period.ends_at > $2
        OR alert_active_period.ends_at IS NULL
    )
ORDER BY alert.id ASC
`

type ListActiveAlertsForRoutesParams struct {
	RoutePks    []int64
	PresentTime sql.NullTime
}

type ListActiveAlertsForRoutesRow struct {
	RoutePk     int64
	Pk          int64
	ID          string
	SourcePk    int64
	SystemPk    int64
	Cause       string
	Effect      string
	Header      string
	Description string
	Url         string
	Hash        string
	StartsAt    sql.NullTime
	EndsAt      sql.NullTime
}

func (q *Queries) ListActiveAlertsForRoutes(ctx context.Context, arg ListActiveAlertsForRoutesParams) ([]ListActiveAlertsForRoutesRow, error) {
	rows, err := q.db.Query(ctx, listActiveAlertsForRoutes, arg.RoutePks, arg.PresentTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListActiveAlertsForRoutesRow
	for rows.Next() {
		var i ListActiveAlertsForRoutesRow
		if err := rows.Scan(
			&i.RoutePk,
			&i.Pk,
			&i.ID,
			&i.SourcePk,
			&i.SystemPk,
			&i.Cause,
			&i.Effect,
			&i.Header,
			&i.Description,
			&i.Url,
			&i.Hash,
			&i.StartsAt,
			&i.EndsAt,
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

const listActiveAlertsForStops = `-- name: ListActiveAlertsForStops :many
SELECT stop.pk stop_pk, alert.pk, alert.id, alert.cause, alert.effect, alert_active_period.starts_at, alert_active_period.ends_at
FROM stop
    INNER JOIN alert_stop ON stop.pk = alert_stop.stop_pk
    INNER JOIN alert ON alert_stop.alert_pk = alert.pk
    INNER JOIN alert_active_period ON alert_active_period.alert_pk = alert.pk
WHERE stop.pk = ANY($1::bigint[])
    AND (
        alert_active_period.starts_at < $2
        OR alert_active_period.starts_at IS NULL
    )
    AND (
        alert_active_period.ends_at > $2
        OR alert_active_period.ends_at IS NULL
    )
ORDER BY alert.id ASC
`

type ListActiveAlertsForStopsParams struct {
	StopPks     []int64
	PresentTime sql.NullTime
}

type ListActiveAlertsForStopsRow struct {
	StopPk   int64
	Pk       int64
	ID       string
	Cause    string
	Effect   string
	StartsAt sql.NullTime
	EndsAt   sql.NullTime
}

func (q *Queries) ListActiveAlertsForStops(ctx context.Context, arg ListActiveAlertsForStopsParams) ([]ListActiveAlertsForStopsRow, error) {
	rows, err := q.db.Query(ctx, listActiveAlertsForStops, arg.StopPks, arg.PresentTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListActiveAlertsForStopsRow
	for rows.Next() {
		var i ListActiveAlertsForStopsRow
		if err := rows.Scan(
			&i.StopPk,
			&i.Pk,
			&i.ID,
			&i.Cause,
			&i.Effect,
			&i.StartsAt,
			&i.EndsAt,
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

const listAgenciesInSystem = `-- name: ListAgenciesInSystem :many
SELECT pk, id, system_pk, source_pk, name, url, timezone, language, phone, fare_url, email FROM agency WHERE system_pk = $1 ORDER BY id
`

func (q *Queries) ListAgenciesInSystem(ctx context.Context, systemPk int64) ([]Agency, error) {
	rows, err := q.db.Query(ctx, listAgenciesInSystem, systemPk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Agency
	for rows.Next() {
		var i Agency
		if err := rows.Scan(
			&i.Pk,
			&i.ID,
			&i.SystemPk,
			&i.SourcePk,
			&i.Name,
			&i.Url,
			&i.Timezone,
			&i.Language,
			&i.Phone,
			&i.FareUrl,
			&i.Email,
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

const listRoutesByPk = `-- name: ListRoutesByPk :many
SELECT pk, id, system_pk, source_pk, color, text_color, short_name, long_name, description, url, sort_order, type, agency_pk, continuous_drop_off, continuous_pickup FROM route WHERE route.pk = ANY($1::bigint[])
`

func (q *Queries) ListRoutesByPk(ctx context.Context, routePks []int64) ([]Route, error) {
	rows, err := q.db.Query(ctx, listRoutesByPk, routePks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Route
	for rows.Next() {
		var i Route
		if err := rows.Scan(
			&i.Pk,
			&i.ID,
			&i.SystemPk,
			&i.SourcePk,
			&i.Color,
			&i.TextColor,
			&i.ShortName,
			&i.LongName,
			&i.Description,
			&i.Url,
			&i.SortOrder,
			&i.Type,
			&i.AgencyPk,
			&i.ContinuousDropOff,
			&i.ContinuousPickup,
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

const listRoutesInAgency = `-- name: ListRoutesInAgency :many
SELECT route.id, route.color FROM route
WHERE route.agency_pk = $1
`

type ListRoutesInAgencyRow struct {
	ID    string
	Color string
}

func (q *Queries) ListRoutesInAgency(ctx context.Context, agencyPk int64) ([]ListRoutesInAgencyRow, error) {
	rows, err := q.db.Query(ctx, listRoutesInAgency, agencyPk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListRoutesInAgencyRow
	for rows.Next() {
		var i ListRoutesInAgencyRow
		if err := rows.Scan(&i.ID, &i.Color); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listRoutesInSystem = `-- name: ListRoutesInSystem :many
SELECT pk, id, system_pk, source_pk, color, text_color, short_name, long_name, description, url, sort_order, type, agency_pk, continuous_drop_off, continuous_pickup FROM route WHERE system_pk = $1 ORDER BY id
`

func (q *Queries) ListRoutesInSystem(ctx context.Context, systemPk int64) ([]Route, error) {
	rows, err := q.db.Query(ctx, listRoutesInSystem, systemPk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Route
	for rows.Next() {
		var i Route
		if err := rows.Scan(
			&i.Pk,
			&i.ID,
			&i.SystemPk,
			&i.SourcePk,
			&i.Color,
			&i.TextColor,
			&i.ShortName,
			&i.LongName,
			&i.Description,
			&i.Url,
			&i.SortOrder,
			&i.Type,
			&i.AgencyPk,
			&i.ContinuousDropOff,
			&i.ContinuousPickup,
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

const listServiceMapsConfigIDsForStops = `-- name: ListServiceMapsConfigIDsForStops :many
SELECT stop.pk, service_map_config.id
FROM service_map_config
    INNER JOIN stop ON service_map_config.system_pk = stop.system_pk
WHERE service_map_config.default_for_routes_at_stop
    AND stop.pk = ANY($1::bigint[])
`

type ListServiceMapsConfigIDsForStopsRow struct {
	Pk int64
	ID string
}

func (q *Queries) ListServiceMapsConfigIDsForStops(ctx context.Context, stopPks []int64) ([]ListServiceMapsConfigIDsForStopsRow, error) {
	rows, err := q.db.Query(ctx, listServiceMapsConfigIDsForStops, stopPks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListServiceMapsConfigIDsForStopsRow
	for rows.Next() {
		var i ListServiceMapsConfigIDsForStopsRow
		if err := rows.Scan(&i.Pk, &i.ID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listServiceMapsForRoute = `-- name: ListServiceMapsForRoute :many
SELECT DISTINCT service_map_config.id config_id, service_map_vertex.position, stop.id stop_id, stop.name stop_name
FROM service_map_config
  INNER JOIN system ON service_map_config.system_pk = system.pk
  INNER JOIN route ON route.system_pk = system.pk
  LEFT JOIN service_map ON service_map.config_pk = service_map_config.pk AND service_map.route_pk = $1
  LEFT JOIN service_map_vertex ON service_map_vertex.map_pk = service_map.pk
  LEFT JOIN stop ON stop.pk = service_map_vertex.stop_pk
WHERE service_map_config.default_for_stops_in_route AND route.pk = $1
ORDER BY service_map_config.id, service_map_vertex.position
`

type ListServiceMapsForRouteRow struct {
	ConfigID string
	Position sql.NullInt32
	StopID   sql.NullString
	StopName sql.NullString
}

func (q *Queries) ListServiceMapsForRoute(ctx context.Context, routePk int64) ([]ListServiceMapsForRouteRow, error) {
	rows, err := q.db.Query(ctx, listServiceMapsForRoute, routePk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListServiceMapsForRouteRow
	for rows.Next() {
		var i ListServiceMapsForRouteRow
		if err := rows.Scan(
			&i.ConfigID,
			&i.Position,
			&i.StopID,
			&i.StopName,
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

const listServiceMapsForStops = `-- name: ListServiceMapsForStops :many
WITH RECURSIVE descendent AS (
	SELECT initial.pk, initial.parent_stop_pk, initial.pk AS descendent_pk
	  FROM stop initial
    WHERE initial.pk = ANY($1::bigint[])
	UNION (
    SELECT parent.pk, parent.parent_stop_pk, descendent.pk AS descendent_pk
      FROM stop parent
      INNER JOIN descendent ON (
        descendent.parent_stop_pk = parent.pk OR
        descendent.pk = parent.pk
      )
  )
)
SELECT descendent.pk stop_pk, service_map_config.id service_map_config_id,
  route.id route_id, route.color route_color, system.id system_id
FROM descendent
  LEFT JOIN service_map_vertex smv ON smv.stop_pk = descendent.descendent_pk
  INNER JOIN service_map ON service_map.pk = smv.map_pk
  INNER JOIN service_map_config ON service_map_config.pk = service_map.config_pk
  LEFT JOIN route ON service_map.route_pk = route.pk
  INNER JOIN system ON system.pk = route.system_pk
WHERE service_map_config.default_for_routes_at_stop
ORDER BY system_id, route_id
`

type ListServiceMapsForStopsRow struct {
	StopPk             int64
	ServiceMapConfigID string
	RouteID            string
	RouteColor         string
	SystemID           string
}

func (q *Queries) ListServiceMapsForStops(ctx context.Context, stopPks []int64) ([]ListServiceMapsForStopsRow, error) {
	rows, err := q.db.Query(ctx, listServiceMapsForStops, stopPks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListServiceMapsForStopsRow
	for rows.Next() {
		var i ListServiceMapsForStopsRow
		if err := rows.Scan(
			&i.StopPk,
			&i.ServiceMapConfigID,
			&i.RouteID,
			&i.RouteColor,
			&i.SystemID,
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

const listStopHeadsignRulesForStops = `-- name: ListStopHeadsignRulesForStops :many
SELECT pk, source_pk, priority, stop_pk, track, headsign FROM stop_headsign_rule
WHERE stop_pk = ANY($1::bigint[])
ORDER BY priority ASC
`

func (q *Queries) ListStopHeadsignRulesForStops(ctx context.Context, stopPks []int64) ([]StopHeadsignRule, error) {
	rows, err := q.db.Query(ctx, listStopHeadsignRulesForStops, stopPks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []StopHeadsignRule
	for rows.Next() {
		var i StopHeadsignRule
		if err := rows.Scan(
			&i.Pk,
			&i.SourcePk,
			&i.Priority,
			&i.StopPk,
			&i.Track,
			&i.Headsign,
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

const listStopTimesAtStops = `-- name: ListStopTimesAtStops :many
SELECT trip_stop_time.pk, trip_stop_time.stop_pk, trip_stop_time.trip_pk, trip_stop_time.arrival_time, trip_stop_time.arrival_delay, trip_stop_time.arrival_uncertainty, trip_stop_time.departure_time, trip_stop_time.departure_delay, trip_stop_time.departure_uncertainty, trip_stop_time.stop_sequence, trip_stop_time.track, trip_stop_time.past, trip.pk, trip.id, trip.route_pk, trip.source_pk, trip.direction_id, trip.started_at, vehicle.id vehicle_id FROM trip_stop_time
    INNER JOIN trip ON trip_stop_time.trip_pk = trip.pk
    LEFT JOIN vehicle ON vehicle.trip_pk = trip.pk
    WHERE trip_stop_time.stop_pk = ANY($1::bigint[])
    AND NOT trip_stop_time.past
    ORDER BY trip_stop_time.departure_time, trip_stop_time.arrival_time
`

type ListStopTimesAtStopsRow struct {
	Pk                   int64
	StopPk               int64
	TripPk               int64
	ArrivalTime          sql.NullTime
	ArrivalDelay         sql.NullInt32
	ArrivalUncertainty   sql.NullInt32
	DepartureTime        sql.NullTime
	DepartureDelay       sql.NullInt32
	DepartureUncertainty sql.NullInt32
	StopSequence         int32
	Track                sql.NullString
	Past                 bool
	Pk_2                 int64
	ID                   string
	RoutePk              int64
	SourcePk             int64
	DirectionID          sql.NullBool
	StartedAt            sql.NullTime
	VehicleID            sql.NullString
}

func (q *Queries) ListStopTimesAtStops(ctx context.Context, stopPks []int64) ([]ListStopTimesAtStopsRow, error) {
	rows, err := q.db.Query(ctx, listStopTimesAtStops, stopPks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListStopTimesAtStopsRow
	for rows.Next() {
		var i ListStopTimesAtStopsRow
		if err := rows.Scan(
			&i.Pk,
			&i.StopPk,
			&i.TripPk,
			&i.ArrivalTime,
			&i.ArrivalDelay,
			&i.ArrivalUncertainty,
			&i.DepartureTime,
			&i.DepartureDelay,
			&i.DepartureUncertainty,
			&i.StopSequence,
			&i.Track,
			&i.Past,
			&i.Pk_2,
			&i.ID,
			&i.RoutePk,
			&i.SourcePk,
			&i.DirectionID,
			&i.StartedAt,
			&i.VehicleID,
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

const listStopsInStopTree = `-- name: ListStopsInStopTree :many
WITH RECURSIVE 
ancestor AS (
	SELECT initial.pk, initial.parent_stop_pk
	  FROM stop initial
	  WHERE	initial.pk = $1
	UNION
	SELECT parent.pk, parent.parent_stop_pk
		FROM stop parent
		INNER JOIN ancestor ON ancestor.parent_stop_pk = parent.pk
),
descendent AS (
	SELECT pk, parent_stop_pk FROM ancestor
	UNION
	SELECT child.pk, child.parent_stop_pk
		FROM stop child
		INNER JOIN descendent ON descendent.pk = child.parent_stop_pk
) 
SELECT stop.pk, stop.id, stop.system_pk, stop.source_pk, stop.parent_stop_pk, stop.name, stop.longitude, stop.latitude, stop.url, stop.code, stop.description, stop.platform_code, stop.timezone, stop.type, stop.wheelchair_boarding, stop.zone_id FROM stop
  INNER JOIN descendent
  ON stop.pk = descendent.pk
`

func (q *Queries) ListStopsInStopTree(ctx context.Context, pk int64) ([]Stop, error) {
	rows, err := q.db.Query(ctx, listStopsInStopTree, pk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Stop
	for rows.Next() {
		var i Stop
		if err := rows.Scan(
			&i.Pk,
			&i.ID,
			&i.SystemPk,
			&i.SourcePk,
			&i.ParentStopPk,
			&i.Name,
			&i.Longitude,
			&i.Latitude,
			&i.Url,
			&i.Code,
			&i.Description,
			&i.PlatformCode,
			&i.Timezone,
			&i.Type,
			&i.WheelchairBoarding,
			&i.ZoneID,
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

const listStopsInSystem = `-- name: ListStopsInSystem :many
SELECT pk, id, system_pk, source_pk, parent_stop_pk, name, longitude, latitude, url, code, description, platform_code, timezone, type, wheelchair_boarding, zone_id FROM stop WHERE system_pk = $1 
    ORDER BY id
`

func (q *Queries) ListStopsInSystem(ctx context.Context, systemPk int64) ([]Stop, error) {
	rows, err := q.db.Query(ctx, listStopsInSystem, systemPk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Stop
	for rows.Next() {
		var i Stop
		if err := rows.Scan(
			&i.Pk,
			&i.ID,
			&i.SystemPk,
			&i.SourcePk,
			&i.ParentStopPk,
			&i.Name,
			&i.Longitude,
			&i.Latitude,
			&i.Url,
			&i.Code,
			&i.Description,
			&i.PlatformCode,
			&i.Timezone,
			&i.Type,
			&i.WheelchairBoarding,
			&i.ZoneID,
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

const listStopsTimesForTrip = `-- name: ListStopsTimesForTrip :many
SELECT trip_stop_time.pk, trip_stop_time.stop_pk, trip_stop_time.trip_pk, trip_stop_time.arrival_time, trip_stop_time.arrival_delay, trip_stop_time.arrival_uncertainty, trip_stop_time.departure_time, trip_stop_time.departure_delay, trip_stop_time.departure_uncertainty, trip_stop_time.stop_sequence, trip_stop_time.track, trip_stop_time.past, stop.id stop_id, stop.name stop_name
FROM trip_stop_time
    INNER JOIN stop ON trip_stop_time.stop_pk = stop.pk
WHERE trip_stop_time.trip_pk = $1
ORDER BY trip_stop_time.stop_sequence ASC
`

type ListStopsTimesForTripRow struct {
	Pk                   int64
	StopPk               int64
	TripPk               int64
	ArrivalTime          sql.NullTime
	ArrivalDelay         sql.NullInt32
	ArrivalUncertainty   sql.NullInt32
	DepartureTime        sql.NullTime
	DepartureDelay       sql.NullInt32
	DepartureUncertainty sql.NullInt32
	StopSequence         int32
	Track                sql.NullString
	Past                 bool
	StopID               string
	StopName             sql.NullString
}

func (q *Queries) ListStopsTimesForTrip(ctx context.Context, tripPk int64) ([]ListStopsTimesForTripRow, error) {
	rows, err := q.db.Query(ctx, listStopsTimesForTrip, tripPk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListStopsTimesForTripRow
	for rows.Next() {
		var i ListStopsTimesForTripRow
		if err := rows.Scan(
			&i.Pk,
			&i.StopPk,
			&i.TripPk,
			&i.ArrivalTime,
			&i.ArrivalDelay,
			&i.ArrivalUncertainty,
			&i.DepartureTime,
			&i.DepartureDelay,
			&i.DepartureUncertainty,
			&i.StopSequence,
			&i.Track,
			&i.Past,
			&i.StopID,
			&i.StopName,
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

const listSystems = `-- name: ListSystems :many
SELECT pk, id, name, timezone, status FROM system ORDER BY id, name
`

func (q *Queries) ListSystems(ctx context.Context) ([]System, error) {
	rows, err := q.db.Query(ctx, listSystems)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []System
	for rows.Next() {
		var i System
		if err := rows.Scan(
			&i.Pk,
			&i.ID,
			&i.Name,
			&i.Timezone,
			&i.Status,
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

const listTransfersFromStops = `-- name: ListTransfersFromStops :many
  SELECT transfer.from_stop_pk,
      transfer.to_stop_pk, stop.id to_id, stop.name to_name, 
      transfer.type, transfer.min_transfer_time, transfer.distance
  FROM transfer
  INNER JOIN stop
    ON stop.pk = transfer.to_stop_pk
  WHERE transfer.from_stop_pk = ANY($1::bigint[])
`

type ListTransfersFromStopsRow struct {
	FromStopPk      int64
	ToStopPk        int64
	ToID            string
	ToName          sql.NullString
	Type            string
	MinTransferTime sql.NullInt32
	Distance        sql.NullInt32
}

func (q *Queries) ListTransfersFromStops(ctx context.Context, fromStopPks []int64) ([]ListTransfersFromStopsRow, error) {
	rows, err := q.db.Query(ctx, listTransfersFromStops, fromStopPks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListTransfersFromStopsRow
	for rows.Next() {
		var i ListTransfersFromStopsRow
		if err := rows.Scan(
			&i.FromStopPk,
			&i.ToStopPk,
			&i.ToID,
			&i.ToName,
			&i.Type,
			&i.MinTransferTime,
			&i.Distance,
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

const listTransfersInSystem = `-- name: ListTransfersInSystem :many
SELECT 
    transfer.pk, transfer.source_pk, transfer.config_source_pk, transfer.system_pk, transfer.from_stop_pk, transfer.to_stop_pk, transfer.type, transfer.min_transfer_time, transfer.distance,
    from_stop.id from_stop_id, from_stop.name from_stop_name, from_system.id from_system_id,
    to_stop.id to_stop_id, to_stop.name to_stop_name, to_system.id to_system_id
FROM transfer
    INNER JOIN stop from_stop ON from_stop.pk = transfer.from_stop_pk
    INNER JOIN system from_system ON from_stop.system_pk = from_system.pk
    INNER JOIN stop to_stop ON to_stop.pk = transfer.to_stop_pk
    INNER JOIN system to_system ON to_stop.system_pk = to_system.pk
WHERE transfer.system_pk = $1 
ORDER BY transfer.pk
`

type ListTransfersInSystemRow struct {
	Pk              int64
	SourcePk        sql.NullInt64
	ConfigSourcePk  sql.NullInt64
	SystemPk        sql.NullInt64
	FromStopPk      int64
	ToStopPk        int64
	Type            string
	MinTransferTime sql.NullInt32
	Distance        sql.NullInt32
	FromStopID      string
	FromStopName    sql.NullString
	FromSystemID    string
	ToStopID        string
	ToStopName      sql.NullString
	ToSystemID      string
}

func (q *Queries) ListTransfersInSystem(ctx context.Context, systemPk sql.NullInt64) ([]ListTransfersInSystemRow, error) {
	rows, err := q.db.Query(ctx, listTransfersInSystem, systemPk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListTransfersInSystemRow
	for rows.Next() {
		var i ListTransfersInSystemRow
		if err := rows.Scan(
			&i.Pk,
			&i.SourcePk,
			&i.ConfigSourcePk,
			&i.SystemPk,
			&i.FromStopPk,
			&i.ToStopPk,
			&i.Type,
			&i.MinTransferTime,
			&i.Distance,
			&i.FromStopID,
			&i.FromStopName,
			&i.FromSystemID,
			&i.ToStopID,
			&i.ToStopName,
			&i.ToSystemID,
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

const listTripsInRoute = `-- name: ListTripsInRoute :many
SELECT trip.pk, trip.id, trip.route_pk, trip.source_pk, trip.direction_id, trip.started_at, vehicle.id vehicle_id FROM trip 
    LEFT JOIN vehicle ON vehicle.trip_pk = trip.pk
WHERE trip.route_pk = $1
ORDER BY trip.id
`

type ListTripsInRouteRow struct {
	Pk          int64
	ID          string
	RoutePk     int64
	SourcePk    int64
	DirectionID sql.NullBool
	StartedAt   sql.NullTime
	VehicleID   sql.NullString
}

func (q *Queries) ListTripsInRoute(ctx context.Context, routePk int64) ([]ListTripsInRouteRow, error) {
	rows, err := q.db.Query(ctx, listTripsInRoute, routePk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListTripsInRouteRow
	for rows.Next() {
		var i ListTripsInRouteRow
		if err := rows.Scan(
			&i.Pk,
			&i.ID,
			&i.RoutePk,
			&i.SourcePk,
			&i.DirectionID,
			&i.StartedAt,
			&i.VehicleID,
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

const listUpdatesInFeed = `-- name: ListUpdatesInFeed :many
SELECT pk, feed_pk, status, started_at, ended_at, result, content_length, content_hash, error_message FROM feed_update 
WHERE feed_pk = $1
ORDER BY pk DESC
LIMIT 100
`

func (q *Queries) ListUpdatesInFeed(ctx context.Context, feedPk int64) ([]FeedUpdate, error) {
	rows, err := q.db.Query(ctx, listUpdatesInFeed, feedPk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FeedUpdate
	for rows.Next() {
		var i FeedUpdate
		if err := rows.Scan(
			&i.Pk,
			&i.FeedPk,
			&i.Status,
			&i.StartedAt,
			&i.EndedAt,
			&i.Result,
			&i.ContentLength,
			&i.ContentHash,
			&i.ErrorMessage,
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
