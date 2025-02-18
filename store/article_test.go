package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
	*gorm.DB
}

/*
ROOST_METHOD_HASH=IsFavorited_799826fee5
ROOST_METHOD_SIG_HASH=IsFavorited_f6d5e67492

FUNCTION_DEF=func (s *ArticleStore) IsFavorited(a *model.Article, u *model.User) (bool, error) // IsFavorited returns whether the article is favorited by the user
*/
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
			name:        "Non-zero count for favorited article",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        &model.User{Model: gorm.Model{ID: 1}},
			countResult: 2,
			want:        true,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{}
			setupMockDB(mockDB)
			store := &ArticleStore{db: mockDB.DB}

			mockDB.On("Table", "favorites").Return(mockDB.DB)
			mockDB.On("Where", "article_id = ? AND user_id = ?", tt.article.ID, tt.user.ID).Return(mockDB.DB)
			mockDB.On("Count", mock.AnythingOfType("*int")).Run(func(args mock.Arguments) {
				arg := args.Get(0).(*int)
				*arg = tt.countResult
			}).Return(mockDB.DB)

			got, err := store.IsFavorited(tt.article, tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.IsFavorited() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ArticleStore.IsFavorited() = %v, want %v", got, tt.want)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func setupMockDB(mockDB *MockDB) {
	db, _, _ := sqlmock.New()
	gdb, _ := gorm.Open("mysql", db)
	mockDB.DB = gdb
}
