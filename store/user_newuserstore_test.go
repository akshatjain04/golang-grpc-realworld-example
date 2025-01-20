package store

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestNewUserStore(t *testing.T) {
	// TODO: Add test cases for different database connection scenarios (valid, nil, invalid) as described in the instructions.

	// Example test case for a valid database connection:
	t.Run("Valid DB", func(t *testing.T) {
		// Arrange
		db, err := gorm.Open("mysql", "user:password@/database")
		if err != nil {
			t.Fatal(err)
		}

		// Act
		userStore := NewUserStore(db)

		// Assert
		assert.NotNil(t, userStore, "Expected UserStore to be created, but got nil")
		assert.Equal(t, db, userStore.db, "Expected UserStore to use the provided database connection, but got %v", userStore.db)
	})

	// Example test case for a nil database connection:
	t.Run("Nil DB", func(t *testing.T) {
		// Arrange
		var db *gorm.DB

		// Act
		userStore := NewUserStore(db)

		// Assert
		assert.Nil(t, userStore, "Expected UserStore to be nil when using a nil database connection, but got a non-nil value")
	})

	// Example test case for an invalid database connection:
	t.Run("Invalid DB", func(t *testing.T) {
		// Arrange
		db, err := gorm.Open("invalid_dialect", "invalid_connection_string")
		if err != nil {
			t.Fatal(err)
		}

		// Act
		userStore := NewUserStore(db)

		// Assert
		assert.Nil(t, userStore, "Expected UserStore to be nil when using an invalid database connection, but got a non-nil value")
	})
}
