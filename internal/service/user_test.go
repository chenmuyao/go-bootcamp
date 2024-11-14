package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordEncrypt(t *testing.T) {
	password := []byte("123456#hello")
	encrypted, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	assert.NoError(t, err)
	t.Log(string(encrypted))

	err = bcrypt.CompareHashAndPassword(encrypted, []byte("wrong"))
	assert.ErrorIs(t, err, bcrypt.ErrMismatchedHashAndPassword)
	err = bcrypt.CompareHashAndPassword(encrypted, password)
	assert.NoError(t, err)
}
