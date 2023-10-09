package service

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPasswordEncrypt(t *testing.T) {
	password := []byte("123456#password")
	encryptd, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	assert.NoError(t, err)
	println(string(encryptd))
	err = bcrypt.CompareHashAndPassword(encryptd, []byte("wrong password"))
	assert.NotNil(t, err)
}
