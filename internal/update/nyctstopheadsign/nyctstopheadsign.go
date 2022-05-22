// Package nyctstopheadsign contains logic for updating the stop headsign rules from the NYCT CSV file.
package nyctstopheadsign

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"

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
	m := map[string]int{}
	for i, colHeader := range records[0] {
		m[colHeader] = i
	}
	return nil
}
