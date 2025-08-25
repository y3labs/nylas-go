package nylas

import (
	"testing"
	"time"
)

func TestBuildRRULE_AllFields(t *testing.T) {
	interval := 2
	count := 5
	wkst := "mo"
	until := time.Date(2025, 8, 24, 17, 30, 5, 0, time.UTC)

	got := BuildRRULE(RRuleOptions{
		Freq:       "weekly",
		Interval:   &interval,
		Count:      &count,
		Until:      &until,
		ByDay:      []string{"MO", "WE", "FR"},
		ByMonth:    []int{1, 6, 12},
		ByMonthDay: []int{1, -1},
		WeekStart:  &wkst,
	})

	want := "FREQ=WEEKLY;INTERVAL=2;COUNT=5;UNTIL=20250824T173005Z;BYDAY=MO,WE,FR;BYMONTH=1,6,12;BYMONTHDAY=1,-1;WKST=MO"
	if got != want {
		t.Fatalf("rrule got %q, want %q", got, want)
	}
}

func TestBuildRRULE_Minimal(t *testing.T) {
	got := BuildRRULE(RRuleOptions{Freq: "daily"})
	if got != "FREQ=DAILY" {
		t.Fatalf("got %q", got)
	}
}
