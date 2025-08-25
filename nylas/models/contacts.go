package models

// ---------- Enums ----------

type SourceType string

const (
	SourceTypeAddressBook SourceType = "address_book"
	SourceTypeInbox       SourceType = "inbox"
	SourceTypeDomain      SourceType = "domain"
)

// ---------- Atomic value objects (read models) ----------

type PhoneNumber struct {
	Number *string `json:"number,omitempty"`
	Type   *string `json:"type,omitempty"`
}

type PhysicalAddress struct {
	Format        *string `json:"format,omitempty"`
	StreetAddress *string `json:"street_address,omitempty"`
	City          *string `json:"city,omitempty"`
	PostalCode    *string `json:"postal_code,omitempty"`
	State         *string `json:"state,omitempty"`
	Country       *string `json:"country,omitempty"`
	Type          *string `json:"type,omitempty"`
}

type WebPage struct {
	URL  *string `json:"url,omitempty"`
	Type *string `json:"type,omitempty"`
}

type ContactEmail struct {
	Email *string `json:"email,omitempty"`
	Type  *string `json:"type,omitempty"`
}

type ContactGroupID struct {
	ID string `json:"id"`
}

type InstantMessagingAddress struct {
	IMAddress *string `json:"im_address,omitempty"`
	Type      *string `json:"type,omitempty"`
}

// ---------- Contact (read model) ----------

type Contact struct {
	ID                string                    `json:"id"`
	GrantID           string                    `json:"grant_id"`
	Object            string                    `json:"object,omitempty"` // "contact"
	Birthday          *string                   `json:"birthday,omitempty"`
	CompanyName       *string                   `json:"company_name,omitempty"`
	DisplayName       *string                   `json:"display_name,omitempty"`
	Emails            []ContactEmail            `json:"emails,omitempty"`
	IMAddresses       []InstantMessagingAddress `json:"im_addresses,omitempty"`
	GivenName         *string                   `json:"given_name,omitempty"`
	JobTitle          *string                   `json:"job_title,omitempty"`
	ManagerName       *string                   `json:"manager_name,omitempty"`
	MiddleName        *string                   `json:"middle_name,omitempty"`
	Nickname          *string                   `json:"nickname,omitempty"`
	Notes             *string                   `json:"notes,omitempty"`
	OfficeLocation    *string                   `json:"office_location,omitempty"`
	PictureURL        *string                   `json:"picture_url,omitempty"`
	Picture           *string                   `json:"picture,omitempty"` // base64 when requested
	Suffix            *string                   `json:"suffix,omitempty"`
	Surname           *string                   `json:"surname,omitempty"`
	Source            *SourceType               `json:"source,omitempty"`
	PhoneNumbers      []PhoneNumber             `json:"phone_numbers,omitempty"`
	PhysicalAddresses []PhysicalAddress         `json:"physical_addresses,omitempty"`
	WebPages          []WebPage                 `json:"web_pages,omitempty"`
	Groups            []ContactGroupID          `json:"groups,omitempty"`
}

// ---------- Find query params ----------

type FindContactQueryParams struct {
	ProfilePicture *bool   `json:"profile_picture,omitempty" url:"profile_picture,omitempty"`
	Select         *string `json:"select,omitempty" url:"select,omitempty"`
}

// ---------- Write models (TypedDicts in Python) ----------

type WriteablePhoneNumber struct {
	Number *string `json:"number,omitempty"`
	Type   *string `json:"type,omitempty"`
}

type WriteablePhysicalAddress struct {
	Format        *string `json:"format,omitempty"`
	StreetAddress *string `json:"street_address,omitempty"`
	City          *string `json:"city,omitempty"`
	PostalCode    *string `json:"postal_code,omitempty"`
	State         *string `json:"state,omitempty"`
	Country       *string `json:"country,omitempty"`
	Type          *string `json:"type,omitempty"`
}

type WriteableWebPage struct {
	URL  *string `json:"url,omitempty"`
	Type *string `json:"type,omitempty"`
}

type WriteableContactEmail struct {
	Email *string `json:"email,omitempty"`
	Type  *string `json:"type,omitempty"`
}

type WriteableContactGroupID struct {
	ID string `json:"id"`
}

type WriteableInstantMessagingAddress struct {
	IMAddress *string `json:"im_address,omitempty"`
	Type      *string `json:"type,omitempty"`
}

// CreateContactRequest mirrors Python; all fields optional (create is sparse-friendly)
type CreateContactRequest struct {
	Birthday          *string                            `json:"birthday,omitempty"`
	CompanyName       *string                            `json:"company_name,omitempty"`
	DisplayName       *string                            `json:"display_name,omitempty"`
	Emails            []WriteableContactEmail            `json:"emails,omitempty"`
	IMAddresses       []WriteableInstantMessagingAddress `json:"im_addresses,omitempty"`
	GivenName         *string                            `json:"given_name,omitempty"`
	JobTitle          *string                            `json:"job_title,omitempty"`
	ManagerName       *string                            `json:"manager_name,omitempty"`
	MiddleName        *string                            `json:"middle_name,omitempty"`
	Nickname          *string                            `json:"nickname,omitempty"`
	Notes             *string                            `json:"notes,omitempty"`
	OfficeLocation    *string                            `json:"office_location,omitempty"`
	PictureURL        *string                            `json:"picture_url,omitempty"`
	Picture           *string                            `json:"picture,omitempty"`
	Suffix            *string                            `json:"suffix,omitempty"`
	Surname           *string                            `json:"surname,omitempty"`
	Source            *SourceType                        `json:"source,omitempty"`
	PhoneNumbers      []WriteablePhoneNumber             `json:"phone_numbers,omitempty"`
	PhysicalAddresses []WriteablePhysicalAddress         `json:"physical_addresses,omitempty"`
	WebPages          []WriteableWebPage                 `json:"web_pages,omitempty"`
	Groups            []WriteableContactGroupID          `json:"groups,omitempty"`
}

// UpdateContactRequest == CreateContactRequest in Python (PUT semantics)
type UpdateContactRequest = CreateContactRequest

// ---------- List query params ----------

type ListContactsQueryParams struct {
	Email       *string     `json:"email,omitempty" url:"email,omitempty"`
	PhoneNumber *string     `json:"phone_number,omitempty" url:"phone_number,omitempty"`
	Source      *SourceType `json:"source,omitempty" url:"source,omitempty"`
	Group       *string     `json:"group,omitempty" url:"group,omitempty"`
	Recurse     *bool       `json:"recurse,omitempty" url:"recurse,omitempty"`

	Select    *string `json:"select,omitempty" url:"select,omitempty"`
	Limit     *int    `json:"limit,omitempty" url:"limit,omitempty"`
	PageToken *string `json:"page_token,omitempty" url:"page_token,omitempty"`
}

// ---------- Contact groups ----------

type GroupType string

const (
	GroupTypeUser   GroupType = "user"
	GroupTypeSystem GroupType = "system"
	GroupTypeOther  GroupType = "other"
)

type ContactGroup struct {
	ID        string     `json:"id"`
	GrantID   string     `json:"grant_id"`
	Object    string     `json:"object,omitempty"` // "contact_group"
	GroupType *GroupType `json:"group_type,omitempty"`
	Name      *string    `json:"name,omitempty"`
	Path      *string    `json:"path,omitempty"`
}

// Simple alias to list-style params if needed elsewhere
type ListContactGroupsQueryParams struct {
	Limit     *int    `json:"limit,omitempty" url:"limit,omitempty"`
	PageToken *string `json:"page_token,omitempty" url:"page_token,omitempty"`
	Select    *string `json:"select,omitempty" url:"select,omitempty"`
}
