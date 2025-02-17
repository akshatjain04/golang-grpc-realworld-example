package store

import (
	"errors"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

// mockDB implements the necessary methods of gorm.DB for testing
type mockDB struct {
	countResult int
	countError  error
	*gorm.DB
}

func (m *mockDB) Table(name string) *gorm.DB {
	return &gorm.DB{}
}

func (m *mockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return &gorm.DB{}
}

func (m *mockDB) Count(value interface{}) *gorm.DB {
	if ptr, ok := value.(*int); ok {
		*ptr = m.countResult
	}
	return &gorm.DB{Error: m.countError}
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
			name:            "Nil Article provided",
			article:         nil,
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 0,
			mockCountError:  nil,
			expectedResult:  false,
			expectedError:   nil,
		},
		{
			name:            "Nil User provided",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            nil,
			mockCountResult: 0,
			mockCountError:  nil,
			expectedResult:  false,
			expectedError:   nil,
		},
		{
			name:            "Both Article and User are nil",
			article:         nil,
			user:            nil,
			mockCountResult: 0,
			mockCountError:  nil,
			expectedResult:  false,
			expectedError:   nil,
		},
		{
			name:            "Multiple favorites for the same Article and User",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 3,
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

			store := &ArticleStore{db: mockDB.DB}

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
