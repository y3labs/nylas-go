package models

import (
	"encoding/json"
	"fmt"
)

// --- literals/enums (use string for compile-time flexibility) ---

type Status string               // "confirmed" | "tentative" | "cancelled"
type Visibility string           // "default" | "public" | "private"
type ParticipantStatus string    // "noreply" | "yes" | "no" | "maybe"
type SendRSVPStatus string       // "yes" | "no" | "maybe"
type EventType string            // "default" | "outOfOffice" | "focusTime" | "workingLocation"
type ConferencingProvider string // "Google Meet" | "Zoom Meeting" | "Microsoft Teams" | "GoToMeeting" | "WebEx" | "unknown"

// --- leaf structs ---

type Participant struct {
	Email       string             `json:"email"`
	Status      *ParticipantStatus `json:"status,omitempty"`
	Name        *string            `json:"name,omitempty"`
	Comment     *string            `json:"comment,omitempty"`
	PhoneNumber *string            `json:"phone_number,omitempty"`
}

// --- When (polymorphic: time | timespan | date | datespan) ---

type whenDiscriminator struct {
	Object string `json:"object"`
}

type Time struct {
	Object   string  `json:"object"` // "time"
	Time     int64   `json:"time"`
	Timezone *string `json:"timezone,omitempty"`
}

type Timespan struct {
	Object        string  `json:"object"` // "timespan"
	StartTime     int64   `json:"start_time"`
	EndTime       int64   `json:"end_time"`
	StartTimezone *string `json:"start_timezone,omitempty"`
	EndTimezone   *string `json:"end_timezone,omitempty"`
}

type Date struct {
	Object string `json:"object"` // "date"
	Date   string `json:"date"`
}

type Datespan struct {
	Object    string `json:"object"` // "datespan"
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// When is a tagged union wrapper that marshals/unmarshals to one of the above.
type When struct {
	Time     *Time
	Timespan *Timespan
	Date     *Date
	Datespan *Datespan
}

func (w When) MarshalJSON() ([]byte, error) {
	switch {
	case w.Time != nil:
		return json.Marshal(w.Time)
	case w.Timespan != nil:
		return json.Marshal(w.Timespan)
	case w.Date != nil:
		return json.Marshal(w.Date)
	case w.Datespan != nil:
		return json.Marshal(w.Datespan)
	default:
		return []byte("null"), nil
	}
}

func (w *When) UnmarshalJSON(b []byte) error {
	var d whenDiscriminator
	if err := json.Unmarshal(b, &d); err != nil {
		return err
	}
	switch d.Object {
	case "time":
		var v Time
		if err := json.Unmarshal(b, &v); err != nil {
			return err
		}
		w.Time = &v
	case "timespan":
		var v Timespan
		if err := json.Unmarshal(b, &v); err != nil {
			return err
		}
		w.Timespan = &v
	case "date":
		var v Date
		if err := json.Unmarshal(b, &v); err != nil {
			return err
		}
		w.Date = &v
	case "datespan":
		var v Datespan
		if err := json.Unmarshal(b, &v); err != nil {
			return err
		}
		w.Datespan = &v
	default:
		return fmt.Errorf("unknown when.object: %q", d.Object)
	}
	return nil
}

// --- Conferencing (union: details | autocreate) ---

type DetailsConfig struct {
	MeetingCode *string  `json:"meeting_code,omitempty"`
	Password    *string  `json:"password,omitempty"`
	URL         *string  `json:"url,omitempty"`
	PIN         *string  `json:"pin,omitempty"`
	Phone       []string `json:"phone,omitempty"`
}

type ConferencingDetails struct {
	Provider ConferencingProvider `json:"provider"`
	Details  map[string]any       `json:"details"`
}

type ConferencingAutocreate struct {
	Provider   ConferencingProvider `json:"provider"`
	Autocreate map[string]any       `json:"autocreate"`
}

type conferencingDiscriminator struct {
	Provider *ConferencingProvider `json:"provider"`
	Details  *json.RawMessage      `json:"details,omitempty"`
	Auto     *json.RawMessage      `json:"autocreate,omitempty"`
}

type Conferencing struct {
	Details    *ConferencingDetails
	Autocreate *ConferencingAutocreate
}

func (c Conferencing) MarshalJSON() ([]byte, error) {
	switch {
	case c.Details != nil:
		return json.Marshal(c.Details)
	case c.Autocreate != nil:
		return json.Marshal(c.Autocreate)
	default:
		return []byte("null"), nil
	}
}

func (c *Conferencing) UnmarshalJSON(b []byte) error {
	if string(b) == "null" || len(b) == 0 {
		return nil
	}
	var d conferencingDiscriminator
	if err := json.Unmarshal(b, &d); err != nil {
		return err
	}
	// If provider missing, accept as nil (back-compat)
	if d.Provider == nil {
		return nil
	}
	// Prefer details if present
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	if _, ok := m["details"]; ok {
		var v ConferencingDetails
		if err := json.Unmarshal(b, &v); err != nil {
			return err
		}
		c.Details = &v
		return nil
	}
	if _, ok := m["autocreate"]; ok {
		var v ConferencingAutocreate
		if err := json.Unmarshal(b, &v); err != nil {
			return err
		}
		c.Autocreate = &v
		return nil
	}
	// provider only => treat as empty details
	c.Details = &ConferencingDetails{
		Provider: *d.Provider,
		Details:  map[string]any{},
	}
	return nil
}

func (c *Conferencing) isEmpty() bool {
	if c == nil {
		return true
	}
	// empty if neither variant was populated
	if c.Details == nil && c.Autocreate == nil {
		return true
	}
	// defensive: treat provider-less variants as empty
	if c.Details != nil && c.Details.Provider == "" {
		return true
	}
	if c.Autocreate != nil && c.Autocreate.Provider == "" {
		return true
	}
	return false
}

func (e *Event) UnmarshalJSON(b []byte) error {
	type alias Event
	aux := (*alias)(e)
	if err := json.Unmarshal(b, aux); err != nil {
		return err
	}
	// Normalize conferencing: if it decoded to a zero-value shell, drop it.
	if e.Conferencing != nil && e.Conferencing.isEmpty() {
		e.Conferencing = nil
	}
	return nil
}

// --- reminders ---

type ReminderOverride struct {
	ReminderMinutes *int    `json:"reminder_minutes,omitempty"`
	ReminderMethod  *string `json:"reminder_method,omitempty"`
}

type Reminders struct {
	UseDefault bool               `json:"use_default"`
	Overrides  []ReminderOverride `json:"overrides,omitempty"`
}

type EventNotetaker struct {
	ID              *string                   `json:"id,omitempty"`
	Name            *string                   `json:"name,omitempty"` // default "Nylas Notetaker"
	MeetingSettings *NotetakerMeetingSettings `json:"meeting_settings,omitempty"`
}

// --- Event model (aligns with Python SDK) ---

type Event struct {
	ID               string          `json:"id"`
	GrantID          string          `json:"grant_id"`
	CalendarID       string          `json:"calendar_id"`
	Busy             bool            `json:"busy"`
	Participants     []Participant   `json:"participants"`
	When             When            `json:"when"`
	Conferencing     *Conferencing   `json:"conferencing,omitempty"`
	Object           string          `json:"object"` // "event"
	Visibility       *Visibility     `json:"visibility,omitempty"`
	ReadOnly         *bool           `json:"read_only,omitempty"`
	Description      *string         `json:"description,omitempty"`
	Location         *string         `json:"location,omitempty"`
	ICALUID          *string         `json:"ical_uid,omitempty"`
	Title            *string         `json:"title,omitempty"`
	HTMLLink         *string         `json:"html_link,omitempty"`
	HideParticipants *bool           `json:"hide_participants,omitempty"`
	Metadata         map[string]any  `json:"metadata,omitempty"`
	Creator          *EmailName      `json:"creator,omitempty"`
	Organizer        *EmailName      `json:"organizer,omitempty"`
	Recurrence       []string        `json:"recurrence,omitempty"`
	Reminders        *Reminders      `json:"reminders,omitempty"`
	Status           *Status         `json:"status,omitempty"`
	Capacity         *int            `json:"capacity,omitempty"`
	CreatedAt        *int64          `json:"created_at,omitempty"`
	UpdatedAt        *int64          `json:"updated_at,omitempty"`
	MasterEventID    *string         `json:"master_event_id,omitempty"`
	Notetaker        *EventNotetaker `json:"notetaker,omitempty"`
}
