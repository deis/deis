package api

// AuthRegisterRequest is the definition of POST /v1/auth/register/.
type AuthRegisterRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

// AuthLoginRequest is the definition of POST /v1/auth/login/.
type AuthLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthLoginResponse is the definition of /v1/auth/login/.
type AuthLoginResponse tokenResponse

// AuthPasswdRequest is the definition of POST /v1/auth/passwd/.
type AuthPasswdRequest struct {
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	NewPassword string `json:"new_password"`
}

// AuthRegenerateRequest is the definition of POST /v1/auth/tokens/.
type AuthRegenerateRequest struct {
	Name string `json:"username,omitempty"`
	All  bool   `json:"all,omitempty"`
}

// AuthCancelRequest is the definition of POST /v1/auth/cancel/.
type AuthCancelRequest struct {
	Username string `json:"username"`
}

// AuthRegenerateResponse is the definition of /v1/auth/tokens/.
type AuthRegenerateResponse tokenResponse

// A generic defenition of a token response.
type tokenResponse struct {
	Token string `json:"token"`
}
