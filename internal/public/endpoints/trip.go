package endpoints

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jamespfennell/transiter/internal/convert"
	"github.com/jamespfennell/transiter/internal/gen/api"
	"github.com/jamespfennell/transiter/internal/gen/db"
	"github.com/jamespfennell/transiter/internal/public/errors"
)

func ListTrips(ctx context.Context, r *Context, req *api.ListTripsRequest) (*api.ListTripsReply, error) {
	system, route, err := getRoute(ctx, r.Querier, req.SystemId, req.RouteId)
	if err != nil {
		return nil, err
	}
	trips, err := r.Querier.ListTrips(ctx, route.Pk)
	if err != nil {
		return nil, err
	}
	apiTrips, err := buildApiTrips(ctx, r, &system, &route, trips)
	if err != nil {
		return nil, err
	}
	return &api.ListTripsReply{
		Trips: apiTrips,
	}, nil
}

func GetTrip(ctx context.Context, r *Context, req *api.GetTripRequest) (*api.Trip, error) {
	system, route, err := getRoute(ctx, r.Querier, req.SystemId, req.RouteId)
	if err != nil {
		return nil, err
	}
	trip, err := r.Querier.GetTrip(ctx, db.GetTripParams{
		TripID:  req.TripId,
		RoutePk: route.Pk,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.NewNotFoundError(fmt.Sprintf("trip %q in route %q in system %q not found",
				req.TripId, req.RouteId, req.SystemId))
		}
		return nil, err
	}
	apiTrips, err := buildApiTrips(ctx, r, &system, &route, []db.Trip{trip})
	if err != nil {
		return nil, err
	}
	return apiTrips[0], nil
}

func buildApiTrips(ctx context.Context, r *Context, system *db.System, route *db.Route, trips []db.Trip) ([]*api.Trip, error) {
	var apiTrips []*api.Trip
	for i := range trips {
		trip := &trips[i]
		stopTimes, err := r.Querier.ListStopsTimesForTrip(ctx, trip.Pk)
		if err != nil {
			return nil, err
		}
		reply := &api.Trip{
			Id:          trip.ID,
			DirectionId: trip.DirectionID.Bool,
			StartedAt:   convert.SQLNullTime(trip.StartedAt),
			Route:       r.Reference.Route(route.ID, system.ID, route.Color),
		}
		// TODO: vechices
		// if trip.VehicleID.Valid {
		//	reply.Vehicle = r.Reference.Vehicle(trip.VehicleID.String)
		//}
		for _, stopTime := range stopTimes {
			reply.StopTimes = append(reply.StopTimes, &api.StopTime{
				StopSequence: stopTime.StopSequence,
				Track:        convert.SQLNullString(stopTime.Track),
				Future:       !stopTime.Past,
				Arrival:      buildEstimatedTime(stopTime.ArrivalTime, stopTime.ArrivalDelay, stopTime.ArrivalUncertainty),
				Departure:    buildEstimatedTime(stopTime.DepartureTime, stopTime.DepartureDelay, stopTime.DepartureUncertainty),
				Stop:         r.Reference.Stop(stopTime.StopID, system.ID, stopTime.StopName),
			})
		}
		apiTrips = append(apiTrips, reply)
	}
	return apiTrips, nil
}
