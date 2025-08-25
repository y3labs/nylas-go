package models

// FreeBusyError mirrors nylas.models.free_busy.FreeBusyError
type FreeBusyError struct {
	Email string `json:"email"`
	Error string `json:"error"`
}

// FreeBusyTimeSlot mirrors nylas.models.free_busy.TimeSlot (note: different from Availability.TimeSlot)
type FreeBusyTimeSlot struct {
	StartTime int64  `json:"start_time"` // Unix seconds
	EndTime   int64  `json:"end_time"`   // Unix seconds
	Status    string `json:"status"`     // typically "busy"
}

// FreeBusy mirrors nylas.models.free_busy.FreeBusy
type FreeBusy struct {
	Email     string             `json:"email"`
	TimeSlots []FreeBusyTimeSlot `json:"time_slots"`
}

// GetFreeBusyResponseItem is a union-friendly struct that can represent either
// a FreeBusy (with time_slots) OR a FreeBusyError (with error).
// If Error is non-empty, treat this item as an error; otherwise use TimeSlots.
type GetFreeBusyResponseItem struct {
	Email     string             `json:"email"`
	TimeSlots []FreeBusyTimeSlot `json:"time_slots,omitempty"`
	Error     string             `json:"error,omitempty"`
}

// GetFreeBusyResponse mirrors the Python alias: List[Union[FreeBusy, FreeBusyError]]
type GetFreeBusyResponse []GetFreeBusyResponseItem

// GetFreeBusyRequest mirrors nylas.models.free_busy.GetFreeBusyRequest
type GetFreeBusyRequest struct {
	StartTime int64    `json:"start_time"` // Unix seconds
	EndTime   int64    `json:"end_time"`   // Unix seconds
	Emails    []string `json:"emails"`
}
