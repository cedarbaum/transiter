-- name: InsertVehicle :exec
INSERT INTO vehicle
    (id, system_pk, trip_id, label, license_plate, current_status, latitude, longitude, bearing, odometer, speed, congestion_level, updated_at, current_stop_pk, current_stop_sequence, occupancy_status, feed_pk, occupancy_percentage)
VALUES
    (sqlc.arg(id), sqlc.arg(system_pk), sqlc.arg(trip_id), sqlc.arg(label), sqlc.arg(license_plate), sqlc.arg(current_status), sqlc.arg(latitude), sqlc.arg(longitude), sqlc.arg(bearing), sqlc.arg(odometer), sqlc.arg(speed), sqlc.arg(congestion_level), sqlc.arg(updated_at), sqlc.arg(current_stop_pk), sqlc.arg(current_stop_sequence), sqlc.arg(occupancy_status), sqlc.arg(feed_pk), sqlc.arg(occupancy_percentage));

-- name: UpdateVehicle :exec
UPDATE vehicle
SET trip_id = sqlc.arg(trip_id),
    label = sqlc.arg(label),
    license_plate = sqlc.arg(license_plate),
    current_status = sqlc.arg(current_status),
    latitude = sqlc.arg(latitude),
    longitude = sqlc.arg(longitude),
    bearing = sqlc.arg(bearing),
    odometer = sqlc.arg(odometer),
    speed = sqlc.arg(speed),
    congestion_level = sqlc.arg(congestion_level),
    updated_at = sqlc.arg(updated_at),
    current_stop_pk = sqlc.arg(current_stop_pk),
    current_stop_sequence = sqlc.arg(current_stop_sequence),
    occupancy_status = sqlc.arg(occupancy_status),
    feed_pk = sqlc.arg(feed_pk),
    occupancy_percentage = sqlc.arg(occupancy_percentage)
WHERE vehicle.pk = sqlc.arg(pk);

-- name: ListVehicles :many
WITH trip_with_system_pk AS (
  SELECT trip.*, feed.system_pk
  FROM trip
  LEFT JOIN feed ON trip.feed_pk = feed.pk
)
SELECT vehicle.*,
       stop.id as stop_id,
       stop.name as stop_name,
       trip_with_system_pk.id as trip_id,
       trip_with_system_pk.direction_id as trip_direction_id,
       route.id as route_id,
       route.color as route_color
FROM vehicle
LEFT JOIN stop ON vehicle.current_stop_pk = stop.pk
LEFT JOIN trip_with_system_pk ON vehicle.trip_id = trip_with_system_pk.id AND vehicle.system_pk = trip_with_system_pk.system_pk
LEFT JOIN route ON trip_with_system_pk.route_pk = route.pk
WHERE vehicle.system_pk = sqlc.arg(system_pk)
  AND vehicle.id >= sqlc.arg(first_vehicle_id)
  AND (
    NOT sqlc.arg(only_return_specified_ids)::bool OR
    vehicle.id = ANY(sqlc.arg(vehicle_ids)::text[])
  )
ORDER BY vehicle.id
LIMIT sqlc.arg(num_vehicles);

-- name: ListVehicles_Geographic :many
WITH distance AS (
  SELECT
  pk vehicle_pk,
  (6371 * acos(cos(radians(latitude)) * cos(radians(sqlc.arg(latitude)::numeric)) * cos(radians(sqlc.arg(longitude)::numeric) - radians(longitude)) + sin(radians(latitude)) * sin(radians(sqlc.arg(latitude)::numeric)))) val
  FROM vehicle
  WHERE vehicle.system_pk = sqlc.arg(system_pk) AND latitude IS NOT NULL AND longitude IS NOT NULL
), trip_with_system_pk AS (
  SELECT trip.*, feed.system_pk
  FROM trip
  LEFT JOIN feed ON trip.feed_pk = feed.pk
)
SELECT vehicle.*,
       stop.id as stop_id,
       stop.name as stop_name,
       trip_with_system_pk.id as trip_id,
       trip_with_system_pk.direction_id as trip_direction_id,
       route.id as route_id,
       route.color as route_color
FROM vehicle
INNER JOIN distance ON vehicle.pk = distance.vehicle_pk
AND distance.val <= sqlc.arg(max_distance)::numeric
LEFT JOIN stop ON vehicle.current_stop_pk = stop.pk
LEFT JOIN trip_with_system_pk ON vehicle.trip_id = trip_with_system_pk.id AND vehicle.system_pk = trip_with_system_pk.system_pk
LEFT JOIN route ON trip_with_system_pk.route_pk = route.pk
ORDER BY distance.val
LIMIT sqlc.arg(num_vehicles);

-- name: GetVehicle :one
WITH trip_with_system_pk AS (
  SELECT trip.*, feed.system_pk
  FROM trip
  LEFT JOIN feed ON trip.feed_pk = feed.pk
)
SELECT vehicle.*,
       stop.id as stop_id,
       stop.name as stop_name,
       trip_with_system_pk.id as trip_id,
       trip_with_system_pk.direction_id as trip_direction_id,
       route.id as route_id,
       route.color as route_color
FROM vehicle
LEFT JOIN stop ON vehicle.current_stop_pk = stop.pk
LEFT JOIN trip_with_system_pk ON vehicle.trip_id = trip_with_system_pk.id AND vehicle.system_pk = trip_with_system_pk.system_pk
LEFT JOIN route ON trip_with_system_pk.route_pk = route.pk
WHERE vehicle.system_pk = sqlc.arg(system_pk) AND vehicle.id = sqlc.arg(vehicle_id);

-- name: ListVehicleIdPkTripPkGtfsHash :many
SELECT id, trip_id, pk FROM vehicle
WHERE id = ANY(sqlc.arg(vehicle_ids)::text[])
AND system_pk = sqlc.arg(system_pk);

-- name: DeleteStaleVehicles :exec
DELETE FROM vehicle
WHERE
  feed_pk = sqlc.arg(feed_pk)
  AND NOT id = ANY(sqlc.arg(active_vehicle_ids)::text[]);

-- name: ClearVehicleFeed :exec
DELETE FROM vehicle
WHERE feed_pk = sqlc.arg(feed_pk);
