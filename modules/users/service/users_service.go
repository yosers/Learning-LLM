package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	db "shofy/db/sqlc"
	model "shofy/modules/users/model"
	"shofy/utils/jwt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService interface {
	Logout(ctx context.Context, token string, userId int) error
	CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error)
	UpdateUser(ctx context.Context, userId int32, req *model.UpdateUserRequest) (*model.UserResponse, error)
	ListUsers(ctx context.Context, req *model.ListUsersRequest) (*model.ListUsersResponse, error)
	GetUserByID(ctx context.Context, id int32) (model.UserResponse, error)
	DeleteUsersByID(ctx context.Context, id int32) error
}

func NewUserService(dbPool *pgxpool.Pool) UserService {
	return &userService{
		queries: db.New(dbPool),
	}
}

type userService struct {
	queries *db.Queries
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

func (s *userService) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error) {

	checkPhone, err := s.queries.FindUserByPhone(ctx, pgtype.Text{String: req.Phone, Valid: true})

	if checkPhone.ID != 0 {
		return nil, fmt.Errorf("Phone already exist")
	}

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

	log.Println("DATA USER:", user.ID, user.ShopID, user.Email.String, user.Phone.String)
	// Check for errors

	if err != nil {
		log.Println("Error creating user:", err)
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
		log.Println("Error creating user:", err)

		return nil, fmt.Errorf("failed to create user profile: %w", err)
	}

	return &model.UserResponse{
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

func (s *userService) UpdateUser(ctx context.Context, userId int32, req *model.UpdateUserRequest) (*model.UserResponse, error) {
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

	return &model.UserResponse{
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

func (s *userService) ListUsers(ctx context.Context, req *model.ListUsersRequest) (*model.ListUsersResponse, error) {
	// Get total count first
	total, err := s.queries.CountUsers(ctx)
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
		Limit:  req.PageSize,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	log.Println("Total users:", users)
	// Get user profiles for all users
	userResponses := make([]model.UserResponse, 0, len(users))
	for _, user := range users {
		profile, err := s.queries.GetUserProfile(ctx, user.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to get user profile: %w", err)
		}

		userResponses = append(userResponses, model.UserResponse{
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
			Shopname:   user.Shopname,
		})
	}

	return &model.ListUsersResponse{
		Users:      userResponses,
		Total:      int32(total),
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *userService) GetUserByID(ctx context.Context, id int32) (model.UserResponse, error) {
	user, err := s.queries.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.UserResponse{}, fmt.Errorf("user not found")
		}
		return model.UserResponse{}, fmt.Errorf("failed to get user: %w", err)
	}

	profile, err := s.queries.GetUserProfile(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.UserResponse{}, fmt.Errorf("user profile not found")
		}
		return model.UserResponse{}, fmt.Errorf("failed to get user profile: %w", err)
	}

	return model.UserResponse{
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
		Shopname:   user.Name,
	}, nil
}

func (s *userService) DeleteUsersByID(ctx context.Context, id int32) error {
	// Check if user exists
	_, err := s.queries.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Proceed to delete user
	err = s.queries.DeleteUserById(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
