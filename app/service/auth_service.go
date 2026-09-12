package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"
)

const refreshTokenBytes = 32

type AuthService struct {
	users      repository.UserRepository
	tokens     repository.TokenRepository
	jwt        *helper.JWTManager
	refreshTTL time.Duration
}

func NewAuthService(
	users repository.UserRepository,
	tokens repository.TokenRepository,
	jwtManager *helper.JWTManager,
	refreshTTL time.Duration,
) *AuthService {
	return &AuthService{users: users, tokens: tokens, jwt: jwtManager, refreshTTL: refreshTTL}
}

func (s *AuthService) Register(ctx context.Context, req model.RegisterRequest) (model.User, error) {
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	if errs := ValidateRegister(req); len(errs) > 0 {
		return model.User{}, newValidationError(errs)
	}

	hashed, err := helper.HashPassword(req.Password)
	if err != nil {
		return model.User{}, err
	}

	return s.users.Create(ctx, model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashed,
		Role:     "user", // selalu ditentukan server
		IsActive: true,
	})
}

var ErrInvalidCredentials = errors.New("username atau password salah")
var ErrAccountDisabled = errors.New("akun dinonaktifkan")

func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (model.TokenPair, error) {
	if errs := ValidateLogin(req); len(errs) > 0 {
		return model.TokenPair{}, newValidationError(errs)
	}

	user, err := s.users.FindByUsername(ctx, strings.TrimSpace(req.Username))
	if err != nil {
		// Username tidak ada: tetap jalankan verifikasi palsu agar waktu
		// tanggapnya mirip kasus password salah (mencegah user enumeration).
		helper.VerifyDummyPassword(req.Password)
		return model.TokenPair{}, ErrInvalidCredentials
	}

	if !helper.VerifyPassword(user.Password, req.Password) {
		return model.TokenPair{}, ErrInvalidCredentials
	}
	if !user.IsActive {
		return model.TokenPair{}, ErrAccountDisabled
	}

	return s.issueTokenPair(ctx, user)
}

func (s *AuthService) Refresh(ctx context.Context, req model.RefreshRequest) (model.TokenPair, error) {
	if strings.TrimSpace(req.RefreshToken) == "" {
		return model.TokenPair{}, newValidationError(map[string]string{"refresh_token": "wajib diisi"})
	}

	hash := helper.SHA256Hex(req.RefreshToken)
	stored, err := s.tokens.FindActive(ctx, hash)
	if err != nil {
		return model.TokenPair{}, ErrInvalidCredentials
	}

	user, err := s.users.FindByID(ctx, stored.UserID)
	if err != nil || !user.IsActive {
		return model.TokenPair{}, ErrAccountDisabled
	}

	// Rotasi: token lama langsung dicabut, diganti yang baru.
	if err := s.tokens.Revoke(ctx, hash); err != nil {
		return model.TokenPair{}, err
	}

	return s.issueTokenPair(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, req model.RefreshRequest) error {
	if strings.TrimSpace(req.RefreshToken) != "" {
		_ = s.tokens.Revoke(ctx, helper.SHA256Hex(req.RefreshToken))
	}
	return nil
}

func (s *AuthService) Me(ctx context.Context, userID int) (model.User, error) {
	return s.users.FindByID(ctx, userID)
}

func (s *AuthService) issueTokenPair(ctx context.Context, user model.User) (model.TokenPair, error) {
	accessToken, err := s.jwt.GenerateAccess(user)
	if err != nil {
		return model.TokenPair{}, err
	}

	refreshToken, err := helper.RandomToken(refreshTokenBytes)
	if err != nil {
		return model.TokenPair{}, err
	}

	err = s.tokens.Save(ctx, model.RefreshToken{
		UserID:    user.ID,
		TokenHash: helper.SHA256Hex(refreshToken),
		ExpiresAt: time.Now().Add(s.refreshTTL),
	})
	if err != nil {
		return model.TokenPair{}, err
	}

	return model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.jwt.AccessTTL().Seconds()),
	}, nil
}