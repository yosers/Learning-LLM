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
	notificationService "shofy/modules/notification/service"
	model "shofy/modules/users/model"
	"shofy/utils/jwt"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthService struct {
	db              *pgxpool.Pool
	whatsappService *notificationService.WhatsAppService
	otpStore        map[string]*model.OTPData // In-memory store for demo, should use Redis/DB in production
	queries         *db.Queries
}

func NewAuthService(pool *pgxpool.Pool) *AuthService {
	return &AuthService{
		db:              pool,
		whatsappService: notificationService.NewWhatsAppService(),
		otpStore:        make(map[string]*model.OTPData),
		queries:         db.New(pool),
	}
}

func (s *AuthService) GenerateAndSendOTP(ctx context.Context, req model.SendOTPRequest) (*model.PhoneResponse, error) {
	// Generate 6 digit OTP
	otp, err := GenerateOTP(6)

	checkPhone, err := s.queries.FindUserByPhoneAndCode(ctx, db.FindUserByPhoneAndCodeParams{
		Phone:    pgtype.Text{String: req.Phone, Valid: true},
		CodeArea: pgtype.Text{String: req.Code, Valid: true},
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Data Not Found: %w", err)
		}
		return nil, fmt.Errorf("Error Database: %w", err)
	}

	dataUser, err := s.queries.FindUserLoginOtpByPhone(ctx, pgtype.Text{String: req.Phone, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			// Jika tidak ada OTP sebelumnya, insert OTP baru
			err = s.queries.InsertUserLoginOtp(ctx, db.InsertUserLoginOtpParams{
				UserID: checkPhone.ID,
				Otp:    otp,
			})
			if err != nil {
				return nil, fmt.Errorf("Failed to store OTP: %w", err)
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
			return &model.PhoneResponse{
				Phone_number: checkPhone.Phone.String,
				Status:       true,
				Remarks:      "User already logged",
			}, nil
		}
	}
	log.Println("WHATSHAP", err)

	// Send OTP via WhatsApp
	// err = s.whatsappService.SendOTP(checkPhone.CodeArea.String+checkPhone.Phone.String, otp)
	// log.Println("WHATSHAP 4")

	// if err != nil {
	// 	return nil, fmt.Errorf("failed to send OTP: %v", err)
	// }

	return &model.PhoneResponse{
		Phone_number: checkPhone.Phone.String,
		Status:       true,
		Otp:          otp,
	}, nil
}

func (s *AuthService) VerifyOTP(ctx context.Context, inputOTP string) (*model.VerifyOTPResponse, error) {

	count, err := s.queries.CountValidOtps(ctx, inputOTP)

	if err != nil {
		log.Println("Error counting OTP:", err)
		return nil, err
	}
	log.Print("Count Valid OTPs: ", count)
	if count > 1 {
		return nil, fmt.Errorf("OTP conflict: multiple valid entries found")
	}

	otpData, err := s.queries.VerifyOtp(ctx, inputOTP)

	if otpData.Status == "EXPIRED" {
		return nil, fmt.Errorf("OTP has expired")
	} else if otpData.Status == "USED" {
		return nil, fmt.Errorf("OTP has already been used")
	}

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

	return &model.VerifyOTPResponse{
		Token: token,
		Role:  roleList,
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
