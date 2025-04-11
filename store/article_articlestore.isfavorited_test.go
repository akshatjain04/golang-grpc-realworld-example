package store

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

func TestArticleStoreIsFavorited(t *testing.T) {
	type args struct {
		a *model.Article
		u *model.User
	}

	// Initialization of arguments and other variables to be used in test scenarios
	var user = &model.User{ID: 2}
	var article = &model.Article{ID: 1}
	testArgs := []args{
		{a: article, u: user}, // Scenario 1
		{a: nil, u: user},     // Scenario 2
		{a: article, u: user}, // Scenario 3
		{a: article, u: user}, // Scenario 4
	}

	tests := []struct {
		name      string
		args      args
		count     int
		mockError bool
		want      bool
		wantErr   bool
	}{
		{"Favorite present", testArgs[0], 1, false, true, false},
		{"Article or User not available", testArgs[1], 1, false, false, false},
		{"Error from database", testArgs[2], 1, true, false, true},
		{"Article not favorited by User", testArgs[3], 0, false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Starting a Panic recover function
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered, failing test. %v\n", r)
					t.Fail()
				}
			}()

			// Mock database initialization
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("Failed to open gorm database: %s", err)
			}

			mock.ExpectBegin()
			mock.ExpectQuery("SELECT count").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(tt.count))

			if tt.mockError {
				mock.ExpectQuery("SELECT count").WillReturnError(fmt.Errorf("mock db error"))
			}

			// Initializing store with mocked database
			store := &ArticleStore{db: gormDB}

			got, err := store.IsFavorited(tt.args.a, tt.args.u)

			// Output the error logs with detailed diagnostic clarity for better understanding.
			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.IsFavorited() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ArticleStore.IsFavorited() = %v, want %v", got, tt.want)
			}
		})
	}
}
