package services

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPasswordService_ValidatePassword_ShouldReturnError_WhenPasswordEmptyString(t *testing.T) {
	passwordService := NewPasswordService()

	err := passwordService.ValidatePassword("")

	assert.NotNil(t, err)
}

func TestPasswordService_ValidatePassword_ShouldReturnError_WhenLengthLessThanSix(t *testing.T) {
	passwordService := NewPasswordService()

	err := passwordService.ValidatePassword("1234")

	assert.NotNil(t, err)
}

func TestPasswordService_ComparePassword_True_WhenPasswordsEqual(t *testing.T) {
	passwordService := NewPasswordService()

	const password = "Hello12345"
	result, _ := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)

	assert.True(t, passwordService.ComparePassword(password, string(result)))
}

func TestPasswordService_ComparePassword_False_WhenPasswordsNotEqual(t *testing.T) {
	passwordService := NewPasswordService()

	const password = "Hello12345"
	result, _ := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)

	assert.False(t, passwordService.ComparePassword("1234", string(result)))
}

func TestPasswordService_HashId_Same_WhenInputSame(t *testing.T) {
	passwordService := NewPasswordService()

	id1 := passwordService.HashId("123456")

	id2 := passwordService.HashId("123456")

	assert.Equal(t, id1, id2)
}

