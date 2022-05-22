// Package nyctstopheadsign contains logic for updating the stop headsign rules from the NYCT CSV file.
package nyctstopheadsign

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"

	"github.com/jamespfennell/gtfs"
	"github.com/jamespfennell/transiter/internal/convert"
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
	// TODO prepend the custom rules
	var rules []rule
	for _, row := range records[1:] {
		rules = append(rules, rule{
			stopID:      row[stopIDCol] + "N",
			directionID: gtfs.DirectionIDFalse,
			headsign:    row[northHeadsignCol],
		})
		rules = append(rules, rule{
			stopID:      row[stopIDCol] + "S",
			directionID: gtfs.DirectionIDTrue,
			headsign:    row[southHeadsignCol],
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
			SourcePk:    updateCtx.UpdatePk,
			Priority:    int32(i),
			StopPk:      stopPk,
			DirectionID: convert.DirectionID(rule.directionID),
			Headsign:    rule.headsign,
		}); err != nil {
			return err
		}
	}
	return nil
}

type rule struct {
	stopID      string
	directionID gtfs.DirectionID
	headsign    string
}

// TODO make the headsign nullable?
/*
def _clean_mta_name(mta_name):
if mta_name.strip() == "":
	return "(Terminating trains)"
return mta_name.strip().replace("&", "and")
*/
