package userservice

import (
	"context"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/pkg/errors"

	"github.com/a-x-a/go-loyalty/internal/config"
	"github.com/a-x-a/go-loyalty/internal/model"
	"github.com/a-x-a/go-loyalty/internal/storage"
	"github.com/a-x-a/go-loyalty/internal/util"
)

type (
	UserStorage interface {
		AddUser(ctx context.Context, login, pwdHash string) error
		GetUser(ctx context.Context, login string) (*storage.DTOUser, error)
	}

	UserService struct {
		storage UserStorage
		cfg     config.ServiceConfig
		l       *zap.Logger
	}
)

func New(storage UserStorage, cfg config.ServiceConfig, l *zap.Logger) *UserService {
	return &UserService{storage, cfg, l}
}

// Регистрация нового пользователя.
func (s *UserService) Register(ctx context.Context, login, password string) error {
	_, err := model.NewUser(login, password)
	if err != nil {
		s.l.Debug("failed to create user", zap.Error(errors.Wrap(err, "model.newuser")))
		return model.ErrInvalidRequestFormat
	}

	// Генерируем хэш пароля.
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.l.Debug("failed to generate hash of password", zap.Error(errors.Wrap(err, "bcrypt.generatefrompassword")))
		return err
	}

	// Сохраняем пользователя в БД.
	err = s.storage.AddUser(ctx, login, string(pwdHash))
	if err != nil {
		s.l.Debug("failed to add user", zap.Error(errors.Wrap(err, "storage.adduser")))
		return model.ErrUsernameAlreadyTaken
	}

	s.l.Info("user created", zap.String("succesful", "userservice.register"))

	return nil
}

// Авторизация пользователя.
func (s *UserService) Login(ctx context.Context, login, password string) (string, error) {
	_, err := model.NewUser(login, password)
	if err != nil {
		s.l.Debug("failed to create user", zap.Error(errors.Wrap(err, "model.newuser")))
		return "", model.ErrInvalidRequestFormat
	}

	user, err := s.storage.GetUser(ctx, login)
	if err != nil {
		s.l.Debug("failed to get user", zap.Error(errors.Wrap(err, "storage.getuser")))
		return "", model.ErrInvalidUsernameOrPassword
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		s.l.Debug("invalid login or password", zap.Error(errors.Wrap(err, "bcrypt.comparehashandpassword")))
		return "", model.ErrInvalidUsernameOrPassword
	}

	// Сгенерировать JWT.
	token, err := util.NewToken(user.ID, s.cfg.Secret)
	if err != nil {
		s.l.Debug("failed to generate JWT", zap.Error(errors.Wrap(err, "util.newtoken")))
		return "", err
	}

	s.l.Info("login user", zap.String("succesful", "userservice.login"))

	return token, nil
}
