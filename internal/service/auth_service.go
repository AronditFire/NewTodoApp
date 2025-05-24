package service

import (
	"errors"
	"os"
	"time"

	"github.com/AronditFire/todo-app/entity"
	"github.com/AronditFire/todo-app/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	AccessSecret    = []byte(os.Getenv("ACCESS_SECRET"))
	RefreshSecret   = []byte(os.Getenv("REFRESH_SECRET"))
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 24 * time.Hour
)

// Added variable for password hashing so it can be overridden in tests.
var generatePasswordHash = bcrypt.GenerateFromPassword
var CompareHashAndPassword = bcrypt.CompareHashAndPassword

type TokenClaims struct {
	jwt.RegisteredClaims
	UserID  int
	IsAdmin bool
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(userReg entity.UserRegisterRequest) error {
	if (len(userReg.Username) < 3) || (len(userReg.Username) > 50) {
		return errors.New("bad username length")
	}
	// Use the variable instead of calling bcrypt.GenerateFromPassword directly.
	hashedPassword, err := generatePasswordHash([]byte(userReg.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Could not hash user password")
	}

	userReg.Password = string(hashedPassword)

	return s.repo.CreateUser(userReg)
}

func (s *AuthService) GetUser(username string) (entity.User, error) {
	return s.repo.GetUser(username)
}

func (s *AuthService) GetUserByID(id int) (entity.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *AuthService) LoginUser(userLogin entity.UserAuthRequest) (string, string, error) {
	if (len(userLogin.Username) < 3) || (len(userLogin.Username) > 50) {
		return "", "", errors.New("bad username length")
	}
	user, err := s.GetUser(userLogin.Username)
	if err != nil {
		return "", "", err
	}

	if err := CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password)); err != nil {
		return "", "", errors.New("Incorrect password")
	}
	// access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID:  user.ID,
		IsAdmin: user.IsAdmin,
	})

	accessTokenSigned, err := accessToken.SignedString(AccessSecret)
	if err != nil {
		return "", "", errors.New("Could not to sign accessToken")
	}

	// refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID:  user.ID,
		IsAdmin: user.IsAdmin,
	})

	refreshTokenSigned, err := refreshToken.SignedString(RefreshSecret)
	if err != nil {
		return "", "", errors.New("Could not to sign refreshToken")
	}

	// summary
	return accessTokenSigned, refreshTokenSigned, err
}

func (s *AuthService) ParseAccessToken(accessTokenStr string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(accessTokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(AccessSecret), nil
	})
	if err != nil {
		return &TokenClaims{}, err
	}
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return &TokenClaims{}, errors.New("token claims are not of type *tokenClaims")
	}

	return claims, nil
}

func (s *AuthService) ParseRefreshToken(refreshTokenStr string) (int, error) {
	token, err := jwt.ParseWithClaims(refreshTokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(AccessSecret), nil

	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserID, nil
}

func (s *AuthService) RenewTokens(id int) (string, string, error) {
	user, err := s.GetUserByID(id)
	if err != nil {
		return "", "", err
	}

	// access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID:  user.ID,
		IsAdmin: user.IsAdmin,
	})

	accessTokenSigned, err := accessToken.SignedString(AccessSecret)
	if err != nil {
		return "", "", errors.New("Could not to sign accessToken")
	}

	// refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID:  user.ID,
		IsAdmin: user.IsAdmin,
	})

	refreshTokenSigned, err := refreshToken.SignedString(RefreshSecret)
	if err != nil {
		return "", "", errors.New("Could not to sign refreshToken")
	}

	// summary
	return accessTokenSigned, refreshTokenSigned, err
}
