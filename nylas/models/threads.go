package models

import (
	"encoding/json"
)

// DraftOrMessage is a tagged-union wrapper for a Message or Draft.
// It inspects the "object" field in the JSON ("message" | "draft").
type DraftOrMessage struct {
	Message *Message
	Draft   *Draft
}

func (u *DraftOrMessage) UnmarshalJSON(b []byte) error {
	// Peek at "object"
	var probe struct {
		Object string `json:"object"`
	}
	if err := json.Unmarshal(b, &probe); err != nil {
		return err
	}
	switch probe.Object {
	case "message":
		var m Message
		if err := json.Unmarshal(b, &m); err != nil {
			return err
		}
		u.Message = &m
		u.Draft = nil
		return nil
	case "draft":
		var d Draft
		if err := json.Unmarshal(b, &d); err != nil {
			return err
		}
		u.Draft = &d
		u.Message = nil
		return nil
	default:
		// Unknown; ignore gracefully
		u.Message, u.Draft = nil, nil
		return nil
	}
}

func (u DraftOrMessage) MarshalJSON() ([]byte, error) {
	switch {
	case u.Message != nil:
		return json.Marshal(u.Message)
	case u.Draft != nil:
		return json.Marshal(u.Draft)
	default:
		return []byte("null"), nil
	}
}

type Thread struct {
	ID                        string         `json:"id"`
	GrantID                   string         `json:"grant_id"`
	HasDrafts                 bool           `json:"has_drafts"`
	Starred                   bool           `json:"starred"`
	Unread                    bool           `json:"unread"`
	MessageIDs                []string       `json:"message_ids"`
	Folders                   []string       `json:"folders"`
	LatestDraftOrMessage      DraftOrMessage `json:"latest_draft_or_message"`
	Object                    string         `json:"object,omitempty"` // "thread"
	EarliestMessageDate       *int64         `json:"earliest_message_date,omitempty"`
	LatestMessageReceivedDate *int64         `json:"latest_message_received_date,omitempty"`
	DraftIDs                  []string       `json:"draft_ids,omitempty"`
	Snippet                   *string        `json:"snippet,omitempty"`
	Subject                   *string        `json:"subject,omitempty"`
	Participants              []EmailName    `json:"participants,omitempty"`
	LatestMessageSentDate     *int64         `json:"latest_message_sent_date,omitempty"`
	HasAttachments            *bool          `json:"has_attachments,omitempty"`
}

// ---- Updates ----

type UpdateThreadRequest struct {
	Starred  *bool          `json:"starred,omitempty"`
	Unread   *bool          `json:"unread,omitempty"`
	Folders  []string       `json:"folders,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// ---- Query params ----
// Mirrors Python ListThreadsQueryParams. Note: many of these are single strings,
// matching the Python API semantics (comma-separated if multiple).

type ListThreadsQueryParams struct {
	// Base list
	Limit     *int    `url:"limit,omitempty"`
	PageToken *string `url:"page_token,omitempty"`
	Select    *string `url:"select,omitempty"`

	Subject             *string `url:"subject,omitempty"`
	AnyEmail            *string `url:"any_email,omitempty"`
	From                *string `url:"from,omitempty"`
	To                  *string `url:"to,omitempty"`
	Cc                  *string `url:"cc,omitempty"`
	Bcc                 *string `url:"bcc,omitempty"`
	In                  *string `url:"in,omitempty"`
	Unread              *bool   `url:"unread,omitempty"`
	Starred             *bool   `url:"starred,omitempty"`
	ThreadID            *string `url:"thread_id,omitempty"`
	EarliestMessageDate *int64  `url:"earliest_message_date,omitempty"`
	LatestMessageBefore *int64  `url:"latest_message_before,omitempty"`
	LatestMessageAfter  *int64  `url:"latest_message_after,omitempty"`
	HasAttachment       *bool   `url:"has_attachment,omitempty"`
	SearchQueryNative   *string `url:"search_query_native,omitempty"`
}
