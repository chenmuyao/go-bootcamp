package repository

import (
	"context"

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
	})
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.userDAOToDomain(u), nil
}

func (repo *UserRepository) userDAOToDomain(u dao.User) domain.User {
	return domain.User{
		ID:       u.ID,
		Password: u.Password,
		Email:    u.Email,
	}
}
