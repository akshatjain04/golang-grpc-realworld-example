package store

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

func TestArticleStoreGetTags(t *testing.T) {
	// Create an in-memory database for the test.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a gorm database object using the mock database connection.
	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a gorm database object", err)
	}

	// Create an ArticleStore object using the gorm database object.
	articleStore := &ArticleStore{db: gormDB}

	// Define test scenarios and expected results.
	tests := []struct {
		name          string
		mockSetup     func()
		expectedTags  []model.Tag
		expectedError error
	}{
		{
			name: "Retrieve all tags successfully",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "tag1").
					AddRow(2, "tag2")
				mock.ExpectQuery("SELECT (.+) FROM `tags`").WillReturnRows(rows)
			},
			expectedTags: []model.Tag{
				{Model: gorm.Model{ID: 1}, Name: "tag1"},
				{Model: gorm.Model{ID: 2}, Name: "tag2"},
			},
		},
		{
			name: "Handle empty database scenario",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "name"})
				mock.ExpectQuery("SELECT (.+) FROM `tags`").WillReturnRows(rows)
			},
			expectedTags: []model.Tag{},
		},
		{
			name: "Handle database error",
			mockSetup: func() {
				mock.ExpectQuery("SELECT (.+) FROM `tags`").WillReturnError(fmt.Errorf("database error"))
			},
			expectedError: fmt.Errorf("database error"),
		},
	}

	// Run tests.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setup mock database expectations.
			test.mockSetup()

			// Call the GetTags function.
			tags, err := articleStore.GetTags()

			// Validate the results.
			if err != test.expectedError {
				t.Errorf("expected error '%v', got '%v'", test.expectedError, err)
			}
			if !reflect.DeepEqual(tags, test.expectedTags) {
				t.Errorf("expected tags '%v', got '%v'", test.expectedTags, tags)
			}

			// Assert that all expectations were met.
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
