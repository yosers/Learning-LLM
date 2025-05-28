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

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService interface {
	GetUserPhoneNumber(ctx context.Context, phoneNumber pgtype.Text) (*PhoneNumberResponse, error)
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
