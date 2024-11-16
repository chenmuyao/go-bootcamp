package repository

import (
	"context"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
)

var (
	ErrDuplicatedEmail = dao.ErrDuplicatedEmail
	// NOTE: Strongly related to the service
	ErrUserNotFound = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
		Birthday: u.Birthday.UnixMilli(), // NOTE: the zero value of time.Time is not that of int64
	})
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.userDAOToDomain(&u), nil
}

func (repo *UserRepository) UpdateProfile(ctx context.Context, user *domain.User) error {
	err := repo.dao.UpdateProfile(ctx, repo.userDomainToDAO(user))
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) GetProfile(ctx context.Context, userID int64) (domain.User, error) {
	u, err := repo.dao.FindByID(ctx, userID)
	if err != nil {
		return domain.User{}, err
	}
	return repo.userDAOToDomain(&u), nil
}

func (repo *UserRepository) userDAOToDomain(u *dao.User) domain.User {
	return domain.User{
		ID:       u.ID,
		Password: u.Password,
		Email:    u.Email,
		Name:     u.Name,
		Birthday: time.Unix(u.Birthday/1000, u.Birthday%1000*10e6),
		Profile:  u.Profile,
	}
}

func (repo *UserRepository) userDomainToDAO(u *domain.User) dao.User {
	return dao.User{
		ID:       u.ID,
		Name:     u.Name,
		Birthday: u.Birthday.UnixMilli(),
		Profile:  u.Profile,
	}
}
