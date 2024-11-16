package service

import (
	"context"
	"errors"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// {{{ Errors

var (
	ErrDuplicatedEmail       = repository.ErrDuplicatedEmail
	ErrInvalidUserOrPassword = errors.New("wrong email or password")
	ErrInvalidUserID         = errors.New("unknown userID")
)

// }}}

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(
	ctx context.Context,
	email string,
	password string,
) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// check the password
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *UserService) EditProfile(
	ctx context.Context,
	user *domain.User,
) error {
	err := svc.repo.UpdateProfile(ctx, user)
	if err == repository.ErrUserNotFound {
		return ErrInvalidUserID
	}
	if err != nil {
		return err
	}

	return nil
}

func (svc *UserService) GetProfile(ctx context.Context, userID int64) (domain.User, error) {
	u, err := svc.repo.GetProfile(ctx, userID)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserID
	}
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}
