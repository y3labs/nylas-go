package models

// ---- Status ----

type WebhookStatus string

const (
	WebhookStatusActive  WebhookStatus = "active"
	WebhookStatusFailing WebhookStatus = "failing"
	WebhookStatusFailed  WebhookStatus = "failed"
	WebhookStatusPause   WebhookStatus = "pause"
)

// ---- Triggers ----

type WebhookTrigger string

const (
	TriggerBookingCreated              WebhookTrigger = "booking.created"
	TriggerBookingPending              WebhookTrigger = "booking.pending"
	TriggerBookingRescheduled          WebhookTrigger = "booking.rescheduled"
	TriggerBookingCancelled            WebhookTrigger = "booking.cancelled"
	TriggerBookingReminder             WebhookTrigger = "booking.reminder"
	TriggerCalendarCreated             WebhookTrigger = "calendar.created"
	TriggerCalendarUpdated             WebhookTrigger = "calendar.updated"
	TriggerCalendarDeleted             WebhookTrigger = "calendar.deleted"
	TriggerContactUpdated              WebhookTrigger = "contact.updated"
	TriggerContactDeleted              WebhookTrigger = "contact.deleted"
	TriggerEventCreated                WebhookTrigger = "event.created"
	TriggerEventUpdated                WebhookTrigger = "event.updated"
	TriggerEventDeleted                WebhookTrigger = "event.deleted"
	TriggerGrantCreated                WebhookTrigger = "grant.created"
	TriggerGrantUpdated                WebhookTrigger = "grant.updated"
	TriggerGrantDeleted                WebhookTrigger = "grant.deleted"
	TriggerGrantExpired                WebhookTrigger = "grant.expired"
	TriggerMessageSendSuccess          WebhookTrigger = "message.send_success"
	TriggerMessageSendFailed           WebhookTrigger = "message.send_failed"
	TriggerMessageBounceDetected       WebhookTrigger = "message.bounce_detected"
	TriggerMessageCreated              WebhookTrigger = "message.created"
	TriggerMessageUpdated              WebhookTrigger = "message.updated"
	TriggerMessageOpened               WebhookTrigger = "message.opened"
	TriggerMessageLinkClicked          WebhookTrigger = "message.link_clicked"
	TriggerMessageOpenedLegacy         WebhookTrigger = "message.opened.legacy"
	TriggerMessageLinkClickedLegacy    WebhookTrigger = "message_link_clicked.legacy"
	TriggerMessageIntelligenceOrder    WebhookTrigger = "message.intelligence.order"
	TriggerMessageIntelligenceTracking WebhookTrigger = "message.intelligence.tracking"
	TriggerMessageIntelligenceReturn   WebhookTrigger = "message.intelligence.return"
	TriggerThreadReplied               WebhookTrigger = "thread.replied"
	TriggerThreadRepliedLegacy         WebhookTrigger = "thread.replied.legacy"
	TriggerFolderCreated               WebhookTrigger = "folder.created"
	TriggerFolderUpdated               WebhookTrigger = "folder.updated"
	TriggerFolderDeleted               WebhookTrigger = "folder.deleted"
)

// ---- Models ----

type Webhook struct {
	ID                         string           `json:"id"`
	TriggerTypes               []WebhookTrigger `json:"trigger_types"`
	WebhookURL                 string           `json:"webhook_url"`
	Status                     WebhookStatus    `json:"status"`
	NotificationEmailAddresses []string         `json:"notification_email_addresses"`
	StatusUpdatedAt            int64            `json:"status_updated_at"`
	CreatedAt                  int64            `json:"created_at"`
	UpdatedAt                  int64            `json:"updated_at"`
	Description                *string          `json:"description,omitempty"`
}

type WebhookWithSecret struct {
	Webhook
	WebhookSecret string `json:"webhook_secret"`
}

type WebhookDeleteData struct {
	Status string `json:"status"`
}

type WebhookDeleteResponse struct {
	RequestID string             `json:"request_id"`
	Data      *WebhookDeleteData `json:"data,omitempty"`
}

type WebhookIpAddressesResponse struct {
	IPAddresses []string `json:"ip_addresses"`
	UpdatedAt   int64    `json:"updated_at"`
}

// ---- Requests ----

type CreateWebhookRequest struct {
	TriggerTypes               []WebhookTrigger `json:"trigger_types"`
	WebhookURL                 string           `json:"webhook_url"`
	Description                *string          `json:"description,omitempty"`
	NotificationEmailAddresses []string         `json:"notification_email_addresses,omitempty"`
}

type UpdateWebhookRequest struct {
	TriggerTypes               []WebhookTrigger `json:"trigger_types,omitempty"`
	WebhookURL                 *string          `json:"webhook_url,omitempty"`
	Description                *string          `json:"description,omitempty"`
	NotificationEmailAddresses []string         `json:"notification_email_addresses,omitempty"`
}
