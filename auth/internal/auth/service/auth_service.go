package service

import (
	"casino/auth/internal/auth/repository"
	"casino/pkg/middleware/auth"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) generateJWT(userID int, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(auth.JwtSecret)
}

func (s *AuthService) Register(name, password string) (string, error) {

	if name == "" || password == "" {
		return "", errors.New("incorrect data")
	}

	candidate, err := s.repo.FindByName(name)
	if err != nil {
		return "", err
	}
	if candidate != nil {
		return "", errors.New("user already exists")
	}

	user, err := s.repo.CreateUser(name, password)
	if err != nil {
		return "", err
	}

	token, err := s.generateJWT(int(user.ID), user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Login(name, password string) (string, error) {

	if name == "" || password == "" {
		return "", errors.New("incorrect data")
	}

	user, err := s.repo.FindByName(name)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := s.generateJWT(int(user.ID), user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}
