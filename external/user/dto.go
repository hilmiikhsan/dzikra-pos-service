package user

type GetTokekenValidationResponse struct {
	UserID   string                  `json:"user_id"`
	Email    string                  `json:"email"`
	FullName string                  `json:"full_name"`
	UserRole []ApplicationPermission `json:"user_roles"`
}

type ApplicationPermission struct {
	ApplicationPermission []UserRoleAppPermission `json:"application_permissions"`
	Roles                 string                  `json:"roles"`
}

type UserRoleAppPermission struct {
	ApplicationID string   `json:"application_id"`
	Name          string   `json:"name"`
	Permissions   []string `json:"permissions"`
}
