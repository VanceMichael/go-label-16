package recovery

import (
	"encoding/json"
	"testing"
	"time"
)

func TestSnapshotLabelsRemainIsolatedAcrossBuildAndRestore(t *testing.T) {
	now := time.Date(2026, 8, 21, 6, 0, 0, 0, time.UTC)
	callerLabels := map[string]string{" Barn ": " north-7 ", "Source": "mobile"}
	records := []Record{{
		Kind: "feed_plan", ID: "plan-16", TenantID: "farm-east", Version: 2, UpdatedAt: now,
		Payload: json.RawMessage(`{"status":"scheduled"}`), Labels: callerLabels,
	}}
	snapshot, err := Build("farm-east", 3, records, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("build snapshot: %v", err)
	}
	if callerLabels[" Barn "] != " north-7 " || callerLabels["Source"] != "mobile" {
		t.Fatalf("snapshot build changed caller labels: %+v", callerLabels)
	}
	if snapshot.Records[0].Labels["barn"] != "north-7" || snapshot.Records[0].Labels["source"] != "mobile" {
		t.Fatalf("snapshot labels were not normalized: %+v", snapshot.Records[0].Labels)
	}

	plan, err := Plan(snapshot, nil, nil)
	if err != nil {
		t.Fatalf("plan restore: %v", err)
	}
	restored, err := Apply(nil, plan)
	if err != nil {
		t.Fatalf("apply restore: %v", err)
	}
	restored[0].Labels["barn"] = "south-2"
	restored[0].Labels["operator_note"] = "local recovery annotation"
	if snapshot.Records[0].Labels["barn"] != "north-7" {
		t.Fatalf("restored record changed snapshot barn: %+v", snapshot.Records[0].Labels)
	}
	if _, polluted := snapshot.Records[0].Labels["operator_note"]; polluted {
		t.Fatalf("restored record added label to snapshot: %+v", snapshot.Records[0].Labels)
	}
	if err := Verify(snapshot); err != nil {
		t.Fatalf("snapshot verification after restore annotation: %v", err)
	}
}
