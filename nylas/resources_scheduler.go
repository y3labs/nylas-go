// nylas-go/nylas/resources_scheduler.go
package nylas

// Scheduler mirrors the Python facade:
//
//	client.Scheduler.Configurations.List(...)
//	client.Scheduler.Bookings.Create(...)
//	client.Scheduler.Sessions.Create(...)
//
// It does not make HTTP calls itself—sub-resources do.

type Scheduler struct{ c *Client }

func (s *Scheduler) Configurations() *ConfigurationsResource { return &ConfigurationsResource{c: s.c} }
func (s *Scheduler) Bookings() *BookingsResource             { return &BookingsResource{c: s.c} }
func (s *Scheduler) Sessions() *SessionsResource             { return &SessionsResource{c: s.c} }
