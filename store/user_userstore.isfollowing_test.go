package store

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

// Test data
var mockUserA = &model.User{ID: 1, Username: "UserA", Bio: "Bio A", Image: "http://example.com/imageA.png", Token: "tokenA"}
var mockUserB = &model.User{ID: 2, Username: "UserB", Bio: "Bio B", Image: "http://example.com/imageB.png", Token: "tokenB"}

func TestUserStoreIsFollowing(t *testing.T) {
	tt := []struct {
		name     string
		userA    *model.User
		userB    *model.User
		want     bool
		mockFunc func(mock sqlmock.Sqlmock, count int, err error)
	}{
		{
			name:     "Test for null User parameters",
			userA:    nil,
			userB:    nil,
			want:     false,
			mockFunc: func(mock sqlmock.Sqlmock, count int, err error) {},
		},
		{
			name:  "Test for User A not following User B",
			userA: mockUserA,
			userB: mockUserB,
			want:  false,
			mockFunc: func(mock sqlmock.Sqlmock, count int, err error) {
				mock.ExpectQuery("^SELECT count(*) FROM follows WHERE from_user_id = (.+) AND to_user_id = (.+)$").
					WithArgs(mockUserA.ID, mockUserB.ID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
			},
		},
		{
			name:  "Test for User A following User B",
			userA: mockUserA,
			userB: mockUserB,
			want:  true,
			mockFunc: func(mock sqlmock.Sqlmock, count int, err error) {
				mock.ExpectQuery("^SELECT count(*) FROM follows WHERE from_user_id = (.+) AND to_user_id = (.+)$").
					WithArgs(mockUserA.ID, mockUserB.ID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
			},
		},
		{
			name:  "Test for database errors",
			userA: mockUserA,
			userB: mockUserB,
			want:  false,
			mockFunc: func(mock sqlmock.Sqlmock, count int, err error) {
				mock.ExpectQuery("^SELECT count(*) FROM follows WHERE from_user_id = (.+) AND to_user_id = (.+)$").
					WithArgs(mockUserA.ID, mockUserB.ID).
					WillReturnError(fmt.Errorf("database error"))
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v", r)
					t.Fail()
				}
			}()

			// Open database mock
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			// Set mock function
			tc.mockFunc(mock, 1, nil)
			// Create guards with mock database
			gormDB, _ := gorm.Open("postgres", db)
			guards := &UserStore{db: gormDB}
			// Call IsFollowing
			got, _ := guards.IsFollowing(tc.userA, tc.userB)

			if got != tc.want {
				t.Errorf("IsFollowing() = %v, want %v", got, tc.want)
			}
		})
	}
}
