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
	GenerateOTPByPhone(ctx context.Context, phoneNumber pgtype.Text) (*PhoneNumberResponse, error)
	VerifyOTP(ctx context.Context, otp string, userId int) (*VerifyOTPResponse, error)
	Logout(ctx context.Context, token string, userId int) error
	CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error)
	UpdateUser(ctx context.Context, userId int32, req *UpdateUserRequest) (*UserResponse, error)
	ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error)
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
	Remarks      string `json:"remarks"`
	UserID       string `json:"user_id"`
}

type VerifyOTPResponse struct {
	Token string `json:"token"`
	//User  *db.User `json:"user"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

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
	PageSize int32 `json:"page_size"`
	ShopID   int32 `json:"shop_id"`
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
}

type ListUsersResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int32          `json:"total"`
	Page       int32          `json:"page"`
	PageSize   int32          `json:"page_size"`
	TotalPages int32          `json:"total_pages"`
}

func (s *userService) GenerateOTPByPhone(ctx context.Context, phoneNumber pgtype.Text) (*PhoneNumberResponse, error) {
	items, err := s.queries.FindUserByPhone(ctx, phoneNumber)

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

	// Simpan atau update OTP ke database
	dataUser, err := s.queries.FindUserLoginOtpByPhone(ctx, phoneNumber)
	log.Println("Error Simpan OTP ke database:", dataUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("Masuk", err)

			// Jika tidak ada OTP sebelumnya, insert OTP baru
			err = s.queries.InsertUserLoginOtp(ctx, db.InsertUserLoginOtpParams{
				UserID: items.ID,
				Otp:    otp,
			})
			if err != nil {

				log.Println("Error inserting OTP to database:", items.ID, otp, err)

				log.Println("Error inserting OTP to database:", err)
				return nil, fmt.Errorf("failed to store OTP: %w", err)
			}
		}
	} else {
		log.Println("Nilai Is Used:", dataUser.IsUsed.Bool)
		log.Println("Nilai Otp:", otp)
		log.Println("Nilai UserID:", dataUser.UserID)

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
			return &PhoneNumberResponse{
				Phone_number: resultPhone,
				Status:       true,
				Remarks:      "User already logged",
			}, nil
		}
	}

	if err != nil {
		log.Println("Error Simpan OTP ke database:", err)
		return nil, fmt.Errorf("failed to store OTP: %w", err)
	}

	return &PhoneNumberResponse{
		Phone_number: resultPhone,
		Status:       flag,
		Otp:          otp,
		UserID:       fmt.Sprintf("%d", items.ID),
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
		log.Println("Failed to verify OTP in User Login OTP:", err)
		return nil, fmt.Errorf("Failed to verify OTP in User Login OTP: %w", err)
	}

	if otpData.Otp == "" {
		log.Println("OTP not found or invalid:", err)
		return nil, errors.New("OTP not found or invalid")
	}

	// Mark OTP as used
	err = s.queries.UpdateIsUsed(ctx, db.UpdateIsUsedParams{
		UserID: int32(userId),
		Otp:    otp,
	})

	if err != nil {
		log.Println("Error update OTP:", err)
		return nil, fmt.Errorf("failed to update OTP: %w", err)
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
		Email: user.Email.String,
		Phone: user.Phone.String,
	}, nil
}

func (s *userService) Logout(ctx context.Context, token string, userId int) error {
	// Invalidate the token
	if err := jwt.InvalidateToken(token); err != nil {
		log.Println("Error invalidating token:", err)
		return fmt.Errorf("failed to invalidate token: %w", err)
	}

	err := s.queries.UpdateIsUsedFalse(ctx, int32(userId))

	if err != nil {
		log.Println("Error update is active UserLoginOtp:", err)
		return fmt.Errorf("Error update is active UserLoginOtp", err)
	}

	return nil
}

func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	// Create user
	user, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		ShopID: req.ShopID,
		Email: pgtype.Text{
			String: req.Email,
			Valid:  req.Email != "",
		},
		Phone: pgtype.Text{
			String: req.Phone,
			Valid:  req.Phone != "",
		},
		IsActive: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create user profile
	profile, err := s.queries.CreateUserProfile(ctx, db.CreateUserProfileParams{
		UserID: user.ID,
		Phone: pgtype.Text{
			String: req.Phone,
			Valid:  req.Phone != "",
		},
		FirstName: pgtype.Text{
			String: req.FirstName,
			Valid:  req.FirstName != "",
		},
		LastName: pgtype.Text{
			String: req.LastName,
			Valid:  req.LastName != "",
		},
		Address: pgtype.Text{
			String: req.Address,
			Valid:  req.Address != "",
		},
		City: pgtype.Text{
			String: req.City,
			Valid:  req.City != "",
		},
		Country: pgtype.Text{
			String: req.Country,
			Valid:  req.Country != "",
		},
		PostalCode: pgtype.Text{
			String: req.PostalCode,
			Valid:  req.PostalCode != "",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user profile: %w", err)
	}

	return &UserResponse{
		ID:         user.ID,
		ShopID:     user.ShopID,
		Email:      user.Email.String,
		Phone:      user.Phone.String,
		FirstName:  profile.FirstName.String,
		LastName:   profile.LastName.String,
		Address:    profile.Address.String,
		City:       profile.City.String,
		Country:    profile.Country.String,
		PostalCode: profile.PostalCode.String,
		IsActive:   user.IsActive.Bool,
	}, nil
}

func (s *userService) UpdateUser(ctx context.Context, userId int32, req *UpdateUserRequest) (*UserResponse, error) {
	// Get existing user
	user, err := s.queries.GetUser(ctx, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update user profile
	err = s.queries.UpdateUserProfile(ctx, db.UpdateUserProfileParams{
		UserID: userId,
		FirstName: pgtype.Text{
			String: req.FirstName,
			Valid:  req.FirstName != "",
		},
		LastName: pgtype.Text{
			String: req.LastName,
			Valid:  req.LastName != "",
		},
		Address: pgtype.Text{
			String: req.Address,
			Valid:  req.Address != "",
		},
		City: pgtype.Text{
			String: req.City,
			Valid:  req.City != "",
		},
		Country: pgtype.Text{
			String: req.Country,
			Valid:  req.Country != "",
		},
		PostalCode: pgtype.Text{
			String: req.PostalCode,
			Valid:  req.PostalCode != "",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	// Get updated profile
	profile, err := s.queries.GetUserProfile(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated user profile: %w", err)
	}

	return &UserResponse{
		ID:         user.ID,
		ShopID:     user.ShopID,
		Email:      user.Email.String,
		Phone:      user.Phone.String,
		FirstName:  profile.FirstName.String,
		LastName:   profile.LastName.String,
		Address:    profile.Address.String,
		City:       profile.City.String,
		Country:    profile.Country.String,
		PostalCode: profile.PostalCode.String,
		IsActive:   user.IsActive.Bool,
	}, nil
}

func (s *userService) ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error) {
	// Get total count first
	total, err := s.queries.CountUsers(ctx, req.ShopID)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Calculate pagination
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}

	offset := (req.Page - 1) * req.PageSize
	totalPages := int32((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	// Get paginated users
	users, err := s.queries.ListUsers(ctx, db.ListUsersParams{
		ShopID: req.ShopID,
		Limit:  req.PageSize,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Get user profiles for all users
	userResponses := make([]UserResponse, 0, len(users))
	for _, user := range users {
		profile, err := s.queries.GetUserProfile(ctx, user.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to get user profile: %w", err)
		}

		userResponses = append(userResponses, UserResponse{
			ID:         user.ID,
			ShopID:     user.ShopID,
			Email:      user.Email.String,
			Phone:      user.Phone.String,
			FirstName:  profile.FirstName.String,
			LastName:   profile.LastName.String,
			Address:    profile.Address.String,
			City:       profile.City.String,
			Country:    profile.Country.String,
			PostalCode: profile.PostalCode.String,
			IsActive:   user.IsActive.Bool,
		})
	}

	return &ListUsersResponse{
		Users:      userResponses,
		Total:      int32(total),
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}
