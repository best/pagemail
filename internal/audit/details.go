package audit

import "encoding/json"

// DetailsType indicates the structure type of audit log details.
type DetailsType string

const (
	DetailsTypeLogin    DetailsType = "login"
	DetailsTypeResource DetailsType = "resource"
	DetailsTypeChange   DetailsType = "change"
	DetailsTypeRaw      DetailsType = "raw"
)

// LoginDetails for user login events.
type LoginDetails struct {
	UserAgent string `json:"user_agent,omitempty"`
}

// ResourceDetails for CRUD operations on resources.
type ResourceDetails struct {
	Name     string   `json:"name,omitempty"`
	Email    string   `json:"email,omitempty"`
	URL      string   `json:"url,omitempty"`
	Role     string   `json:"role,omitempty"`
	Host     string   `json:"host,omitempty"`
	Port     int      `json:"port,omitempty"`
	Formats  []string `json:"formats,omitempty"`
	IsActive *bool    `json:"is_active,omitempty"`
}

// ChangeDetails for a single field change.
type ChangeDetails struct {
	Field    string `json:"field"`
	OldValue string `json:"old_value,omitempty"`
	NewValue string `json:"new_value,omitempty"`
}

// ChangeSetDetails for multiple field changes.
type ChangeSetDetails struct {
	Changes []ChangeDetails `json:"changes"`
}

// DetailsTypeForAction maps action to its expected details type.
func DetailsTypeForAction(action string) DetailsType {
	switch action {
	case ActionUserLogin:
		return DetailsTypeLogin
	case ActionUserUpdate, ActionSettingsUpdate:
		return DetailsTypeChange
	case ActionUserCreate, ActionUserDelete,
		ActionSMTPCreate, ActionSMTPUpdate, ActionSMTPDelete,
		ActionWebhookCreate, ActionWebhookUpdate, ActionWebhookDelete,
		ActionCaptureCreate, ActionCaptureDelete,
		ActionDeliveryCreate:
		return DetailsTypeResource
	default:
		return DetailsTypeRaw
	}
}

// NormalizeDetails parses stored JSON and returns typed details.
func NormalizeDetails(action, rawJSON, userAgent string) (detailsType DetailsType, details interface{}) {
	if rawJSON != "" {
		var parsed interface{}
		if err := json.Unmarshal([]byte(rawJSON), &parsed); err == nil {
			return DetailsTypeForAction(action), parsed
		}
		return DetailsTypeRaw, rawJSON
	}
	if action == ActionUserLogin && userAgent != "" {
		return DetailsTypeLogin, LoginDetails{UserAgent: userAgent}
	}
	return DetailsTypeForAction(action), nil
}

// BoolPtr returns a pointer to a bool value.
func BoolPtr(v bool) *bool {
	return &v
}
