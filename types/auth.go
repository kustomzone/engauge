package types

// AdminCreds --
type AdminCreds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// NewAdminCreds --
func NewAdminCreds() *AdminCreds {
	return &AdminCreds{}
}
