package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type WhatsAppService struct {
	accessToken   string
	phoneNumberID string
}

type WhatsAppMessage struct {
	MessagingProduct string   `json:"messaging_product"`
	RecipientType    string   `json:"recipient_type"`
	To               string   `json:"to"`
	Type             string   `json:"type"`
	Template         Template `json:"template"`
}

type Template struct {
	Name       string      `json:"name"`
	Language   Language    `json:"language"`
	Components []Component `json:"components"`
}

type Language struct {
	Code string `json:"code"`
}

type Component struct {
	Type       string      `json:"type"`
	Parameters []Parameter `json:"parameters"`
}

type Parameter struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func NewWhatsAppService() *WhatsAppService {
	return &WhatsAppService{
		accessToken:   os.Getenv("WHATSAPP_ACCESS_TOKEN"),
		phoneNumberID: os.Getenv("WHATSAPP_PHONE_NUMBER_ID"),
	}
}

func (s *WhatsAppService) SendOTP(phoneNumber string, otp string) error {

	log.Println("WHATSHAP 3" + phoneNumber + " " + otp)

	message := WhatsAppMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               phoneNumber,
		Type:             "template",
		Template: Template{
			Name:     "otp_notification",
			Language: Language{Code: "id"},
			Components: []Component{
				{
					Type: "body",
					Parameters: []Parameter{
						{
							Type: "text",
							Text: otp,
						},
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(message)
	log.Println("jsonData:", jsonData)

	if err != nil {
		return fmt.Errorf("error marshaling message: %v", err)
	}

	url := fmt.Sprintf("https://graph.facebook.com/v17.0/%s/messages", s.phoneNumberID)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	defer resp.Body.Close()
	log.Println("Masuk resp.StatusCode:", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error response from WhatsApp API: %d", resp.StatusCode)
	}

	return nil
}
