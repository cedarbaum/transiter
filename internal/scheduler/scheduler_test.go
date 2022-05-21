package scheduler

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
)

const systemID1 = "systemID1"
const feedId1 = "feedId1"
const feedId2 = "feedId2"
const systemID2 = "systemID2"
const feedId3 = "feedId3"

func TestScheduler(t *testing.T) {
	resetSystem1 := func(s *Scheduler) error {
		return s.Reset(context.Background(), systemID1)
	}
	resetSystem2 := func(s *Scheduler) error {
		return s.Reset(context.Background(), systemID2)
	}
	resetAll := func(s *Scheduler) error {
		return s.ResetAll(context.Background())
	}
	testCases := []struct {
		description     string
		update          []SystemConfig
		resetF          func(*Scheduler) error
		runningPeriod   time.Duration
		expectedUpdates map[systemAndFeed]int
	}{
		{
			description: "just change periodicity of one feed",
			update: []SystemConfig{
				{
					Id: systemID1,
					FeedConfigs: []FeedConfig{
						{
							Id:     feedId1,
							Period: 1000 * time.Millisecond,
						},
					},
				},
			},
			resetF:        resetSystem1,
			runningPeriod: 2000 * time.Millisecond,
			expectedUpdates: map[systemAndFeed]int{
				{systemID: systemID1, feedId: feedId1}: 2,
			},
		},
		{
			description: "new feed in same system",
			update: []SystemConfig{
				{
					Id: systemID1,
					FeedConfigs: []FeedConfig{
						{
							Id:     feedId1,
							Period: 500 * time.Millisecond,
						},
						{
							Id:     feedId2,
							Period: 500 * time.Millisecond,
						},
					},
				},
			},
			resetF:        resetSystem1,
			runningPeriod: 2000 * time.Millisecond,
			expectedUpdates: map[systemAndFeed]int{
				{systemID: systemID1, feedId: feedId1}: 4,
				{systemID: systemID1, feedId: feedId2}: 4,
			},
		},
		{
			description: "remove feed in system",
			update: []SystemConfig{
				{
					Id:          systemID1,
					FeedConfigs: []FeedConfig{},
				},
			},
			resetF:          resetSystem1,
			runningPeriod:   2000 * time.Millisecond,
			expectedUpdates: map[systemAndFeed]int{},
		},
		{
			description:     "remove system",
			update:          []SystemConfig{},
			resetF:          resetSystem1,
			runningPeriod:   2000 * time.Millisecond,
			expectedUpdates: map[systemAndFeed]int{},
		},
		{
			description: "new system, only reset the new one",
			update: []SystemConfig{
				{
					Id: systemID1,
					FeedConfigs: []FeedConfig{
						{
							Id:     feedId1,
							Period: 500 * time.Millisecond,
						},
					},
				},
				{
					Id: systemID2,
					FeedConfigs: []FeedConfig{
						{
							Id:     feedId3,
							Period: 500 * time.Millisecond,
						},
					},
				},
			},
			resetF:        resetSystem2,
			runningPeriod: 2000 * time.Millisecond,
			expectedUpdates: map[systemAndFeed]int{
				{systemID: systemID1, feedId: feedId1}: 4,
				{systemID: systemID2, feedId: feedId3}: 4,
			},
		},
		{
			description: "new system, reset all",
			update: []SystemConfig{
				{
					Id: systemID1,
					FeedConfigs: []FeedConfig{
						{
							Id:     feedId1,
							Period: 500 * time.Millisecond,
						},
					},
				},
				{
					Id: systemID2,
					FeedConfigs: []FeedConfig{
						{
							Id:     feedId3,
							Period: 500 * time.Millisecond,
						},
					},
				},
			},
			resetF:        resetAll,
			runningPeriod: 2000 * time.Millisecond,
			expectedUpdates: map[systemAndFeed]int{
				{systemID: systemID1, feedId: feedId1}: 4,
				{systemID: systemID2, feedId: feedId3}: 4,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			ctx, cancelFunc := context.WithCancel(context.Background())
			defer cancelFunc()

			ops := testOps{
				currentConfig: []SystemConfig{
					{
						Id: systemID1,
						FeedConfigs: []FeedConfig{
							{
								Id:     feedId1,
								Period: 500 * time.Millisecond,
							},
						},
					},
				},
				updateChan: make(chan systemAndFeed, 20),
			}
			clock := clock.NewMock()

			scheduler := New()
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				scheduler.RunWithClockAndOps(ctx, clock, &ops)
				wg.Done()
			}()
			if err := scheduler.ResetAll(ctx); err != nil {
				t.Fatalf("failed to create scheduler: %s", err)
			}

			clock.Add(2000 * time.Millisecond)

			updates := ops.getUpdates(4)
			expected := map[systemAndFeed]int{
				{systemID: systemID1, feedId: feedId1}: 4,
			}
			if !reflect.DeepEqual(expected, updates) {
				t.Errorf("Updates got = %+v, want = %+v", updates, expected)
			}

			ops.currentConfig = tc.update

			if err := tc.resetF(scheduler); err != nil {
				t.Errorf("Unexpected error when reseting: %s", err)
			}

			clock.Add(tc.runningPeriod)

			var numExpectedUpdates int
			for _, val := range tc.expectedUpdates {
				numExpectedUpdates += val
			}
			updates = ops.getUpdates(numExpectedUpdates)
			if !reflect.DeepEqual(tc.expectedUpdates, updates) {
				t.Errorf("Updates got = %+v, want = %+v", updates, tc.expectedUpdates)
			}

			cancelFunc()
			wg.Wait()
			close(ops.updateChan)
			for update := range ops.updateChan {
				t.Errorf("Unexpected update after scheduler stopped: %s", update)
			}
		})
	}
}

type testOps struct {
	currentConfig []SystemConfig
	updateChan    chan systemAndFeed
}

func (ops *testOps) ListSystemConfigs(ctx context.Context) ([]SystemConfig, error) {
	return ops.currentConfig, nil
}

func (ops *testOps) GetSystemConfig(ctx context.Context, systemID string) (SystemConfig, error) {
	for _, config := range ops.currentConfig {
		if config.Id == systemID {
			return config, nil
		}
	}
	return SystemConfig{Id: systemID}, nil
}

func (ops *testOps) UpdateFeed(ctx context.Context, systemID, feedId string) error {
	ops.updateChan <- systemAndFeed{systemID: systemID, feedId: feedId}
	return nil
}

func (ops *testOps) getUpdates(num int) map[systemAndFeed]int {
	if num == 0 {
		return map[systemAndFeed]int{}
	}
	timeoutChan := time.After(time.Duration(num+1) * 500 * time.Millisecond)
	m := map[systemAndFeed]int{}
	seen := 0
	for {
		select {
		case key := <-ops.updateChan:
			m[key] = m[key] + 1
			seen += 1
			if seen == num {
				return m
			}
		case <-timeoutChan:
			return m
		}
	}
}

type systemAndFeed struct {
	systemID, feedId string
}
