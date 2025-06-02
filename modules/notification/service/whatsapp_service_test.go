package service

import (
	"testing"
)

func TestWhatsAppService_SendOTP(t *testing.T) {
	// Create mock service instance
	service := NewMockWhatsAppService()

	tests := []struct {
		name        string
		phoneNumber string
		otp         string
		wantErr     bool
	}{
		{
			name:        "Valid Indonesian number",
			phoneNumber: "081234567890",
			otp:         "123456",
			wantErr:     false,
		},
		{
			name:        "Valid International format",
			phoneNumber: "6281234567890",
			otp:         "123456",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.SendOTP(tt.phoneNumber, tt.otp)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendOTP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
