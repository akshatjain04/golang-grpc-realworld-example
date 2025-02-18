package store

import (
	"errors"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

type mockDB struct {
	countResult int
	countError  error
	*gorm.DB
}

/*
ROOST_METHOD_HASH=IsFavorited_799826fee5
ROOST_METHOD_SIG_HASH=IsFavorited_f6d5e67492

FUNCTION_DEF=func (s *ArticleStore) IsFavorited(a *model.Article, u *model.User) (bool, error) // IsFavorited returns whether the article is favorited by the user
*/
func (m *mockDB) Count(value interface{}) *gorm.DB {
	*value.(*int) = m.countResult
	return &gorm.DB{Error: m.countError}
}

func (m *mockDB) Table(name string) *gorm.DB {
	return &gorm.DB{}
}

func TestArticleStoreIsFavorited(t *testing.T) {
	tests := []struct {
		name        string
		article     *model.Article
		user        *model.User
		countResult int
		countError  error
		want        bool
		wantErr     bool
	}{
		{
			name:        "Article is favorited by the user",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        &model.User{Model: gorm.Model{ID: 1}},
			countResult: 1,
			want:        true,
			wantErr:     false,
		},
		{
			name:        "Article is not favorited by the user",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        &model.User{Model: gorm.Model{ID: 1}},
			countResult: 0,
			want:        false,
			wantErr:     false,
		},
		{
			name:    "Nil Article parameter",
			article: nil,
			user:    &model.User{Model: gorm.Model{ID: 1}},
			want:    false,
			wantErr: false,
		},
		{
			name:    "Nil User parameter",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			user:    nil,
			want:    false,
			wantErr: false,
		},
		{
			name:       "Database error",
			article:    &model.Article{Model: gorm.Model{ID: 1}},
			user:       &model.User{Model: gorm.Model{ID: 1}},
			countError: errors.New("database error"),
			want:       false,
			wantErr:    true,
		},
		{
			name:        "Zero count for non-favorited article",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        &model.User{Model: gorm.Model{ID: 1}},
			countResult: 0,
			want:        false,
			wantErr:     false,
		},
		{
			name:        "Multiple favorites for the same article and user",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        &model.User{Model: gorm.Model{ID: 1}},
			countResult: 2,
			want:        true,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{
				countResult: tt.countResult,
				countError:  tt.countError,
				DB:          &gorm.DB{},
			}
			store := &ArticleStore{db: mockDB.DB}
			got, err := store.IsFavorited(tt.article, tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.IsFavorited() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ArticleStore.IsFavorited() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (m *mockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return &gorm.DB{}
}
