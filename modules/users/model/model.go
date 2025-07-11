package service

import "time"

type CreateUserRequest struct {
	ShopID     int32  `json:"shop_id"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Address    string `json:"address"`
	City       string `json:"city"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
	CodeArea   string `json:"code_area"` // Optional, can be empty
}

type UpdateUserRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Address    string `json:"address"`
	City       string `json:"city"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
}

type ListUsersRequest struct {
	Page     int32 `json:"page"`
	PageSize int32 `json:"limit"`
	// ShopID   int32 `json:"shop_id"`
}

type UserResponse struct {
	ID         int32  `json:"id"`
	ShopID     int32  `json:"shop_id"`
	Email      string `json:"email,omitempty"`
	Phone      string `json:"phone,omitempty"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Address    string `json:"address"`
	City       string `json:"city"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
	IsActive   bool   `json:"is_active"`
	Shopname   string `json:"shop_name"`
}

type ListUsersResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int32          `json:"total_items"`
	Page       int32          `json:"current_page"`
	PageSize   int32          `json:"page_size"`
	TotalPages int32          `json:"total_pages"`
}

type LogoutRequest struct {
	UserID string `json:"user_id"`
}

type OTPData struct {
	PhoneNumber string
	OTP         string
	ExpiresAt   time.Time
}

type SendOTPRequest struct {
	Code  string `json:"code"`
	Phone string `json:"phone"`
}

type VerifyOTP struct {
	Otp string `json:"otp"`
}

type PhoneResponse struct {
	Phone_number string `json:"phone_number"`
	Status       bool   `json:"status"`
	Otp          string `json:"otp"`
	Remarks      string `json:"remarks"`
	UserID       string `json:"user_id"`
}

type VerifyOTPResponse struct {
	Token string   `json:"token"`
	Role  []string `json:"role"`
}
