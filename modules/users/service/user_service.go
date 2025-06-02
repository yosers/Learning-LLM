package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	db "shofy/db/sqlc"
	notificationService "shofy/modules/notification/service"
	"shofy/utils/jwt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OTPData struct {
	PhoneNumber string
	OTP         string
	ExpiresAt   time.Time
}

type AuthService struct {
	db              *pgxpool.Pool
	whatsappService *notificationService.WhatsAppService
	otpStore        map[string]*OTPData // In-memory store for demo, should use Redis/DB in production
	queries         *db.Queries
}

type PhoneResponse struct {
	Phone_number string `json:"phone_number"`
	Status       bool   `json:"status"`
	Otp          string `json:"otp"`
	Remarks      string `json:"remarks"`
	UserID       string `json:"user_id"`
}

func NewAuthService(pool *pgxpool.Pool) *AuthService {
	return &AuthService{
		db:              pool,
		whatsappService: notificationService.NewWhatsAppService(),
		otpStore:        make(map[string]*OTPData),
		queries:         db.New(pool),
	}
}

func (s *AuthService) GenerateAndSendOTP(ctx context.Context, phoneNumber string) (*PhoneResponse, error) {
	// Generate 6 digit OTP
	rand.Seed(time.Now().UnixNano())
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	// Store OTP with 5 minutes expiration
	s.otpStore[phoneNumber] = &OTPData{
		PhoneNumber: phoneNumber,
		OTP:         otp,
		ExpiresAt:   time.Now().Add(5 * time.Minute),
	}

	getPhone, err := s.queries.FindUserByPhone(ctx, pgtype.Text{String: phoneNumber, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to update OTP: %w", err)
	}

	dataUser, err := s.queries.FindUserLoginOtpByPhone(ctx, pgtype.Text{String: phoneNumber, Valid: true})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("Masuk", err)

			// Jika tidak ada OTP sebelumnya, insert OTP baru
			err = s.queries.InsertUserLoginOtp(ctx, db.InsertUserLoginOtpParams{
				UserID: getPhone.ID,
				Otp:    otp,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to store OTP: %w", err)
			}
		}
	} else {

		if dataUser.IsUsed.Valid && !dataUser.IsUsed.Bool {
			// Jika OTP sudah ada, update OTP
			err = s.queries.UpdateOTPByUserId(ctx, db.UpdateOTPByUserIdParams{
				UserID: int32(dataUser.UserID),
				Otp:    otp,
			})

			if err != nil {
				log.Println("Error updating OTP in database:", err)
				return nil, fmt.Errorf("failed to update OTP: %w", err)
			}
		} else {
			return &PhoneResponse{
				Phone_number: phoneNumber,
				Status:       true,
				Remarks:      "User already logged",
			}, nil
		}
	}
	// Send OTP via WhatsApp
	err = s.whatsappService.SendOTP(phoneNumber, otp)

	if err != nil {
		return nil, fmt.Errorf("failed to send OTP: %v", err)
	}

	return &PhoneResponse{
		Phone_number: phoneNumber,
		Status:       true,
		Otp:          otp,
	}, nil
}

func (s *AuthService) VerifyOTP(ctx context.Context, phoneNumber, inputOTP string) (*VerifyOTPResponse, error) {
	otpData, err := s.queries.VerifyOtp(ctx, inputOTP)

	if err != nil {
		log.Println("Failed to verify OTP in User Login OTP:", err)
		return nil, fmt.Errorf("Failed to verify OTP in User Login OTP: %w", err)
	}

	// Mark OTP as used
	err = s.queries.UpdateIsUsed(ctx, db.UpdateIsUsedParams{
		UserID: int32(otpData.UserID),
		Otp:    inputOTP,
	})

	if err != nil {
		log.Println("Error update OTP:", err)
		return nil, fmt.Errorf("failed to update OTP: %w", err)
	}

	//GENERATE ROLE JWT
	rolesFromDB, err := s.queries.ListUserRole(ctx, otpData.UserID)

	if err != nil {
		log.Println("failed to get user roles:", err)
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	fmt.Printf("Role DB: %+v\n", rolesFromDB)

	var roleList []string
	for _, r := range rolesFromDB {
		roleList = append(roleList, r.Name)
	}
	// Generate JWT token
	token, err := jwt.GenerateToken(otpData.UserID, roleList)

	if err != nil {
		log.Println("failed to generate token:", err)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &VerifyOTPResponse{
		Token: token,
	}, nil
}
