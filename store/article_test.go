// ********RoostGPT********
/*

roost_feedback [2/13/2025, 4:22:32 PM]:a
*/

// ********RoostGPT********

package store

import (
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

func TestArticleStoreAddFavorite(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error %v occurred when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("an error %v occurred when opening gorm database", err)
	}

	store := &ArticleStore{db: gormDB}

	type testdata struct {
		Article  *model.Article
		User     *model.User
		Err      error
		Scenario string
	}

	tests := []testdata{
		{
			Article:  &model.Article{Model: gorm.Model{ID: 1}},
			User:     &model.User{Model: gorm.Model{ID: 1}},
			Err:      nil,
			Scenario: "Successfully Adding a Favorite Article to a User's List",
		},
		{
			Article:  &model.Article{Model: gorm.Model{ID: 2}},
			User:     &model.User{Model: gorm.Model{ID: 2}},
			Err:      gorm.ErrRecordNotFound,
			Scenario: "Rollback on Append Failure",
		},
		{
			Article:  &model.Article{Model: gorm.Model{ID: 3}},
			User:     &model.User{Model: gorm.Model{ID: 3}},
			Err:      gorm.ErrRecordNotFound,
			Scenario: "Rollback on Update Failure",
		},
		{
			Article:  &model.Article{Model: gorm.Model{ID: 4}},
			User:     &model.User{Model: gorm.Model{ID: 4}},
			Err:      nil,
			Scenario: "Increment of FavoritesCount in Successful Operation",
		},
	}

	for _, test := range tests {
		t.Run(test.Scenario, func(t *testing.T) {
			type testQuery struct {
				Query string
				Args  []driver.Value
				Rows  *sqlmock.Rows
			}

			testQueries := []testQuery{
				{
					Query: "SELECT * FROM articles WHERE id = \\?",
					Args:  []driver.Value{test.Article.ID},
					Rows:  sqlmock.NewRows([]string{"id"}).AddRow(test.Article.ID),
				},
				{
					Query: "SELECT * FROM users WHERE id = \\?",
					Args:  []driver.Value{test.User.ID},
					Rows:  sqlmock.NewRows([]string{"id"}).AddRow(test.User.ID),
				},
				{
					Query: "INSERT INTO favorites (user_id, article_id) VALUES (?, ?)",
					Args:  []driver.Value{test.User.ID, test.Article.ID},
				},
				{
					Query: "UPDATE articles SET favorites_count = favorites_count + 1 WHERE id = ?",
					Args:  []driver.Value{test.Article.ID},
				},
			}

			for _, q := range testQueries {
				mock.ExpectQuery(q.Query).WithArgs(q.Args...).WillReturnRows(q.Rows)
			}

			err = store.AddFavorite(test.Article, test.User)
			if !errors.Is(err, test.Err) {
				t.Errorf("expected %v, but got %v", test.Err, err)
			}
			if err == nil && test.Article.FavoritesCount != 1 {
				t.Errorf("expected favorites count to be equal 1, but got %d", test.Article.FavoritesCount)
			}
		})
	}
}
