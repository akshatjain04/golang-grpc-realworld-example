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

func NewMockDB(countResult int, countError error) *gorm.DB {
	return &gorm.DB{Value: &mockDB{
		countResult: countResult,
		countError:  countError,
	}}
}

func (m *mockDB) Table(name string) *gorm.DB {
	return &gorm.DB{Value: m}
}

func TestArticleStoreIsFavorited(t *testing.T) {
	tests := []struct {
		name        string
		article     *model.Article
		user        *model.User
		mockDB      *gorm.DB
		expected    bool
		expectedErr error
	}{
		{
			name:        "Article is favorited by the user",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        &model.User{Model: gorm.Model{ID: 1}},
			mockDB:      NewMockDB(1, nil),
			expected:    true,
			expectedErr: nil,
		},
		{
			name:        "Article is not favorited by the user",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        &model.User{Model: gorm.Model{ID: 1}},
			mockDB:      NewMockDB(0, nil),
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "Nil Article parameter",
			article:     nil,
			user:        &model.User{Model: gorm.Model{ID: 1}},
			mockDB:      NewMockDB(0, nil),
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "Nil User parameter",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        nil,
			mockDB:      NewMockDB(0, nil),
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "Database error",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        &model.User{Model: gorm.Model{ID: 1}},
			mockDB:      NewMockDB(0, errors.New("database error")),
			expected:    false,
			expectedErr: errors.New("database error"),
		},
		{
			name:        "Empty database table",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        &model.User{Model: gorm.Model{ID: 1}},
			mockDB:      NewMockDB(0, nil),
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "Multiple favorites for the same article-user pair",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        &model.User{Model: gorm.Model{ID: 1}},
			mockDB:      NewMockDB(2, nil),
			expected:    true,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &ArticleStore{
				db: tt.mockDB,
			}

			result, err := store.IsFavorited(tt.article, tt.user)

			if result != tt.expected {
				t.Errorf("Expected result %v, but got %v", tt.expected, result)
			}

			if (err != nil && tt.expectedErr == nil) || (err == nil && tt.expectedErr != nil) || (err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error()) {
				t.Errorf("Expected error %v, but got %v", tt.expectedErr, err)
			}
		})
	}
}

func (m *mockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return &gorm.DB{Value: m}
}
