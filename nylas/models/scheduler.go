package models

// --- Enums ---

type BookingType string

const (
	BookingTypeBooking               BookingType = "booking"
	BookingTypeOrganizerConfirmation BookingType = "organizer-confirmation"
)

type BookingReminderType string

const (
	BookingReminderEmail   BookingReminderType = "email"
	BookingReminderWebhook BookingReminderType = "webhook"
)

type BookingRecipientType string

const (
	BookingRecipientHost  BookingRecipientType = "host"
	BookingRecipientGuest BookingRecipientType = "guest"
	BookingRecipientAll   BookingRecipientType = "all"
)

type EmailLanguage string

const (
	EmailLangEN EmailLanguage = "en"
	EmailLangES EmailLanguage = "es"
	EmailLangFR EmailLanguage = "fr"
	EmailLangDE EmailLanguage = "de"
	EmailLangNL EmailLanguage = "nl"
	EmailLangSV EmailLanguage = "sv"
	EmailLangJA EmailLanguage = "ja"
	EmailLangZH EmailLanguage = "zh"
)

type AdditionalFieldType string

const (
	AdditionalFieldText        AdditionalFieldType = "text"
	AdditionalFieldMultiLine   AdditionalFieldType = "multi_line_text"
	AdditionalFieldEmail       AdditionalFieldType = "email"
	AdditionalFieldPhoneNumber AdditionalFieldType = "phone_number"
	AdditionalFieldDropdown    AdditionalFieldType = "dropdown"
	AdditionalFieldDate        AdditionalFieldType = "date"
	AdditionalFieldCheckbox    AdditionalFieldType = "checkbox"
	AdditionalFieldRadioButton AdditionalFieldType = "radio_button"
)

// NOTE: Python’s type hints say “list of options”, so model this as []string.
type AdditionalFieldOption = string

// --- Templates & Settings ---

type BookingConfirmedTemplate struct {
	Title *string `json:"title,omitempty"`
	Body  *string `json:"body,omitempty"`
}

type EmailTemplate struct {
	// Logo intentionally omitted in the Python version; keep parity.
	BookingConfirmed *BookingConfirmedTemplate `json:"booking_confirmed,omitempty"`
}

type AdditionalField struct {
	Label    string                  `json:"label"`
	Type     AdditionalFieldType     `json:"type"`
	Required bool                    `json:"required"`
	Pattern  *string                 `json:"pattern,omitempty"`
	Order    *int                    `json:"order,omitempty"`
	Options  []AdditionalFieldOption `json:"options,omitempty"` // used for dropdown/radio_button
}

type SchedulerSettings struct {
	AdditionalFields         map[string]AdditionalField `json:"additional_fields,omitempty"`
	AvailableDaysInFuture    *int                       `json:"available_days_in_future,omitempty"`
	MinBookingNotice         *int                       `json:"min_booking_notice,omitempty"`
	MinCancellationNotice    *int                       `json:"min_cancellation_notice,omitempty"`
	CancellationPolicy       *string                    `json:"cancellation_policy,omitempty"`
	ReschedulingURL          *string                    `json:"rescheduling_url,omitempty"`
	CancellationURL          *string                    `json:"cancellation_url,omitempty"`
	OrganizerConfirmationURL *string                    `json:"organizer_confirmation_url,omitempty"`
	ConfirmationRedirectURL  *string                    `json:"confirmation_redirect_url,omitempty"`
	HideReschedulingOptions  *bool                      `json:"hide_rescheduling_options,omitempty"`
	HideCancellationOptions  *bool                      `json:"hide_cancellation_options,omitempty"`
	HideAdditionalGuests     *bool                      `json:"hide_additional_guests,omitempty"`
	EmailTemplate            *EmailTemplate             `json:"email_template,omitempty"`
}

// --- Booking configuration ---

type BookingReminder struct {
	Type               BookingReminderType   `json:"type"`
	MinutesBeforeEvent int                   `json:"minutes_before_event"`
	Recipient          *BookingRecipientType `json:"recipient,omitempty"`
	EmailSubject       *string               `json:"email_subject,omitempty"`
}

type EventBooking struct {
	Title         string            `json:"title"`
	Description   *string           `json:"description,omitempty"`
	Location      *string           `json:"location,omitempty"`
	Timezone      *string           `json:"timezone,omitempty"`
	BookingType   *BookingType      `json:"booking_type,omitempty"`
	Conferencing  *Conferencing     `json:"conferencing,omitempty"`
	DisableEmails *bool             `json:"disable_emails,omitempty"`
	Reminders     []BookingReminder `json:"reminders,omitempty"`
}

type Availability struct {
	DurationMinutes   int                `json:"duration_minutes"`
	IntervalMinutes   *int               `json:"interval_minutes,omitempty"`
	RoundTo           *int               `json:"round_to,omitempty"`
	AvailabilityRules *AvailabilityRules `json:"availability_rules,omitempty"`
}

type ParticipantBooking struct {
	CalendarID string `json:"calendar_id"`
}

type ParticipantAvailability struct {
	CalendarIDs []string    `json:"calendar_ids"`
	OpenHours   []OpenHours `json:"open_hours,omitempty"`
}

type ConfigParticipant struct {
	Email        string                  `json:"email"`
	Availability ParticipantAvailability `json:"availability"`
	Booking      ParticipantBooking      `json:"booking"`
	Name         *string                 `json:"name,omitempty"`
	IsOrganizer  *bool                   `json:"is_organizer,omitempty"`
	Timezone     *string                 `json:"timezone,omitempty"`
}

type Configuration struct {
	ID                  string              `json:"id"`
	Participants        []ConfigParticipant `json:"participants"`
	Availability        Availability        `json:"availability"`
	EventBooking        EventBooking        `json:"event_booking"`
	Slug                *string             `json:"slug,omitempty"`
	RequiresSessionAuth *bool               `json:"requires_session_auth,omitempty"`
	Scheduler           *SchedulerSettings  `json:"scheduler,omitempty"`
	Appearance          map[string]string   `json:"appearance,omitempty"`
}

// --- Create/Update Configuration payloads ---

type CreateConfigurationRequest struct {
	Participants        []ConfigParticipant `json:"participants"`
	Availability        Availability        `json:"availability"`
	EventBooking        EventBooking        `json:"event_booking"`
	Slug                *string             `json:"slug,omitempty"`
	RequiresSessionAuth *bool               `json:"requires_session_auth,omitempty"`
	Scheduler           *SchedulerSettings  `json:"scheduler,omitempty"`
	Appearance          map[string]string   `json:"appearance,omitempty"`
}

type UpdateConfigurationRequest struct {
	Participants        *[]ConfigParticipant `json:"participants,omitempty"`
	Availability        *Availability        `json:"availability,omitempty"`
	EventBooking        *EventBooking        `json:"event_booking,omitempty"`
	Slug                *string              `json:"slug,omitempty"`
	RequiresSessionAuth *bool                `json:"requires_session_auth,omitempty"`
	Scheduler           *SchedulerSettings   `json:"scheduler,omitempty"`
	Appearance          *map[string]string   `json:"appearance,omitempty"`
}

// --- Session ---

type CreateSessionRequest struct {
	ConfigurationID *string `json:"configuration_id,omitempty"`
	Slug            *string `json:"slug,omitempty"`
	TimeToLive      *int    `json:"time_to_live,omitempty"` // seconds
}

type Session struct {
	SessionID string `json:"session_id"`
}

// --- Booking (create/confirm/cancel/reschedule) ---

type BookingGuest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type BookingParticipant struct {
	Email string `json:"email"`
}

type CreateBookingRequest struct {
	StartTime        int64                `json:"start_time"`
	EndTime          int64                `json:"end_time"`
	Guest            BookingGuest         `json:"guest"`
	Participants     []BookingParticipant `json:"participants,omitempty"`
	Timezone         *string              `json:"timezone,omitempty"`
	EmailLanguage    *EmailLanguage       `json:"email_language,omitempty"`
	AdditionalGuests []BookingGuest       `json:"additional_guests,omitempty"`
	AdditionalFields map[string]string    `json:"additional_fields,omitempty"`
}

type BookingOrganizer struct {
	Email string  `json:"email"`
	Name  *string `json:"name,omitempty"`
}

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
)

type Booking struct {
	BookingID   string           `json:"booking_id"`
	EventID     string           `json:"event_id"`
	Title       string           `json:"title"`
	Organizer   BookingOrganizer `json:"organizer"`
	Status      BookingStatus    `json:"status"`
	Description *string          `json:"description,omitempty"`
}

type ConfirmBookingStatus string

const (
	ConfirmBookingConfirm ConfirmBookingStatus = "confirm"
	ConfirmBookingCancel  ConfirmBookingStatus = "cancel"
)

type ConfirmBookingRequest struct {
	Salt               string               `json:"salt"`
	Status             ConfirmBookingStatus `json:"status"`
	CancellationReason *string              `json:"cancellation_reason,omitempty"`
}

type DeleteBookingRequest struct {
	CancellationReason *string `json:"cancellation_reason,omitempty"`
}

type RescheduleBookingRequest struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

// --- Query params (for future resources) ---

type CreateBookingQueryParams struct {
	ConfigurationID *string `url:"configuration_id,omitempty"`
	Slug            *string `url:"slug,omitempty"`
	Timezone        *string `url:"timezone,omitempty"`
}

type FindBookingQueryParams struct {
	ConfigurationID *string `url:"configuration_id,omitempty"`
	Slug            *string `url:"slug,omitempty"`
	ClientID        *string `url:"client_id,omitempty"`
}

// Aliases
type ConfirmBookingQueryParams = FindBookingQueryParams
type RescheduleBookingQueryParams = FindBookingQueryParams
type DestroyBookingQueryParams = FindBookingQueryParams
