package payload

import "mime/multipart"

type UpdateUserPayload struct {
	DisplayName *string        `json:"displayName"`
	Profile     multipart.File `file:"profile"`
}

type UpdateUserPasswordPayload struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}
