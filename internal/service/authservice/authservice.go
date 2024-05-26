package authservice

import (
	"context"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/pkg/errors"

	"github.com/a-x-a/go-loyalty/internal/config"
	"github.com/a-x-a/go-loyalty/internal/customerrors"
	"github.com/a-x-a/go-loyalty/internal/model"
	"github.com/a-x-a/go-loyalty/internal/storage"
	"github.com/a-x-a/go-loyalty/internal/util"
)

type (
	UserStorage interface {
		Add(ctx context.Context, login, pwdHash string) error
		Get(ctx context.Context, login string) (*storage.DTOUser, error)
	}

	AuthService struct {
		storage UserStorage
		cfg     config.ServiceConfig
		l       *zap.Logger
	}
)

func New(storage UserStorage, cfg config.ServiceConfig, l *zap.Logger) *AuthService {
	return &AuthService{storage, cfg, l}
}

// Регистрация нового пользователя.
func (s *AuthService) Register(ctx context.Context, login, password string) error {
	_, err := model.NewUser(login, password)
	if err != nil {
		s.l.Debug("failed to create user", zap.Error(errors.Wrap(err, "model.newuser")))
		return customerrors.ErrInvalidRequestFormat
	}

	// Генерируем хэш пароля.
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.l.Debug("failed to generate hash of password", zap.Error(errors.Wrap(err, "bcrypt.generatefrompassword")))
		return err
	}

	// Сохраняем пользователя в БД.
	err = s.storage.Add(ctx, login, string(pwdHash))
	if err != nil {
		s.l.Debug("failed to add user", zap.Error(errors.Wrap(err, "storage.adduser")))
		return customerrors.ErrUsernameAlreadyTaken
	}

	s.l.Info("user created", zap.String("successful", "authservice.register"))

	return nil
}

// Авторизация пользователя.
func (s *AuthService) Login(ctx context.Context, login, password string) (string, error) {
	_, err := model.NewUser(login, password)
	if err != nil {
		s.l.Debug("failed to create user", zap.Error(errors.Wrap(err, "model.newuser")))
		return "", customerrors.ErrInvalidRequestFormat
	}

	user, err := s.storage.Get(ctx, login)
	if err != nil {
		s.l.Debug("failed to get user", zap.Error(errors.Wrap(err, "storage.getuser")))
		return "", customerrors.ErrInvalidUsernameOrPassword
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		s.l.Debug("invalid login or password", zap.Error(errors.Wrap(err, "bcrypt.comparehashandpassword")))
		return "", customerrors.ErrInvalidUsernameOrPassword
	}

	// Сгенерировать JWT.
	token, err := util.NewToken(user.ID, s.cfg.Secret)
	if err != nil {
		s.l.Debug("failed to generate JWT", zap.Error(errors.Wrap(err, "util.newtoken")))
		return "", err
	}

	s.l.Info("login user", zap.String("successful", "authservice.login"))

	return token, nil
}
