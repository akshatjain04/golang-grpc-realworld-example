package store

import (
	"errors"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

type MockArticleStore struct {
	db *mockDB
}
type mockDB struct {
	countResult int
	countError  error
}

/*
ROOST_METHOD_HASH=ArticleStore_IsFavorited_799826fee5
ROOST_METHOD_SIG_HASH=ArticleStore_IsFavorited_f6d5e67492

FUNCTION_DEF=func (s *ArticleStore) IsFavorited(a *model.Article, u *model.User) (bool, error) // IsFavorited returns whether the article is favorited by the user
*/
func (s *MockArticleStore) IsFavorited(a *model.Article, u *model.User) (bool, error) {
	if a == nil || u == nil {
		return false, nil
	}
	var count int
	err := s.db.Table("favorite_articles").Where("article_id = ? AND user_id = ?", a.ID, u.ID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func TestArticleStoreIsFavorited(t *testing.T) {
	tests := []struct {
		name            string
		article         *model.Article
		user            *model.User
		mockCountResult int
		mockCountError  error
		expectedResult  bool
		expectedError   error
	}{
		{
			name:            "Article is favorited by the user",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 1,
			mockCountError:  nil,
			expectedResult:  true,
			expectedError:   nil,
		},
		{
			name:            "Article is not favorited by the user",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 0,
			mockCountError:  nil,
			expectedResult:  false,
			expectedError:   nil,
		},
		{
			name:            "Database error occurs",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 0,
			mockCountError:  errors.New("database error"),
			expectedResult:  false,
			expectedError:   errors.New("database error"),
		},
		{
			name:            "Nil article parameter",
			article:         nil,
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 0,
			mockCountError:  nil,
			expectedResult:  false,
			expectedError:   nil,
		},
		{
			name:            "Nil user parameter",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            nil,
			mockCountResult: 0,
			mockCountError:  nil,
			expectedResult:  false,
			expectedError:   nil,
		},
		{
			name:            "Both article and user parameters are nil",
			article:         nil,
			user:            nil,
			mockCountResult: 0,
			mockCountError:  nil,
			expectedResult:  false,
			expectedError:   nil,
		},
		{
			name:            "Multiple favorites for the same article and user",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 2,
			mockCountError:  nil,
			expectedResult:  true,
			expectedError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{
				countResult: tt.mockCountResult,
				countError:  tt.mockCountError,
			}

			store := &MockArticleStore{
				db: mockDB,
			}

			result, err := store.IsFavorited(tt.article, tt.user)

			if result != tt.expectedResult {
				t.Errorf("Expected result %v, but got %v", tt.expectedResult, result)
			}

			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("Expected error %v, but got %v", tt.expectedError, err)
			}
		})
	}
}

func (m *mockDB) Count(value interface{}) *gorm.DB {
	*value.(*int) = m.countResult
	return &gorm.DB{Error: m.countError}
}

func (m *mockDB) Table(name string) *gorm.DB {
	return &gorm.DB{Error: m.countError}
}

func (m *mockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return &gorm.DB{Error: m.countError}
}
