package models

// ---------- Enums ----------

type NotetakerState string

const (
	NotetakerStateScheduled       NotetakerState = "scheduled"
	NotetakerStateConnecting      NotetakerState = "connecting"
	NotetakerStateWaitingForEntry NotetakerState = "waiting_for_entry"
	NotetakerStateFailedEntry     NotetakerState = "failed_entry"
	NotetakerStateAttending       NotetakerState = "attending"
	NotetakerStateMediaProcessing NotetakerState = "media_processing"
	NotetakerStateMediaAvailable  NotetakerState = "media_available"
	NotetakerStateMediaError      NotetakerState = "media_error"
	NotetakerStateMediaDeleted    NotetakerState = "media_deleted"
)

type NotetakerOrderBy string

const (
	NotetakerOrderByName      NotetakerOrderBy = "name"
	NotetakerOrderByJoinTime  NotetakerOrderBy = "join_time"
	NotetakerOrderByCreatedAt NotetakerOrderBy = "created_at"
)

type NotetakerOrderDirection string

const (
	NotetakerOrderDirectionASC  NotetakerOrderDirection = "asc"
	NotetakerOrderDirectionDESC NotetakerOrderDirection = "desc"
)

type MeetingProvider string

const (
	MeetingProviderGoogleMeet     MeetingProvider = "Google Meet"
	MeetingProviderZoom           MeetingProvider = "Zoom Meeting"
	MeetingProviderMicrosoftTeams MeetingProvider = "Microsoft Teams"
)

// ---------- Requests (TypedDict equivalents) ----------

// NotetakerMeetingSettingsRequest: optional fields for create/update requests
type NotetakerMeetingSettingsRequest struct {
	VideoRecording *bool `json:"video_recording,omitempty"`
	AudioRecording *bool `json:"audio_recording,omitempty"`
	Transcription  *bool `json:"transcription,omitempty"`
}

// InviteNotetakerRequest mirrors the Python creation payload
type InviteNotetakerRequest struct {
	MeetingLink     string                           `json:"meeting_link"`
	JoinTime        *int64                           `json:"join_time,omitempty"` // Unix timestamp (seconds)
	Name            *string                          `json:"name,omitempty"`
	MeetingSettings *NotetakerMeetingSettingsRequest `json:"meeting_settings,omitempty"`
}

// UpdateNotetakerRequest mirrors the Python update payload
type UpdateNotetakerRequest struct {
	JoinTime        *int64                           `json:"join_time,omitempty"` // Unix timestamp (seconds)
	Name            *string                          `json:"name,omitempty"`
	MeetingSettings *NotetakerMeetingSettingsRequest `json:"meeting_settings,omitempty"`
}

// ListNotetakerQueryParams extends typical list params with filters/sorting.
// Include url tags so your EncodeQuery helper can serialize these properly.
type ListNotetakerQueryParams struct {
	// Generic list params (if you support these across SDK)
	Limit         *int    `json:"limit,omitempty" url:"limit,omitempty"`
	PageToken     *string `json:"page_token,omitempty" url:"page_token,omitempty"`
	PrevPageToken *string `json:"prev_page_token,omitempty" url:"prev_page_token,omitempty"`

	// Notetaker-specific filters
	State          *NotetakerState          `json:"state,omitempty" url:"state,omitempty"`
	JoinTimeStart  *int64                   `json:"join_time_start,omitempty" url:"join_time_start,omitempty"`
	JoinTimeEnd    *int64                   `json:"join_time_end,omitempty" url:"join_time_end,omitempty"`
	OrderBy        *NotetakerOrderBy        `json:"order_by,omitempty" url:"order_by,omitempty"`
	OrderDirection *NotetakerOrderDirection `json:"order_direction,omitempty" url:"order_direction,omitempty"`
}

// FindNotetakerQueryParams
type FindNotetakerQueryParams struct {
	Select *string `json:"select,omitempty" url:"select,omitempty"`
}

// ---------- Responses (dataclasses) ----------

type NotetakerMeetingSettings struct {
	// In Python these default to true. The API will return explicit values;
	// for outbound requests, use NotetakerMeetingSettingsRequest.
	VideoRecording bool `json:"video_recording"`
	AudioRecording bool `json:"audio_recording"`
	Transcription  bool `json:"transcription"`
}

type NotetakerMediaRecording struct {
	Size      int64  `json:"size"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	CreatedAt int64  `json:"created_at"` // Unix timestamp (seconds)
	ExpiresAt int64  `json:"expires_at"` // Unix timestamp (seconds)
	URL       string `json:"url"`
	TTL       int64  `json:"ttl"`
}

type NotetakerMedia struct {
	Recording  *NotetakerMediaRecording `json:"recording,omitempty"`
	Transcript *NotetakerMediaRecording `json:"transcript,omitempty"`
}

type Notetaker struct {
	ID              string                   `json:"id"`
	Name            string                   `json:"name"`
	JoinTime        int64                    `json:"join_time"` // Unix timestamp (seconds)
	MeetingLink     string                   `json:"meeting_link"`
	State           NotetakerState           `json:"state"`
	MeetingSettings NotetakerMeetingSettings `json:"meeting_settings"`
	MeetingProvider *MeetingProvider         `json:"meeting_provider,omitempty"`
	Message         *string                  `json:"message,omitempty"`
	Object          string                   `json:"object,omitempty"` // usually "notetaker"
}

// Convenience helpers to mirror Python methods.
func (n Notetaker) IsState(s NotetakerState) bool { return n.State == s }
func (n Notetaker) IsScheduled() bool             { return n.IsState(NotetakerStateScheduled) }
func (n Notetaker) IsAttending() bool             { return n.IsState(NotetakerStateAttending) }
func (n Notetaker) HasMediaAvailable() bool       { return n.IsState(NotetakerStateMediaAvailable) }

type NotetakerLeaveResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Object  string `json:"object,omitempty"` // usually "notetaker_leave_response"
}
