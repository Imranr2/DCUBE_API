package user

import (
	"errors"
	"net/http"

	dcubeerrs "github.com/Imranr2/DCUBE_API/internal/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserManager interface {
	Register(Request) (*Response, dcubeerrs.Error)
	Login(Request) (*Response, dcubeerrs.Error)
}

type UserManagerImpl struct {
	database *gorm.DB
}

func NewUserManager(database *gorm.DB) UserManager {
	return &UserManagerImpl{
		database: database,
	}
}

func (m *UserManagerImpl) Register(req Request) (*Response, dcubeerrs.Error) {
	var user User
	err := m.database.First(&user, User{Username: req.Username}).Error
	
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dcubeerrs.New(http.StatusInternalServerError, "An error occured while creating new user")
		}
	}

	if user.Username == req.Username {
		return nil, dcubeerrs.New(http.StatusBadRequest, "Username already exists")
	}

	pwHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, dcubeerrs.New(http.StatusInternalServerError, "An error occured while hashing password")
	}

	newUser := User{
		Username: req.Username,
		Password: string(pwHash),
	}

	err = m.database.Create(&newUser).Error

	if err != nil {
		return nil, dcubeerrs.New(http.StatusInternalServerError, "An error occured while creating new user")
	}

	return &Response{User: newUser}, nil
}

func (m *UserManagerImpl) Login(req Request) (*Response, dcubeerrs.Error) {
	var user User
	err := m.database.First(&user, User{Username: req.Username}).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dcubeerrs.New(http.StatusUnauthorized, "Invalid credentials")
		}
		return nil, dcubeerrs.New(http.StatusInternalServerError, "An error occured while authenticating user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))

	if err != nil {
		return nil, dcubeerrs.New(http.StatusUnauthorized, "Invalid credentials")
	}

	return &Response{User: user}, nil
}