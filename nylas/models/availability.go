package models

// ----- Enums -----

type AvailabilityMethod string

const (
	AvailabilityMethodMaxFairness     AvailabilityMethod = "max-fairness"
	AvailabilityMethodMaxAvailability AvailabilityMethod = "max-availability"
)

// ----- Core response models -----

type TimeSlot struct {
	// Emails of participants available in this slot.
	Emails    []string `json:"emails"`
	StartTime int64    `json:"start_time"` // Unix seconds
	EndTime   int64    `json:"end_time"`   // Unix seconds
}

type GetAvailabilityResponse struct {
	TimeSlots []TimeSlot `json:"time_slots"`
	// Only populated for round-robin events (order in which accounts are next in line)
	Order []string `json:"order,omitempty"`
}

// ----- Request models -----

type MeetingBuffer struct {
	// Increments of 5 minutes; defaults to 0 when omitted
	Before *int `json:"before,omitempty"`
	After  *int `json:"after,omitempty"`
}

type OpenHours struct {
	// Sunday=0 ... Saturday=6
	Days     []int    `json:"days"`
	Timezone string   `json:"timezone"`          // IANA TZ, e.g. "America/New_York"
	Start    string   `json:"start"`             // "9:00", "13:30" (24h, no leading zeros)
	End      string   `json:"end"`               // "17:00"
	Exdates  []string `json:"exdates,omitempty"` // YYYY-MM-DD
}

type AvailabilityRules struct {
	AvailabilityMethod *AvailabilityMethod `json:"availability_method,omitempty"`
	Buffer             *MeetingBuffer      `json:"buffer,omitempty"`
	DefaultOpenHours   []OpenHours         `json:"default_open_hours,omitempty"`
	RoundRobinGroupID  *string             `json:"round_robin_group_id,omitempty"`
	// Only for Microsoft/EWS; whether tentative events are busy (defaults true)
	TentativeAsBusy *bool `json:"tentative_as_busy,omitempty"`
}

type AvailabilityParticipant struct {
	Email       string      `json:"email"`
	CalendarIDs []string    `json:"calendar_ids,omitempty"`
	OpenHours   []OpenHours `json:"open_hours,omitempty"`
}

type GetAvailabilityRequest struct {
	StartTime         int64                     `json:"start_time"` // Unix seconds
	EndTime           int64                     `json:"end_time"`   // Unix seconds
	Participants      []AvailabilityParticipant `json:"participants"`
	DurationMinutes   int                       `json:"duration_minutes"`
	IntervalMinutes   *int                      `json:"interval_minutes,omitempty"`
	RoundTo30Minutes  *bool                     `json:"round_to_30_minutes,omitempty"` // deprecated; use RoundTo
	AvailabilityRules *AvailabilityRules        `json:"availability_rules,omitempty"`
	// Round each returned slot to a multiple of N minutes (5..60). Overrides RoundTo30Minutes when set.
	RoundTo *int `json:"round_to,omitempty"`
}
