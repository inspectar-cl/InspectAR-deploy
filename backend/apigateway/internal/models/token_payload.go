package models

type TokenPayload struct {
	Username   string      `json:"username"`
	Email      string      `json:"email"`
	Device     string      `json:"device"`
	Scope      string      `json:"scope"`
	Expiration interface{} `json:"exp"`
}
