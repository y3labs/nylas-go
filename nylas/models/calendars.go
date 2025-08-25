package models

// ---------- Enums ----------

type EventSelection string

const (
	EventSelectionInternal        EventSelection = "internal"
	EventSelectionExternal        EventSelection = "external"
	EventSelectionOwnEvents       EventSelection = "own_events"
	EventSelectionParticipantOnly EventSelection = "participant_only"
	EventSelectionAll             EventSelection = "all"
)

// ---------- Notetaker config (read models) ----------

type NotetakerParticipantFilter struct {
	ParticipantsGTE *int `json:"participants_gte,omitempty"`
	ParticipantsLTE *int `json:"participants_lte,omitempty"`
}

type NotetakerRules struct {
	EventSelection    []EventSelection            `json:"event_selection,omitempty"`
	ParticipantFilter *NotetakerParticipantFilter `json:"participant_filter,omitempty"`
}

type CalendarNotetaker struct {
	Name            *string                   `json:"name,omitempty"` // default "Nylas Notetaker"
	MeetingSettings *NotetakerMeetingSettings `json:"meeting_settings,omitempty"`
	Rules           *NotetakerRules           `json:"rules,omitempty"`
}

// ---------- Calendar (read model) ----------

type Calendar struct {
	ID                 string             `json:"id"`
	GrantID            string             `json:"grant_id"`
	Name               string             `json:"name"`
	ReadOnly           bool               `json:"read_only"`
	IsOwnedByUser      bool               `json:"is_owned_by_user"`
	Object             string             `json:"object,omitempty"` // "calendar"
	Timezone           *string            `json:"timezone,omitempty"`
	Description        *string            `json:"description,omitempty"`
	Location           *string            `json:"location,omitempty"`
	HexColor           *string            `json:"hex_color,omitempty"`
	HexForegroundColor *string            `json:"hex_foreground_color,omitempty"`
	IsPrimary          *bool              `json:"is_primary,omitempty"`
	Metadata           map[string]any     `json:"metadata,omitempty"`
	Notetaker          *CalendarNotetaker `json:"notetaker,omitempty"`
}

// ---------- Query params ----------

type ListCalendarsQueryParams struct {
	Limit     *int    `json:"limit,omitempty" url:"limit,omitempty"`
	PageToken *string `json:"page_token,omitempty" url:"page_token,omitempty"`
	Select    *string `json:"select,omitempty" url:"select,omitempty"`
	// NOTE: Python uses a dict here. If your EncodeQuery doesn't yet support map→"key:value"
	// pairs for this field, add that behavior (see earlier suggestion) or pre-build a string.
	MetadataPair map[string]string `json:"metadata_pair,omitempty" url:"metadata_pair,omitempty"`
}

type FindCalendarQueryParams struct {
	Select *string `json:"select,omitempty" url:"select,omitempty"`
}

// ---------- Notetaker (write models for requests) ----------

type NotetakerCalendarSettings struct {
	VideoRecording *bool `json:"video_recording,omitempty"`
	AudioRecording *bool `json:"audio_recording,omitempty"`
	Transcription  *bool `json:"transcription,omitempty"`
}

type NotetakerCalendarParticipantFilter struct {
	ParticipantsGTE *int `json:"participants_gte,omitempty"`
	ParticipantsLTE *int `json:"participants_lte,omitempty"`
}

type NotetakerCalendarRules struct {
	EventSelection    []EventSelection                    `json:"event_selection,omitempty"`
	ParticipantFilter *NotetakerCalendarParticipantFilter `json:"participant_filter,omitempty"`
}

type NotetakerCalendarRequest struct {
	Name            *string                    `json:"name,omitempty"`
	MeetingSettings *NotetakerCalendarSettings `json:"meeting_settings,omitempty"`
	Rules           *NotetakerCalendarRules    `json:"rules,omitempty"`
}

// ---------- Create/Update payloads ----------

type CreateCalendarRequest struct {
	Name        string                    `json:"name"`
	Description *string                   `json:"description,omitempty"`
	Location    *string                   `json:"location,omitempty"`
	Timezone    *string                   `json:"timezone,omitempty"`
	Metadata    map[string]string         `json:"metadata,omitempty"`
	Notetaker   *NotetakerCalendarRequest `json:"notetaker,omitempty"`
}

type UpdateCalendarRequest struct {
	// Reuse Create fields as optional for PATCH/PUT semantics
	Name        *string                   `json:"name,omitempty"`
	Description *string                   `json:"description,omitempty"`
	Location    *string                   `json:"location,omitempty"`
	Timezone    *string                   `json:"timezone,omitempty"`
	Metadata    map[string]string         `json:"metadata,omitempty"`
	Notetaker   *NotetakerCalendarRequest `json:"notetaker,omitempty"`

	// Color options (Google-only for foreground)
	HexColor           *string `json:"hex_color,omitempty"`
	HexForegroundColor *string `json:"hex_foreground_color,omitempty"`
}
