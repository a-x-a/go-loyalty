package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/a-x-a/go-loyalty/internal/customerrors"
)

func TestNewUser(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		name         string
		login        string
		password     string
		expectedUser *User
		expectedErr  error
	}{
		{
			name: "valid login and password",
			login: "testuser",
			password: "testpassword",
			expectedUser: &User{
				login: "testuser",
				password: "testpassword",
			},
			expectedErr: nil,
		},
		{
			name: "invalid login or password",
			login: "",
			password: "",
			expectedUser: nil,
			expectedErr: customerrors.ErrInvalidUsernameOrPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := NewUser(tc.login, tc.password)
			assert.Equal(tc.expectedUser, user)
			assert.Equal(tc.expectedErr, err)
		})
	}

}
