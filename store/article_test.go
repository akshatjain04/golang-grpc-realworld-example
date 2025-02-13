package store

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

/*
ROOST_METHOD_HASH=AddFavorite_9460fca478
ROOST_METHOD_SIG_HASH=AddFavorite_c13a109f91

FUNCTION_DEF=func (s *ArticleStore) AddFavorite(a *model.Article, u *model.User) error // AddFavorite favorite an article
*/
func TestArticleStoreAddFavorite(t *testing.T) {
	db, _, err := sqlmock.New()
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
			err = store.AddFavorite(test.Article, test.User)
			if err != nil {
				if err != test.Err {
					t.Errorf("expected %v, but got %v", test.Err, err)
				}
				return
			}
			if test.Article.FavoritesCount != 1 {
				t.Errorf("expected favorites count to be equal 1, but got %d", test.Article.FavoritesCount)
			}
		})
	}
}
