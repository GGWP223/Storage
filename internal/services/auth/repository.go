package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(request Request) error
	GetUser(request Request) (User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(request Request) error {
	pass, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
		return errors.New("failed to generate password")
	}

	user := User{
		Uid:       uuid.NewString(),
		CreatedAt: time.DateTime,
		Login:     request.Login,
		Password:  string(pass),
	}

	return r.db.Create(&user).Error
}

func (r *repository) GetUser(request Request) (User, error) {
	var user User

	if err := r.db.First(&user, "login = ?", request.Login).Error; err != nil {
		return User{}, errors.New("login not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return User{}, errors.New("wrong password")
	}

	return user, nil
}
