package auth

// AuthRequest is contract for vk or gmail auth request
type AuthRequest struct {
	AuthToken string `validate:"required"`
}

// AuthRequest is response for vk or gmail auth
type AuthResponse struct {
	AuthToken string
}
