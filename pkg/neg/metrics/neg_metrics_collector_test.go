package metrics

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"k8s.io/ingress-gce/pkg/neg/types"
	"k8s.io/utils/clock"
)

func TestUpdateMigrationStartAndEndTime(t *testing.T) {
	syncerKey := syncerKey(1)
	testCases := []struct {
		name                        string
		migrationCount              int
		dualStackMigrationStartTime map[types.NegSyncerKey]time.Time
		dualStackMigrationEndTime   map[types.NegSyncerKey]time.Time
		wantChangedStartTime        bool
		wantChangedEndTime          bool
	}{
		{
			name:                        "start time should get set since this is the first time when migration starts",
			migrationCount:              5,
			dualStackMigrationStartTime: map[types.NegSyncerKey]time.Time{},
			dualStackMigrationEndTime:   map[types.NegSyncerKey]time.Time{},
			wantChangedStartTime:        true,
		},
		{
			name:           "start time is already present so no action required",
			migrationCount: 5,
			dualStackMigrationStartTime: map[types.NegSyncerKey]time.Time{
				syncerKey: time.Now(),
			},
			dualStackMigrationEndTime: map[types.NegSyncerKey]time.Time{},
			wantChangedStartTime:      false,
		},
		{
			name:           "end time should get unset because migration is still in progress",
			migrationCount: 5,
			dualStackMigrationStartTime: map[types.NegSyncerKey]time.Time{
				syncerKey: time.Now(),
			},
			dualStackMigrationEndTime: map[types.NegSyncerKey]time.Time{
				syncerKey: time.Now(),
			},
			wantChangedEndTime: true,
		},
		{
			name:           "end time should get set because it's not currently set",
			migrationCount: 0,
			dualStackMigrationStartTime: map[types.NegSyncerKey]time.Time{
				syncerKey: time.Now(),
			},
			dualStackMigrationEndTime: map[types.NegSyncerKey]time.Time{},
			wantChangedEndTime:        true,
		},
		{
			name:           "should not set new end time if end time already present",
			migrationCount: 0,
			dualStackMigrationStartTime: map[types.NegSyncerKey]time.Time{
				syncerKey: time.Now(),
			},
			dualStackMigrationEndTime: map[types.NegSyncerKey]time.Time{
				syncerKey: time.Now(),
			},
		},
		{
			name:                        "migration was not in progress so end time should not get set",
			migrationCount:              0,
			dualStackMigrationStartTime: map[types.NegSyncerKey]time.Time{},
			dualStackMigrationEndTime:   map[types.NegSyncerKey]time.Time{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sm := &SyncerMetrics{
				dualStackMigrationStartTime: tc.dualStackMigrationStartTime,
				dualStackMigrationEndTime:   tc.dualStackMigrationEndTime,
				clock:                       clock.RealClock{},
			}

			clonedStartTime := tc.dualStackMigrationStartTime[syncerKey]
			clonedEndTime := tc.dualStackMigrationEndTime[syncerKey]

			sm.updateMigrationStartAndEndTime(syncerKey, tc.migrationCount)

			if tc.wantChangedStartTime == (sm.dualStackMigrationStartTime[syncerKey] == clonedStartTime) {
				t.Errorf("updateMigrationStartAndEndTime(%v, %v): startTimeBefore=%v, startTimeAfter=%v; want change=%v", syncerKey, tc.migrationCount, clonedStartTime.UnixNano(), sm.dualStackMigrationStartTime[syncerKey].UnixNano(), tc.wantChangedStartTime)
			}
			if tc.wantChangedEndTime == (sm.dualStackMigrationEndTime[syncerKey] == clonedEndTime) {
				t.Errorf("updateMigrationStartAndEndTime(%v, %v): endTimeBefore=%v, endTimeAfter=%v; want change=%v", syncerKey, tc.migrationCount, clonedEndTime.UnixNano(), sm.dualStackMigrationStartTime[syncerKey].UnixNano(), tc.wantChangedEndTime)
			}
		})
	}
}

func TestUpdateEndpointsCountPerType(t *testing.T) {
	syncerKey := syncerKey(1)
	inputCommittedEndpoints := map[string]types.NetworkEndpointSet{
		"zone1": types.NewNetworkEndpointSet([]types.NetworkEndpoint{
			{IP: "ipv4only-1"}, {IP: "ipv4only-2"}, {IP: "ipv4only-3"},
			{IPv6: "ipv6only-1"},
			{IP: "dualStack-1a", IPv6: "dualStack-1b"}, {IP: "dualStack-2a", IPv6: "dualStack-2b"},
		}...),
		"zone2": types.NewNetworkEndpointSet([]types.NetworkEndpoint{
			{IP: "ipv4only-4"},
			{IP: "dualStack-3a", IPv6: "dualStack-3b"},
		}...),
	}
	inputMigrationCount := 10
	inputEndpointsCountPerType := map[types.NegSyncerKey]map[string]int{
		syncerKey: {
			ipv4EndpointType:      10,
			ipv6EndpointType:      9,
			dualStackEndpointType: 15,
			migrationEndpointType: 7,
		},
	}
	sm := &SyncerMetrics{
		endpointsCountPerType: inputEndpointsCountPerType,
	}

	wantEndpointsCountPerType := map[types.NegSyncerKey]map[string]int{
		syncerKey: {
			ipv4EndpointType:      4,
			ipv6EndpointType:      1,
			dualStackEndpointType: 3,
			migrationEndpointType: 10,
		},
	}

	sm.updateEndpointsCountPerType(syncerKey, inputCommittedEndpoints, inputMigrationCount)

	if diff := cmp.Diff(wantEndpointsCountPerType, sm.endpointsCountPerType); diff != "" {
		t.Errorf("updateEndpointsCountPerType(...) = Unexpected diff in endpointsCountPerType: (-want +got):\n%s", diff)
	}
}

func TestComputeDualStackMigrationDurations(t *testing.T) {
	curTime := time.Unix(50, 0)
	inputDualStackMigrationStartTime := map[types.NegSyncerKey]time.Time{
		syncerKey(1): time.Unix(5, 0),
		syncerKey(2): time.Unix(7, 0),
		syncerKey(3): time.Unix(6, 0),
		syncerKey(4): time.Unix(10, 0),
	}
	inputDualStackMigrationEndTime := map[types.NegSyncerKey]time.Time{
		syncerKey(1): time.Unix(10, 0),
		syncerKey(3): time.Unix(9, 0),
	}
	sm := &SyncerMetrics{
		dualStackMigrationStartTime: inputDualStackMigrationStartTime,
		dualStackMigrationEndTime:   inputDualStackMigrationEndTime,
		clock:                       &fakeClock{curTime: curTime},
	}

	wantFinishedDurations := []int{
		int(time.Unix(10, 0).Sub(time.Unix(5, 0)).Seconds()), // 10 - 5
		int(time.Unix(9, 0).Sub(time.Unix(6, 0)).Seconds()),  // 9 - 6
	}
	wantLongestUnfinishedDuration := int(curTime.Sub(time.Unix(7, 0)).Seconds()) // 50 - 7

	gotFinishedDurations, gotLongestUnfinishedDuration := sm.computeDualStackMigrationDurations()

	sortSlices := cmpopts.SortSlices(func(a, b int) bool { return a < b })
	if diff := cmp.Diff(wantFinishedDurations, gotFinishedDurations, sortSlices); diff != "" {
		t.Errorf("computeDualStackMigrationDurations() = Unexpected diff in finishedDurations: (-want +got):\n%s", diff)
	}
	if gotLongestUnfinishedDuration != wantLongestUnfinishedDuration {
		t.Errorf("computeDualStackMigrationDurations() returned longestUnfinishedDuration=%v; want=%v", gotLongestUnfinishedDuration, wantLongestUnfinishedDuration)
	}

	// Ensure that finished durations are not returned more than once but
	// longestUnfinishedDuration duration is returned until it completes.
	gotFinishedDurations, gotLongestUnfinishedDuration = sm.computeDualStackMigrationDurations()
	if len(gotFinishedDurations) != 0 {
		t.Errorf("computeDualStackMigrationDurations() returned non-empty finishedDurations; want finishedDurations to be empty if computeDualStackMigrationDurations is invoked more than once.")
	}
	if gotLongestUnfinishedDuration != wantLongestUnfinishedDuration {
		t.Errorf("computeDualStackMigrationDurations() returned longestUnfinishedDuration=%v; want=%v", gotLongestUnfinishedDuration, wantLongestUnfinishedDuration)
	}
}

func TestComputeDualStackMigrationCounts(t *testing.T) {
	inputEndpointsCountPerType := map[types.NegSyncerKey]map[string]int{
		syncerKeyWithPort(1, 8080): {ipv4EndpointType: 1, ipv6EndpointType: 5, dualStackEndpointType: 9, migrationEndpointType: 13},
		syncerKeyWithPort(2, 8080): {ipv6EndpointType: 6, migrationEndpointType: 14},
		syncerKeyWithPort(1, 8443): {ipv4EndpointType: 3, ipv6EndpointType: 7, dualStackEndpointType: 11},
		syncerKeyWithPort(2, 8443): {ipv6EndpointType: 8, dualStackEndpointType: 12, migrationEndpointType: 16},
		syncerKeyWithPort(3, 80):   {ipv4EndpointType: 10, ipv6EndpointType: 20, dualStackEndpointType: 30},
	}
	sm := &SyncerMetrics{
		endpointsCountPerType: inputEndpointsCountPerType,
	}

	wantSyncerCountByEndpointType := map[string]int{
		ipv4EndpointType:      3,
		ipv6EndpointType:      5,
		dualStackEndpointType: 4,
		migrationEndpointType: 3,
	}
	wantMigrationEndpointCount := 13 + 14 + 16
	wantMigrationServicesCount := 2

	gotSyncerCountByEndpointType, gotMigrationEndpointCount, gotMigrationServicesCount := sm.computeDualStackMigrationCounts()

	if diff := cmp.Diff(wantSyncerCountByEndpointType, gotSyncerCountByEndpointType); diff != "" {
		t.Errorf("computeDualStackMigrationCounts() = Unexpected diff in negCountByEndpointType: (-want +got):\n%s", diff)
	}
	if gotMigrationEndpointCount != wantMigrationEndpointCount {
		t.Errorf("computeDualStackMigrationCounts() returned migrationEndpointCount=%v: want=%v", gotMigrationEndpointCount, wantMigrationEndpointCount)
	}
	if gotMigrationServicesCount != wantMigrationServicesCount {
		t.Errorf("computeDualStackMigrationCounts() returned migrationServicesCount=%v: want=%v", gotMigrationServicesCount, wantMigrationServicesCount)
	}
}

type fakeClock struct {
	clock.Clock
	curTime time.Time
}

func (f *fakeClock) Since(t time.Time) time.Duration {
	return f.curTime.Sub(t)
}

func syncerKey(i int32) types.NegSyncerKey {
	return syncerKeyWithPort(i, i)
}

func syncerKeyWithPort(i int32, port int32) types.NegSyncerKey {
	return types.NegSyncerKey{
		Namespace: fmt.Sprintf("ns-%v", i),
		Name:      fmt.Sprintf("name-%v", i),
		NegName:   fmt.Sprintf("neg-name-%v", i),
		PortTuple: types.SvcPortTuple{Port: port},
	}
}
