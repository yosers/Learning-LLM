package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/big"
	db "shofy/db/sqlc"
	"shofy/utils/jwt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService interface {
	GetUserPhoneNumber(ctx context.Context, phoneNumber pgtype.Text) (*PhoneNumberResponse, error)
	VerifyOTP(ctx context.Context, otp string, userId int) (*VerifyOTPResponse, error)
}

func NewUserService(dbPool *pgxpool.Pool) UserService {
	return &userService{
		queries: db.New(dbPool),
	}
}

type userService struct {
	queries *db.Queries
}

type PhoneNumberResponse struct {
	Phone_number string `json:"phone_number"`
	Status       bool   `json:"status"`
	Otp          string `json:"otp"`
}

type VerifyOTPResponse struct {
	Token string   `json:"token"`
	User  *db.User `json:"user"`
}

func (s *userService) GetUserPhoneNumber(ctx context.Context, phoneNumber pgtype.Text) (*PhoneNumberResponse, error) {
	items, err := s.queries.GetPhoneNumber(ctx, phoneNumber)
	log.Println("log ini pasti muncul")

	if err != nil {
		log.Println("Error fetching phone number:", err)
		// Kalau tidak ditemukan
		if errors.Is(err, sql.ErrNoRows) {
			return &PhoneNumberResponse{
				Phone_number: "",
				Status:       false,
				Otp:          "",
			}, nil
		}
		// Kalau error selain not found
		return nil, err
	}

	var resultPhone string
	var flag bool

	if items.Phone.Valid && items.Phone.String != "" {
		resultPhone = items.Phone.String
		flag = true
	} else {
		resultPhone = ""
		flag = false
	}

	otp, err := GenerateOTP(6)
	if err != nil {
		// Jika terjadi error saat generate OTP
		log.Println("Error generating OTP:", err)
		log.Fatal(err)
	}

	// Simpan OTP ke database
	err = s.queries.InsertUserLoginOtp(ctx, db.InsertUserLoginOtpParams{
		UserID: items.ID,
		Otp:    otp,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to store OTP: %w", err)
	}

	return &PhoneNumberResponse{
		Phone_number: resultPhone,
		Status:       flag,
		Otp:          otp, // Dummy OTP sementara
	}, nil
}

// GenerateOTP menghasilkan OTP numerik dengan panjang tertentu (4 atau 6 digit)
func GenerateOTP(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid OTP length")
	}

	otp := ""
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10)) // angka 0–9
		if err != nil {
			return "", err
		}
		otp += n.String()
	}
	return otp, nil
}

func (s *userService) VerifyOTP(ctx context.Context, otp string, userId int) (*VerifyOTPResponse, error) {
	otpData, err := s.queries.VerifyOtp(ctx, db.VerifyOtpParams{
		UserID: int32(userId),
		Otp:    otp,
	})

	if err != nil {
		log.Println("Error verifying OTP:", err)
		return nil, fmt.Errorf("failed to verify OTP: %w", err)
	}

	if otpData.Otp == "" {
		log.Println("OTP not found or invalid:", err)

		return nil, errors.New("OTP not found or invalid")
	}

	// Get user data
	user, err := s.queries.GetUser(ctx, otpData.UserID)
	if err != nil {
		log.Println("failed to get user data:", err)
		return nil, fmt.Errorf("failed to get user data: %w", err)
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(otpData.UserID)
	if err != nil {
		log.Println("failed to generate token:", err)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &VerifyOTPResponse{
		Token: token,
		User:  &user,
	}, nil
}
