package model

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserValidate(t *testing.T) {
	// Table-driven tests for various scenarios
	type testCase struct {
		name        string
		user        User
		expectedErr error
	}

	testCases := []testCase{
		// Scenario 1: Valid User
		{
			name: "Valid User",
			user: User{
				Username: "john.doe",
				Email:    "john.doe@example.com",
				Password: "password123",
				Bio:      "Software engineer",
				Image:    "https://example.com/avatar.jpg",
			},
			expectedErr: nil,
		},

		// Scenario 2: Empty Username
		{
			name: "Empty Username",
			user: User{
				Email:    "john.doe@example.com",
				Password: "password123",
				Bio:      "Software engineer",
				Image:    "https://example.com/avatar.jpg",
			},
			expectedErr: errors.New("Username cannot be empty"),
		},

		// Scenario 3: Invalid Email Format
		{
			name: "Invalid Email Format",
			user: User{
				Username: "john.doe",
				Email:    "invalid_email",
				Password: "password123",
				Bio:      "Software engineer",
				Image:    "https://example.com/avatar.jpg",
			},
			expectedErr: errors.New("Email must be a valid email address"),
		},

		// Scenario 4: Empty Password
		{
			name: "Empty Password",
			user: User{
				Username: "john.doe",
				Email:    "john.doe@example.com",
				Bio:      "Software engineer",
				Image:    "https://example.com/avatar.jpg",
			},
			expectedErr: errors.New("Password cannot be empty"),
		},

		// Scenario 5: Password Length Less Than Minimum
		{
			name: "Password Length Less Than Minimum",
			user: User{
				Username: "john.doe",
				Email:    "john.doe@example.com",
				Password: "short",
				Bio:      "Software engineer",
				Image:    "https://example.com/avatar.jpg",
			},
			expectedErr: errors.New("Password must be at least 8 characters long"),
		},

		// Scenario 6: Invalid Username Format
		{
			name: "Invalid Username Format",
			user: User{
				Username: "john.doe!",
				Email:    "john.doe@example.com",
				Password: "password123",
				Bio:      "Software engineer",
				Image:    "https://example.com/avatar.jpg",
			},
			expectedErr: errors.New("Username must only contain alphanumeric characters"),
		},

		// Scenario 7: Duplicate Username
		{
			name: "Duplicate Username",
			user: User{
				Username: "existing_username",
				Email:    "john.doe@example.com",
				Password: "password123",
				Bio:      "Software engineer",
				Image:    "https://example.com/avatar.jpg",
			},
			// TODO: Simulate duplicate username error
			expectedErr: errors.New("Username already exists"),
		},

		// Scenario 8: Duplicate Email
		{
			name: "Duplicate Email",
			user: User{
				Username: "john.doe",
				Email:    "existing_email@example.com",
				Password: "password123",
				Bio:      "Software engineer",
				Image:    "https://example.com/avatar.jpg",
			},
			// TODO: Simulate duplicate email error
			expectedErr: errors.New("Email already exists"),
		},

		// Scenario 9: Internal Server Error
		{
			name: "Internal Server Error",
			user: User{
				Username: "john.doe",
				Email:    "john.doe@example.com",
				Password: "password123",
				Bio:      "Software engineer",
				Image:    "https://example.com/avatar.jpg",
			},
			// TODO: Simulate internal server error
			expectedErr: errors.New("Internal server error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.user.Validate()

			// Log detailed success or failure reasons for diagnostic clarity
			if err == nil {
				t.Log("Test passed successfully")
			} else {
				t.Logf("Test failed with error: %v", err)
			}

			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
