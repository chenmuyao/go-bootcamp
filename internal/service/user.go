package service

import (
	"context"
	"errors"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// {{{ Consts

// }}}
// {{{ Global Varirables

var (
	ErrDuplicatedUser        = repository.ErrDuplicatedUser
	ErrInvalidUserOrPassword = errors.New("wrong email or password")
	ErrInvalidUserID         = errors.New("unknown userID")
)

// }}}
// {{{ Interface

// NOTE:
// 1. For test purpose
// 2. UserServiceV1 V2, etc
// 3. UserServiceVIP ...

//go:generate mockgen -source=./user.go -package=svcmocks -destination=./mocks/user.mock.go
type UserService interface {
	SignUp(ctx context.Context, u domain.User) (domain.User, error)
	Login(ctx context.Context, email string, password string) (domain.User, error)
	EditProfile(ctx context.Context, user *domain.User) error
	GetProfile(ctx context.Context, userID int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	FindOrCreateByGitea(ctx context.Context, info domain.GiteaInfo) (domain.User, error)
}

// }}}
// {{{ Struct

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

// }}}
// {{{ Other structs

// }}}
// {{{ Struct Methods

func (svc *userService) SignUp(ctx context.Context, u domain.User) (domain.User, error) {
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

func (svc *userService) Login(
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

func (svc *userService) EditProfile(
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

func (svc *userService) GetProfile(ctx context.Context, userID int64) (domain.User, error) {
	u, err := svc.repo.FindByID(ctx, userID)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserID
	}
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
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

func (svc *userService) FindOrCreateByGitea(
	ctx context.Context,
	info domain.GiteaInfo,
) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, info.Email)
	if err != repository.ErrUserNotFound {
		return u, err
	}
	// Create a new user
	u, err = svc.repo.Create(ctx, domain.User{
		Email: info.Email,
		Name:  info.Login,
	})
	// system error
	if err != nil && err != repository.ErrDuplicatedUser {
		return domain.User{}, err
	}
	if err == repository.ErrDuplicatedUser {
		// TODO: should query the master database
		return svc.repo.FindByEmail(ctx, info.Email)
	}
	// err == nil
	return u, nil
}

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}
