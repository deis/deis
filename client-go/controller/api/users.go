package api

// User is the definition of the user object.
type User struct {
	ID          int    `json:"id"`
	LastLogin   string `json:"last_login"`
	IsSuperuser bool   `json:"is_superuser"`
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	IsStaff     bool   `json:"is_staff"`
	IsActive    bool   `json:"is_active"`
	DateJoined  string `json:"date_joined"`
}
