// Package static contains the code for updating the database from a GTFS static feed.
package static

import (
	"context"
	"fmt"

	"github.com/jamespfennell/gtfs"
	"github.com/jamespfennell/transiter/internal/convert"
	"github.com/jamespfennell/transiter/internal/db/dbwrappers"
	"github.com/jamespfennell/transiter/internal/gen/db"
	"github.com/jamespfennell/transiter/internal/servicemaps"
	"github.com/jamespfennell/transiter/internal/update/common"
)

func Parse(content []byte) (*gtfs.Static, error) {
	// TODO: support custom GTFS static options
	return gtfs.ParseStatic(content, gtfs.ParseStaticOptions{})
}

func Update(ctx context.Context, updateCtx common.UpdateContext, data *gtfs.Static) error {
	agencyIDToPk, err := updateAgencies(ctx, updateCtx, data.Agencies)
	if err != nil {
		return err
	}
	routeIDToPk, err := updateRoutes(ctx, updateCtx, data.Routes, agencyIDToPk)
	if err != nil {
		return err
	}
	stopIDToPk, err := updateStops(ctx, updateCtx, data.Stops)
	if err != nil {
		return err
	}
	if err := updateTransfers(ctx, updateCtx, data.Transfers, stopIDToPk); err != nil {
		return err
	}
	if err := servicemaps.UpdateStaticMaps(ctx, updateCtx.Querier, updateCtx.Logger, servicemaps.UpdateStaticMapsArgs{
		SystemPk:    updateCtx.SystemPk,
		Trips:       data.Trips,
		RouteIDToPk: routeIDToPk,
	}); err != nil {
		return err
	}
	return nil
}

func updateAgencies(ctx context.Context, updateCtx common.UpdateContext, agencies []gtfs.Agency) (map[string]int64, error) {
	idToPk, err := dbwrappers.MapAgencyIDToPk(ctx, updateCtx.Querier, updateCtx.SystemPk)
	if err != nil {
		return nil, err
	}
	for _, agency := range agencies {
		var err error
		pk, ok := idToPk[agency.Id]
		if ok {
			err = updateCtx.Querier.UpdateAgency(ctx, db.UpdateAgencyParams{
				Pk:       pk,
				SourcePk: updateCtx.UpdatePk,
				Name:     agency.Name,
				Url:      agency.Url,
				Timezone: agency.Timezone,
				Language: convert.NullString(agency.Language),
				Phone:    convert.NullString(agency.Phone),
				FareUrl:  convert.NullString(agency.FareUrl),
				Email:    convert.NullString(agency.Email),
			})
		} else {
			pk, err = updateCtx.Querier.InsertAgency(ctx, db.InsertAgencyParams{
				ID:       agency.Id,
				SystemPk: updateCtx.SystemPk,
				SourcePk: updateCtx.UpdatePk,
				Name:     agency.Name,
				Url:      agency.Url,
				Timezone: agency.Timezone,
				Language: convert.NullString(agency.Language),
				Phone:    convert.NullString(agency.Phone),
				FareUrl:  convert.NullString(agency.FareUrl),
				Email:    convert.NullString(agency.Email),
			})
			idToPk[agency.Id] = pk
		}
		if err != nil {
			return nil, err
		}
	}
	deletedIds, err := updateCtx.Querier.DeleteStaleAgencies(ctx, db.DeleteStaleAgenciesParams{
		FeedPk:   updateCtx.FeedPk,
		UpdatePk: updateCtx.UpdatePk,
	})
	if err != nil {
		return nil, err
	}
	for _, id := range deletedIds {
		delete(idToPk, id)
	}
	return idToPk, nil
}

func updateRoutes(ctx context.Context, updateCtx common.UpdateContext, routes []gtfs.Route, agencyIDToPk map[string]int64) (map[string]int64, error) {
	idToPk, err := dbwrappers.MapRouteIDToPkInSystem(ctx, updateCtx.Querier, updateCtx.SystemPk)
	if err != nil {
		return nil, err
	}
	for _, route := range routes {
		agencyPk, ok := agencyIDToPk[route.Agency.Id]
		if !ok {
			updateCtx.Logger.WarnCtx(ctx, fmt.Sprintf("no agency %q in the database; skipping route %q", route.Agency.Id, route.Id))
			continue
		}
		pk, ok := idToPk[route.Id]
		if ok {
			err = updateCtx.Querier.UpdateRoute(ctx, db.UpdateRouteParams{
				Pk:                pk,
				SourcePk:          updateCtx.UpdatePk,
				Color:             route.Color,
				TextColor:         route.TextColor,
				ShortName:         convert.NullString(route.ShortName),
				LongName:          convert.NullString(route.LongName),
				Description:       convert.NullString(route.Description),
				Url:               convert.NullString(route.Url),
				SortOrder:         convert.NullInt32(route.SortOrder),
				Type:              route.Type.String(),
				ContinuousPickup:  route.ContinuousPickup.String(),
				ContinuousDropOff: route.ContinuousDropOff.String(),
				AgencyPk:          agencyPk,
			})
		} else {
			pk, err = updateCtx.Querier.InsertRoute(ctx, db.InsertRouteParams{
				ID:                route.Id,
				SystemPk:          updateCtx.SystemPk,
				SourcePk:          updateCtx.UpdatePk,
				Color:             route.Color,
				TextColor:         route.TextColor,
				ShortName:         convert.NullString(route.ShortName),
				LongName:          convert.NullString(route.LongName),
				Description:       convert.NullString(route.Description),
				Url:               convert.NullString(route.Url),
				SortOrder:         convert.NullInt32(route.SortOrder),
				Type:              route.Type.String(),
				ContinuousPickup:  route.ContinuousPickup.String(),
				ContinuousDropOff: route.ContinuousDropOff.String(),
				AgencyPk:          agencyPk,
			})
			idToPk[route.Id] = pk
		}
		if err != nil {
			return nil, err
		}
	}
	deletedIds, err := updateCtx.Querier.DeleteStaleRoutes(ctx, db.DeleteStaleRoutesParams{
		FeedPk:   updateCtx.FeedPk,
		UpdatePk: updateCtx.UpdatePk,
	})
	if err != nil {
		return nil, err
	}
	for _, id := range deletedIds {
		delete(idToPk, id)
	}
	return idToPk, nil
}

func updateStops(ctx context.Context, updateCtx common.UpdateContext, stops []gtfs.Stop) (map[string]int64, error) {
	idToPk, err := dbwrappers.MapStopIDToPkInSystem(ctx, updateCtx.Querier, updateCtx.SystemPk)
	if err != nil {
		return nil, err
	}
	for _, stop := range stops {
		var wheelchairBoarding *bool
		switch stop.WheelchairBoarding {
		case gtfs.Possible:
			t := true
			wheelchairBoarding = &t
		case gtfs.NotPossible:
			f := false
			wheelchairBoarding = &f
		}
		pk, ok := idToPk[stop.Id]
		if ok {
			err = updateCtx.Querier.UpdateStop(ctx, db.UpdateStopParams{
				Pk:                 pk,
				SourcePk:           updateCtx.UpdatePk,
				Name:               convert.NullString(stop.Name),
				Type:               stop.Type.String(),
				Longitude:          convert.Gps(stop.Longitude),
				Latitude:           convert.Gps(stop.Latitude),
				Url:                convert.NullString(stop.Url),
				Code:               convert.NullString(stop.Code),
				Description:        convert.NullString(stop.Description),
				PlatformCode:       convert.NullString(stop.PlatformCode),
				Timezone:           convert.NullString(stop.Timezone),
				WheelchairBoarding: convert.NullBool(wheelchairBoarding),
				ZoneID:             convert.NullString(stop.ZoneId),
			})
		} else {
			pk, err = updateCtx.Querier.InsertStop(ctx, db.InsertStopParams{
				ID:                 stop.Id,
				SystemPk:           updateCtx.SystemPk,
				SourcePk:           updateCtx.UpdatePk,
				Name:               convert.NullString(stop.Name),
				Type:               stop.Type.String(),
				Longitude:          convert.Gps(stop.Longitude),
				Latitude:           convert.Gps(stop.Latitude),
				Url:                convert.NullString(stop.Url),
				Code:               convert.NullString(stop.Code),
				Description:        convert.NullString(stop.Description),
				PlatformCode:       convert.NullString(stop.PlatformCode),
				Timezone:           convert.NullString(stop.Timezone),
				WheelchairBoarding: convert.NullBool(wheelchairBoarding),
				ZoneID:             convert.NullString(stop.ZoneId),
			})
			idToPk[stop.Id] = pk
		}
		if err != nil {
			return nil, err
		}
	}
	deletedIds, err := updateCtx.Querier.DeleteStaleStops(ctx, db.DeleteStaleStopsParams{
		FeedPk:   updateCtx.FeedPk,
		UpdatePk: updateCtx.UpdatePk,
	})
	if err != nil {
		return nil, err
	}
	for _, id := range deletedIds {
		delete(idToPk, id)
	}
	// We now populate the parent stop field
	for _, stop := range stops {
		if stop.Parent == nil {
			continue
		}
		parentStopPk, ok := idToPk[stop.Parent.Id]
		if !ok {
			continue
		}
		if err := updateCtx.Querier.UpdateStop_Parent(ctx, db.UpdateStop_ParentParams{
			Pk:           idToPk[stop.Id],
			ParentStopPk: convert.NullInt64(&parentStopPk),
		}); err != nil {
			return nil, err
		}
	}
	return idToPk, nil
}

func updateTransfers(ctx context.Context, updateCtx common.UpdateContext, transfers []gtfs.Transfer, stopIDToPk map[string]int64) error {
	if err := updateCtx.Querier.DeleteStaleTransfers(ctx, db.DeleteStaleTransfersParams{
		FeedPk:   updateCtx.FeedPk,
		UpdatePk: updateCtx.UpdatePk,
	}); err != nil {
		return err
	}
	for _, transfer := range transfers {
		fromPk, ok := stopIDToPk[transfer.From.Id]
		if !ok {
			continue
		}
		toPk, ok := stopIDToPk[transfer.To.Id]
		if !ok {
			continue
		}
		if err := updateCtx.Querier.InsertTransfer(ctx, db.InsertTransferParams{
			SystemPk:        convert.NullInt64(&updateCtx.SystemPk),
			SourcePk:        convert.NullInt64(&updateCtx.UpdatePk),
			FromStopPk:      fromPk,
			ToStopPk:        toPk,
			Type:            transfer.Type.String(),
			MinTransferTime: convert.NullInt32(transfer.MinTransferTime),
		}); err != nil {
			return err
		}
	}
	return nil
}
