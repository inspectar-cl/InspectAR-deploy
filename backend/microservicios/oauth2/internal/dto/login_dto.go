package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	DeviceID string `json:"device_id" binding:"required"`
}
