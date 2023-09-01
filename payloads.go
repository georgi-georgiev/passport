package passport

type RecoveryEmailPayload struct {
	Email string `json:"email" binding:"required"`
}

type ExchangeCodeRequestPayload struct {
	Email string `json:"email" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type PasswordResetPayload struct {
	Email    string `json:"email" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type InviteUserPayload struct {
	Email string `json:"email"`
}

type CreateUserPayload struct {
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required"`
	Username string   `json:"username" binding:"required"`
	Role     string   `json:"role" binding:"required"`
	Rights   []string `json:"rights"`
}

type UpdateUserPayload struct {
	Email                string `json:"email"`
	Username             string `json:"username"`
	IsActive             bool   `json:"isActive"`
	ShouldChangePassword bool   `json:"shoudlChangePassword"`
}

type CreateRolePayload struct {
	Name string `json:"name" binding:"required"`
}

type UpdateRolePayload struct {
	Name string `json:"name" binding:"required"`
}

type CreateRightPayload struct {
	Name string `json:"name" binding:"required"`
}

type UpdateRightPayload struct {
	Name string `json:"name" binding:"required"`
}
