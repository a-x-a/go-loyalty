package userservice

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/a-x-a/go-loyalty/internal/model"
	"github.com/a-x-a/go-loyalty/internal/util"
)

type (
	UserStorage interface {
		AddUser(ctx context.Context, login, pwdHash string) error
		GetUserID(ctx context.Context, login, pwdHash string) (int64, error)
	}

	UserService struct {
		storage UserStorage
	}
)

const secret = "secret"

var ErrNotRegisteredUser = errors.New("Пользователь не зарегистрирован")

func New(storage UserStorage) *UserService {
	return &UserService{
		storage: storage,
	}
}

// Регистрация нового пользователя.
func (s *UserService) Register(ctx context.Context, login, password string) error {
	_, err := model.NewUser(login, password)
	if err != nil {
		return err
	}

	// Генерируем хэш пароля.
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ErrNotRegisteredUser
	}

	// Сохраняем пользователя в БД.
	err = s.storage.AddUser(ctx, login, string(pwdHash))
	if err != nil {
		return ErrNotRegisteredUser
	}

	return nil
}

// Авторизация пользователя.
func (s *UserService) Login(ctx context.Context, login, password string) (string, error) {
	_, err := model.NewUser(login, password)
	if err != nil {
		return "", err
	}

	// Генерируем хэш пароля.
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrNotRegisteredUser
	}

	userID, err := s.storage.GetUserID(ctx, login, string(pwdHash))
	if err != nil {
		return "", ErrNotRegisteredUser
	}

	// Сгенерировать JWT.
	token, err := util.NewToken(userID, secret)
	if err != nil {
		return "", ErrNotRegisteredUser
	}

	return token, nil
}
