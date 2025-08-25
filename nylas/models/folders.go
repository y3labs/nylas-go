package models

import (
	"bytes"
	"encoding/json"
)

// FolderAttributes is tolerant to API variations:
// it accepts either a JSON string (e.g. "\\Sent") or an array of strings (["\\Sent"]).
type FolderAttributes []string

func (fa *FolderAttributes) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if bytes.Equal(b, []byte("null")) || len(b) == 0 {
		*fa = nil
		return nil
	}
	// Array form
	if len(b) > 0 && b[0] == '[' {
		var vs []string
		if err := json.Unmarshal(b, &vs); err != nil {
			return err
		}
		*fa = FolderAttributes(vs)
		return nil
	}
	// Single string form
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*fa = FolderAttributes([]string{s})
	return nil
}

func (fa FolderAttributes) MarshalJSON() ([]byte, error) {
	if fa == nil {
		return []byte("null"), nil
	}
	return json.Marshal([]string(fa))
}

type Folder struct {
	ID              string           `json:"id"`
	GrantID         string           `json:"grant_id"`
	Name            string           `json:"name"`
	Object          string           `json:"object,omitempty"`           // "folder"
	ParentID        *string          `json:"parent_id,omitempty"`        // MS only
	BackgroundColor *string          `json:"background_color,omitempty"` // Google only
	TextColor       *string          `json:"text_color,omitempty"`       // Google only
	SystemFolder    *bool            `json:"system_folder,omitempty"`    // Google only
	ChildCount      *int             `json:"child_count,omitempty"`      // MS only
	UnreadCount     *int             `json:"unread_count,omitempty"`
	TotalCount      *int             `json:"total_count,omitempty"`
	Attributes      FolderAttributes `json:"attributes,omitempty"` // tolerant string|[]string
}
