package service

import (
	"context"
	"errors"
	"testing"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository"
	repomocks "github.com/chenmuyao/go-bootcamp/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordEncrypt(t *testing.T) {
	password := []byte("password123!")
	encrypted, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	assert.NoError(t, err)
	t.Log(string(encrypted))

	err = bcrypt.CompareHashAndPassword(encrypted, []byte("wrong"))
	assert.ErrorIs(t, err, bcrypt.ErrMismatchedHashAndPassword)
	err = bcrypt.CompareHashAndPassword(encrypted, password)
	assert.NoError(t, err)
}

func TestUserLogin(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) repository.UserRepository

		// Input
		ctx      context.Context
		email    string
		password string

		// Output
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "login ok",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(gomock.Any(), "ok@test.com").Return(
					domain.User{
						ID:       123,
						Email:    "ok@test.com",
						Password: "$2a$10$ak/.qMW4bKq3ksEXokuuquyXXNjHAv1t8wqwWvGVbje0rjyrZTqgy",
					}, nil)
				return repo
			},
			email:    "ok@test.com",
			password: "password123!",
			wantUser: domain.User{
				ID:       123,
				Email:    "ok@test.com",
				Password: "$2a$10$ak/.qMW4bKq3ksEXokuuquyXXNjHAv1t8wqwWvGVbje0rjyrZTqgy",
			},
		},
		{
			name: "user not found",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(gomock.Any(), "usernotfound@test.com").Return(
					domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email:    "usernotfound@test.com",
			password: "password123!",
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "db error",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(gomock.Any(), "dberror@test.com").Return(
					domain.User{}, errors.New("db error"))
				return repo
			},
			email:    "dberror@test.com",
			password: "password123!",
			wantErr:  errors.New("db error"),
		},
		{
			name: "wrong password",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(gomock.Any(), "wrong@test.com").Return(
					domain.User{
						ID:       123,
						Email:    "wrong@test.com",
						Password: "$2a$10$ak/.qMW4bKq3ksEXokuuquyXXNjHAv1t8wqwWvGVbje0rjyrZTqgy",
					}, nil)
				return repo
			},
			email:    "wrong@test.com",
			password: "wrong",
			wantErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewUserService(tc.mock(ctrl))

			u, err := svc.Login(tc.ctx, tc.email, tc.password)

			assert.Equal(t, tc.wantUser, u)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
