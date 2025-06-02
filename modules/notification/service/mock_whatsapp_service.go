package service

type MockWhatsAppService struct {
	accessToken   string
	phoneNumberID string
}

func NewMockWhatsAppService() *MockWhatsAppService {
	return &MockWhatsAppService{
		accessToken:   "MOCK_ACCESS_TOKEN",
		phoneNumberID: "MOCK_PHONE_NUMBER_ID",
	}
}

func (s *MockWhatsAppService) SendOTP(phoneNumber string, otp string) error {
	// Mock implementation that always succeeds
	return nil
}
