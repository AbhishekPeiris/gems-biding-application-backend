package service

import (
	"errors"
	"time"

	"github.com/boswin/gems-auction-backend/config"
	"github.com/boswin/gems-auction-backend/internal/domain"
	"github.com/boswin/gems-auction-backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

type RegisterRequest struct {
	FullName string          `json:"full_name"`
	Email    string          `json:"email"`
	Password string          `json:"password"`
	Role     domain.UserRole `json:"role"` // ADMIN/SELLER/BUYER
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  domain.User `json:"user"`
}

func (s *AuthService) Register(req RegisterRequest) (*domain.User, error) {
	// check existing email
	existing, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existing != nil && existing.ID != 0 {
		return nil, errors.New("email already registered")
	}

	hashed, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: hashed,
		Role:     req.Role,
		IsActive: true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// do not return password
	user.Password = ""
	return user, nil
}

func (s *AuthService) Login(req LoginRequest) (*AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("user is disabled")
	}

	if err := comparePassword(user.Password, req.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := generateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return &AuthResponse{Token: token, User: *user}, nil
}

func hashPassword(raw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	return string(b), err
}

func comparePassword(hashed string, raw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw))
}

func generateJWT(userID int64, email string, role domain.UserRole) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"role":  string(role),
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(config.AppConfig.JWTSecret))
}
