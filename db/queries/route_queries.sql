-- name: InsertRoute :one
INSERT INTO route
    (id, system_pk, source_pk, color, text_color,
     short_name, long_name, description, url, sort_order,
     type, continuous_pickup, continuous_drop_off, agency_pk)
VALUES
    (sqlc.arg(id), sqlc.arg(system_pk), sqlc.arg(source_pk), sqlc.arg(color), sqlc.arg(text_color),
     sqlc.arg(short_name), sqlc.arg(long_name), sqlc.arg(description), sqlc.arg(url), sqlc.arg(sort_order),
     sqlc.arg(type), sqlc.arg(continuous_pickup),sqlc.arg(continuous_drop_off), sqlc.arg(agency_pk))
RETURNING pk;

-- name: UpdateRoute :exec
UPDATE route SET
    source_pk = sqlc.arg(source_pk),
    color = sqlc.arg(color),
    text_color = sqlc.arg(text_color),
    short_name = sqlc.arg(short_name),
    long_name = sqlc.arg(long_name),
    description = sqlc.arg(description),
    url = sqlc.arg(url),
    sort_order = sqlc.arg(sort_order),
    type = sqlc.arg(type),
    continuous_pickup = sqlc.arg(continuous_pickup),
    continuous_drop_off = sqlc.arg(continuous_drop_off),
    agency_pk = sqlc.arg(agency_pk)
WHERE
    pk = sqlc.arg(pk);

-- name: DeleteStaleRoutes :many
DELETE FROM route
USING feed_update
WHERE
    feed_update.pk = route.source_pk
    AND feed_update.feed_pk = sqlc.arg(feed_pk)
    AND feed_update.pk != sqlc.arg(update_pk)
RETURNING route.id;

-- name: ListRoutes :many
SELECT * FROM route WHERE system_pk = $1 ORDER BY id;

-- name: GetRoute :one
SELECT route.* FROM route
    INNER JOIN system ON route.system_pk = system.pk
    WHERE system.pk = sqlc.arg(system_pk)
    AND route.id = sqlc.arg(route_id);

-- name: MapRouteIDToPkInSystem :many
SELECT pk, id from route
WHERE
    system_pk = sqlc.arg(system_pk)
    AND (
        NOT sqlc.arg(filter_by_route_id)::bool
        OR id = ANY(sqlc.arg(route_ids)::text[])
    );

-- name: EstimateHeadwaysForRoutes :many
WITH per_stop_data AS (
    SELECT
        trip.route_pk route_pk,
        EXTRACT(epoch FROM MAX(trip_stop_time.arrival_time) - MIN(trip_stop_time.arrival_time)) total_diff,
        COUNT(*)-1 num_diffs
    FROM trip_stop_time
        INNER JOIN trip ON trip.pk = trip_stop_time.trip_pk
    WHERE trip.route_pk = ANY(sqlc.arg(route_pks)::bigint[])
        AND NOT trip_stop_time.past
        AND trip_stop_time.arrival_time IS NOT NULL
        AND trip_stop_time.arrival_time >= sqlc.arg(present_time)
    GROUP BY trip_stop_time.stop_pk, trip.route_pk
        HAVING COUNT(*) > 1
)
SELECT
    route_pk,
    COALESCE(ROUND(SUM(total_diff) / (SUM(num_diffs)))::integer, -1)::integer estimated_headway
FROM per_stop_data
GROUP BY route_pk;

-- name: ListRoutesByPk :many
SELECT route.pk, route.id id, route.color, system.id system_id
FROM route
    INNER JOIN system on route.system_pk = system.pk
WHERE route.pk = ANY(sqlc.arg(route_pks)::bigint[]);

-- name: ListRoutesInAgency :many
SELECT route.id, route.color FROM route
WHERE route.agency_pk = sqlc.arg(agency_pk);

-- name: ListTripsForRoutes :many
SELECT route.pk route_pk, trip.id trip_id FROM route
    INNER JOIN trip ON route.pk = trip.route_pk
WHERE route_pk = ANY(sqlc.arg(route_pks)::bigint[])
GROUP BY route.pk, trip.id;
