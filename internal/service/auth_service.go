package service

import (
	"errors"

	"github.com/AronditFire/todo-app/entity"
	"github.com/AronditFire/todo-app/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReg.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Could not hash user password")
	}

	userReg.Password = string(hashedPassword)

	return s.repo.CreateUser(userReg)
}

func (s *AuthService) GetUser(username string) (entity.User, error) {
	return s.repo.GetUser(username)
}

func (s *AuthService) LoginUser(userLogin entity.UserAuthRequest) (entity.User, error) {
	if (len(userLogin.Username) < 3) || (len(userLogin.Username) > 50) {
		return entity.User{}, errors.New("bad username length")
	}
	user, err := s.GetUser(userLogin.Username)
	if err != nil {
		return entity.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password)); err != nil {
		return entity.User{}, errors.New("Incorrect password")
	}

	return s.repo.GetUser(user.Username)
}
