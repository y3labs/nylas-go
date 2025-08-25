package nylas

import (
	"strconv"
	"strings"
	"time"
)

// RRuleOptions provides a typed way to construct an RFC5545 RRULE string.
// Only commonly used parts are included; extend as needed.
type RRuleOptions struct {
	Freq       string // SECONDLY, MINUTELY, HOURLY, DAILY, WEEKLY, MONTHLY, YEARLY
	Interval   *int
	Count      *int
	Until      *time.Time // will be formatted in UTC as YYYYMMDD'T'HHMMSS'Z'
	ByDay      []string   // MO, TU, WE, TH, FR, SA, SU
	ByMonth    []int      // 1..12
	ByMonthDay []int      // 1..31 or negative for from-end
	WeekStart  *string    // e.g., "MO"
}

func BuildRRULE(o RRuleOptions) string {
	parts := []string{}
	if o.Freq != "" {
		parts = append(parts, "FREQ="+strings.ToUpper(o.Freq))
	}
	if o.Interval != nil {
		parts = append(parts, "INTERVAL="+strconv.Itoa(*o.Interval))
	}
	if o.Count != nil {
		parts = append(parts, "COUNT="+strconv.Itoa(*o.Count))
	}
	if o.Until != nil {
		parts = append(parts, "UNTIL="+o.Until.UTC().Format("20060102T150405Z"))
	}
	if len(o.ByDay) > 0 {
		parts = append(parts, "BYDAY="+strings.Join(o.ByDay, ","))
	}
	if len(o.ByMonth) > 0 {
		parts = append(parts, "BYMONTH="+joinInts(o.ByMonth))
	}
	if len(o.ByMonthDay) > 0 {
		parts = append(parts, "BYMONTHDAY="+joinInts(o.ByMonthDay))
	}
	if o.WeekStart != nil {
		parts = append(parts, "WKST="+strings.ToUpper(*o.WeekStart))
	}
	return strings.Join(parts, ";")
}

func joinInts(xs []int) string {
	ss := make([]string, 0, len(xs))
	for _, v := range xs {
		ss = append(ss, strconv.Itoa(v))
	}
	return strings.Join(ss, ",")
}
