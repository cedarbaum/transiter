// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	CalculatePeriodicityForRoute(ctx context.Context, routePk int64) (interface{}, error)
	CountAgenciesInSystem(ctx context.Context, systemPk int64) (int64, error)
	CountFeedsInSystem(ctx context.Context, systemPk int64) (int64, error)
	CountRoutesInSystem(ctx context.Context, systemPk int64) (int64, error)
	CountStopsInSystem(ctx context.Context, systemPk int64) (int64, error)
	CountSystems(ctx context.Context) (int64, error)
	CountTransfersInSystem(ctx context.Context, systemPk sql.NullInt64) (int64, error)
	DeleteAgency(ctx context.Context, pk int64) error
	DeleteFeed(ctx context.Context, pk int64) error
	DeleteSystem(ctx context.Context, pk int64) error
	GetAgencyInSystem(ctx context.Context, arg GetAgencyInSystemParams) (Agency, error)
	GetFeedForUpdate(ctx context.Context, updatePk int64) (Feed, error)
	GetFeedInSystem(ctx context.Context, arg GetFeedInSystemParams) (Feed, error)
	GetLastStopsForTrips(ctx context.Context, tripPks []int64) ([]GetLastStopsForTripsRow, error)
	GetRouteInSystem(ctx context.Context, arg GetRouteInSystemParams) (GetRouteInSystemRow, error)
	GetStopInSystem(ctx context.Context, arg GetStopInSystemParams) (GetStopInSystemRow, error)
	GetSystem(ctx context.Context, id string) (System, error)
	GetTrip(ctx context.Context, arg GetTripParams) (GetTripRow, error)
	InsertAgency(ctx context.Context, arg InsertAgencyParams) error
	InsertFeed(ctx context.Context, arg InsertFeedParams) error
	InsertFeedUpdate(ctx context.Context, arg InsertFeedUpdateParams) (int64, error)
	InsertSystem(ctx context.Context, arg InsertSystemParams) error
	ListActiveAlertsForAgency(ctx context.Context, arg ListActiveAlertsForAgencyParams) ([]ListActiveAlertsForAgencyRow, error)
	ListActiveAlertsForRoutes(ctx context.Context, arg ListActiveAlertsForRoutesParams) ([]ListActiveAlertsForRoutesRow, error)
	ListActiveAlertsForStops(ctx context.Context, arg ListActiveAlertsForStopsParams) ([]ListActiveAlertsForStopsRow, error)
	ListAgenciesInSystem(ctx context.Context, systemPk int64) ([]Agency, error)
	ListAutoUpdateFeedsForSystem(ctx context.Context, systemID string) ([]ListAutoUpdateFeedsForSystemRow, error)
	ListDirectionNameRulesForStops(ctx context.Context, stopPks []int64) ([]DirectionNameRule, error)
	ListFeedsInSystem(ctx context.Context, systemPk int64) ([]Feed, error)
	ListMessagesForAlerts(ctx context.Context, alertPks []int64) ([]AlertMessage, error)
	ListRoutesByPk(ctx context.Context, routePks []int64) ([]Route, error)
	ListRoutesInAgency(ctx context.Context, agencyPk int64) ([]ListRoutesInAgencyRow, error)
	ListRoutesInSystem(ctx context.Context, systemPk int64) ([]Route, error)
	ListServiceMapsForRoute(ctx context.Context, routePk int64) ([]ListServiceMapsForRouteRow, error)
	ListServiceMapsForStops(ctx context.Context, stopPks []int64) ([]ListServiceMapsForStopsRow, error)
	ListServiceMapsGroupIDsForStops(ctx context.Context, stopPks []int64) ([]ListServiceMapsGroupIDsForStopsRow, error)
	ListStopTimesAtStops(ctx context.Context, stopPks []int64) ([]ListStopTimesAtStopsRow, error)
	ListStopsInStopTree(ctx context.Context, pk int64) ([]Stop, error)
	ListStopsInSystem(ctx context.Context, systemPk int64) ([]Stop, error)
	ListStopsTimesForTrip(ctx context.Context, tripPk int64) ([]ListStopsTimesForTripRow, error)
	ListSystems(ctx context.Context) ([]System, error)
	ListTransfersFromStops(ctx context.Context, fromStopPks []int64) ([]ListTransfersFromStopsRow, error)
	ListTransfersInSystem(ctx context.Context, systemPk sql.NullInt64) ([]ListTransfersInSystemRow, error)
	ListTripsInRoute(ctx context.Context, routePk int64) ([]ListTripsInRouteRow, error)
	ListUpdatesInFeed(ctx context.Context, feedPk int64) ([]FeedUpdate, error)
	MapAgencyPkToIdInSystem(ctx context.Context, systemPk int64) ([]MapAgencyPkToIdInSystemRow, error)
	UpdateAgency(ctx context.Context, arg UpdateAgencyParams) error
	UpdateFeed(ctx context.Context, arg UpdateFeedParams) error
	UpdateSystem(ctx context.Context, arg UpdateSystemParams) error
}

var _ Querier = (*Queries)(nil)
