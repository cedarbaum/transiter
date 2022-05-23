// Package nyctsubwaycsv contains logic for updating the stop headsign rules from the NYCT CSV file.
package nyctsubwaycsv

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"

	"github.com/jamespfennell/transiter/internal/db/dbwrappers"
	"github.com/jamespfennell/transiter/internal/gen/db"
	"github.com/jamespfennell/transiter/internal/update/common"
)

func ParseAndUpdate(ctx context.Context, updateCtx common.UpdateContext, content []byte) error {
	csvReader := csv.NewReader(bytes.NewReader(content))
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return fmt.Errorf("file contains no header row")
	}
	stopIDCol := -1
	northHeadsignCol := -1
	southHeadsignCol := -1
	for i, header := range records[0] {
		switch header {
		case "GTFS Stop ID":
			stopIDCol = i
		case "North Direction Label":
			northHeadsignCol = i
		case "South Direction Label":
			southHeadsignCol = i
		}
	}
	if stopIDCol < 0 {
		return fmt.Errorf("CSV file missing stop ID column")
	}
	if northHeadsignCol < 0 {
		return fmt.Errorf("CSV file missing north headsign/label column")
	}
	if southHeadsignCol < 0 {
		return fmt.Errorf("CSV file missing south headsign/label column")
	}
	rules := customRules()
	for _, row := range records[1:] {
		rules = append(rules, rule{
			stopID:   row[stopIDCol] + "N",
			headsign: row[northHeadsignCol],
		})
		rules = append(rules, rule{
			stopID:   row[stopIDCol] + "S",
			headsign: row[southHeadsignCol],
		})
	}
	if err := updateCtx.Querier.DeleteStopHeadsignRules(ctx, updateCtx.FeedPk); err != nil {
		return err
	}
	stopIDsSet := map[string]bool{}
	var stopIDs []string
	for _, rule := range rules {
		if stopIDsSet[rule.stopID] {
			continue
		}
		stopIDsSet[rule.stopID] = true
		stopIDs = append(stopIDs, rule.stopID)
	}
	stopIDToPk, err := dbwrappers.MapStopIDToPkInSystem(ctx, updateCtx.Querier, updateCtx.SystemPk, stopIDs)
	if err != nil {
		return err
	}
	for i, rule := range rules {
		stopPk, ok := stopIDToPk[rule.stopID]
		if !ok {
			continue
		}
		if err := updateCtx.Querier.InsertStopHeadSignRule(ctx, db.InsertStopHeadSignRuleParams{
			SourcePk: updateCtx.UpdatePk,
			Priority: int32(i),
			StopPk:   stopPk,
			Headsign: rule.headsign,
		}); err != nil {
			return err
		}
	}
	return nil
}

type rule struct {
	stopID   string
	track    *string
	headsign string
}

// TODO make the headsign nullable?
/*
def _clean_mta_name(mta_name):
if mta_name.strip() == "":
	return "(Terminating trains)"
return mta_name.strip().replace("&", "and")
*/
const (
	eastSideAndQueens = "East Side and Queens"
	manhattan         = "Manhattan"
	rockaways         = "Euclid - Lefferts - Rockaways" // To be consistent with the MTA
	uptown            = "Uptown"
	uptownAndTheBronx = "Uptown and The Bronx"
	queens            = "Queens"
)

func customRules() []rule {
	optOf := func(s string) *string {
		return &s
	}
	return []rule{
		// Hoyt-Schermerhorn Sts station
		{
			stopID:   "A42N",
			track:    optOf("E2"),
			headsign: "Court Sq, Queens",
		},
		{
			stopID:   "A42N",
			headsign: manhattan,
		},
	}
}

/*
	  //,{MANHATTAN}
SPECIAL_STOPS_CSV = f"""
stop_id,track,track_name,basic_name
A42S,E1,"Church Av, Brooklyn",{ROCKAWAYS}
A41S,B1,Coney Island,{ROCKAWAYS}
A25N,D4,{EAST_SIDE_AND_QUEENS},{UPTOWN_AND_THE_BRONX}
D15N,B2,{EAST_SIDE_AND_QUEENS},{UPTOWN_AND_THE_BRONX}
R14N,A4,{UPTOWN},{QUEENS}
B08N,T2,{QUEENS},{UPTOWN}
D14N,D4,{EAST_SIDE_AND_QUEENS},{UPTOWN_AND_THE_BRONX}
D26N,A2,Franklin Avenue,{MANHATTAN}
"""
*/
