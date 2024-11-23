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
	ErrDuplicatedUser        = repository.ErrDuplicatedUser
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

func (svc *UserService) SignUp(ctx context.Context, u domain.User) (domain.User, error) {
	// default cost perf test 13.28 req/s
	// Mincost 127.40 req/s
	// cost 12 perf test 3.60 req/s
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
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
	u, err := svc.repo.FindById(ctx, userID)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserID
	}
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (svc *UserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		return u, err
	}
	// Create a new user
	u, err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	// system error
	if err != nil && err != repository.ErrDuplicatedUser {
		return domain.User{}, err
	}
	if err == repository.ErrDuplicatedUser {
		// TODO: should query the master database
		return svc.repo.FindByPhone(ctx, phone)
	}
	// err == nil
	return u, nil
}
