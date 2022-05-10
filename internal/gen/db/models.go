// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package db

import (
	"database/sql"
	"time"

	"github.com/jackc/pgtype"
)

type Agency struct {
	Pk       int64
	ID       string
	SystemPk int64
	SourcePk int64
	Name     string
	Url      string
	Timezone string
	Language sql.NullString
	Phone    sql.NullString
	FareUrl  sql.NullString
	Email    sql.NullString
}

type Alert struct {
	Pk        int64
	ID        string
	SourcePk  int64
	SystemPk  int64
	Cause     string
	Effect    string
	CreatedAt sql.NullTime
	SortOrder sql.NullInt32
	UpdatedAt sql.NullTime
}

type AlertActivePeriod struct {
	Pk       int64
	AlertPk  int64
	StartsAt sql.NullTime
	EndsAt   sql.NullTime
}

type AlertAgency struct {
	AlertPk  int64
	AgencyPk int64
}

type AlertMessage struct {
	Pk          int64
	AlertPk     int64
	Header      string
	Description string
	Url         sql.NullString
	Language    sql.NullString
}

type AlertRoute struct {
	AlertPk int64
	RoutePk int64
}

type AlertStop struct {
	AlertPk int64
	StopPk  int64
}

type AlertTrip struct {
	AlertPk int64
	TripPk  int64
}

type DirectionNameRule struct {
	Pk          int64
	ID          sql.NullString
	StopPk      int64
	SourcePk    int64
	Priority    int32
	DirectionID sql.NullBool
	Track       sql.NullString
	Name        string
}

type Feed struct {
	Pk                int64
	ID                string
	SystemPk          int64
	AutoUpdateEnabled bool
	AutoUpdatePeriod  sql.NullInt32
	Config            string
}

type FeedUpdate struct {
	Pk               int64
	FeedPk           int64
	Status           string
	CreatedAt        sql.NullTime
	CompletedAt      sql.NullTime
	ContentCreatedAt sql.NullTime
	ContentHash      sql.NullString
	ContentLength    sql.NullInt32
	Result           sql.NullString
	ResultMessage    sql.NullString
	TotalDuration    sql.NullInt32
}

type Route struct {
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
}

type ScheduledService struct {
	Pk        int64
	ID        string
	SystemPk  int64
	SourcePk  int64
	Monday    sql.NullBool
	Tuesday   sql.NullBool
	Wednesday sql.NullBool
	Thursday  sql.NullBool
	Friday    sql.NullBool
	Saturday  sql.NullBool
	Sunday    sql.NullBool
	EndDate   sql.NullTime
	StartDate sql.NullTime
}

type ScheduledServiceAddition struct {
	Pk        int64
	ServicePk int64
	Date      time.Time
}

type ScheduledServiceRemoval struct {
	Pk        int64
	ServicePk int64
	Date      time.Time
}

type ScheduledTrip struct {
	Pk                   int64
	ID                   string
	RoutePk              int64
	ServicePk            int64
	DirectionID          sql.NullBool
	BikesAllowed         string
	BlockID              sql.NullString
	Headsign             sql.NullString
	ShortName            sql.NullString
	WheelchairAccessible string
}

type ScheduledTripFrequency struct {
	Pk             int64
	TripPk         int64
	StartTime      time.Time
	EndTime        time.Time
	Headway        int32
	FrequencyBased bool
}

type ScheduledTripStopTime struct {
	Pk                    int64
	TripPk                int64
	StopPk                int64
	ArrivalTime           sql.NullTime
	DepartureTime         sql.NullTime
	StopSequence          int32
	ContinuousDropOff     string
	ContinuousPickup      string
	DropOffType           string
	ExactTimes            bool
	Headsign              sql.NullString
	PickupType            string
	ShapeDistanceTraveled sql.NullFloat64
}

type ServiceMap struct {
	Pk       int64
	RoutePk  int64
	ConfigPk int64
}

type ServiceMapConfig struct {
	Pk                     int64
	ID                     string
	SystemPk               int64
	Config                 []byte
	DefaultForRoutesAtStop bool
	DefaultForStopsInRoute bool
}

type ServiceMapVertex struct {
	Pk       int64
	StopPk   int64
	MapPk    int64
	Position int32
}

type Stop struct {
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
}

type System struct {
	Pk       int64
	ID       string
	Name     string
	Timezone sql.NullString
	Status   string
}

type SystemUpdate struct {
	Pk               int64
	SystemPk         int64
	Status           string
	StatusMessage    sql.NullString
	TotalDuration    sql.NullFloat64
	ScheduledAt      sql.NullTime
	CompletedAt      sql.NullTime
	Config           sql.NullString
	ConfigTemplate   sql.NullString
	ConfigParameters sql.NullString
	ConfigSourceUrl  sql.NullString
	TransiterVersion sql.NullString
}

type Transfer struct {
	Pk              int64
	SourcePk        sql.NullInt64
	ConfigSourcePk  sql.NullInt64
	SystemPk        sql.NullInt64
	FromStopPk      int64
	ToStopPk        int64
	Type            string
	MinTransferTime sql.NullInt32
	Distance        sql.NullInt32
}

type TransfersConfig struct {
	Pk       int64
	Distance pgtype.Numeric
}

type TransfersConfigSystem struct {
	TransfersConfigPk sql.NullInt64
	SystemPk          sql.NullInt64
}

type Trip struct {
	Pk                  int64
	ID                  string
	RoutePk             int64
	SourcePk            int64
	DirectionID         sql.NullBool
	Delay               sql.NullInt32
	StartedAt           sql.NullTime
	UpdatedAt           sql.NullTime
	CurrentStopSequence sql.NullInt32
}

type TripStopTime struct {
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
}

type Vehicle struct {
	Pk                  int64
	ID                  sql.NullString
	SourcePk            int64
	SystemPk            int64
	TripPk              sql.NullInt64
	Label               sql.NullString
	LicensePlate        sql.NullString
	CurrentStatus       string
	Latitude            sql.NullFloat64
	Longitude           sql.NullFloat64
	Bearing             sql.NullFloat64
	Odometer            sql.NullFloat64
	Speed               sql.NullFloat64
	CongestionLevel     string
	UpdatedAt           sql.NullTime
	CurrentStopPk       sql.NullInt64
	CurrentStopSequence sql.NullInt32
	OccupancyStatus     string
}
