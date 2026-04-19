package modules

type User struct {
	ID       int    `db:"id"       json:"id"`
	Username string `db:"username" json:"username"`
	Email    string `db:"email"    json:"email"`
	Password string `db:"password" json:"-"`
	Role     string `db:"role"     json:"role"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"     default:"user"`
}

type AuthResponse struct {
	SessionID string `json:"session_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}
