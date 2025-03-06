package store

import (
	errors "errors"
	testing "testing"

	gorm "github.com/jinzhu/gorm"
	model "github.com/raahii/golang-grpc-realworld-example/model"
)

type MockArticleStore struct {
	db *MockDB
}
type MockDB struct {
	countError error
	count      int
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

func (m *MockDB) Count(value interface{}) *gorm.DB {
	*value.(*int) = m.count
	return &gorm.DB{Error: m.countError}
}

func (m *MockDB) Table(name string) *gorm.DB {
	return &gorm.DB{Error: m.countError}
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return &gorm.DB{Error: m.countError}
}

func TestArticleStoreIsFavorited(t *testing.T) {
	tests := []struct {
		name        string
		article     *model.Article
		user        *model.User
		dbCount     int
		dbError     error
		expected    bool
		expectedErr error
	}{
		{
			name:        "Article is favorited by the user",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        &model.User{Model: gorm.Model{ID: 1}},
			dbCount:     1,
			dbError:     nil,
			expected:    true,
			expectedErr: nil,
		},
		{
			name:        "Article is not favorited by the user",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        &model.User{Model: gorm.Model{ID: 1}},
			dbCount:     0,
			dbError:     nil,
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "Nil Article parameter",
			article:     nil,
			user:        &model.User{Model: gorm.Model{ID: 1}},
			dbCount:     0,
			dbError:     nil,
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "Nil User parameter",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        nil,
			dbCount:     0,
			dbError:     nil,
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "Database error",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        &model.User{Model: gorm.Model{ID: 1}},
			dbCount:     0,
			dbError:     errors.New("database error"),
			expected:    false,
			expectedErr: errors.New("database error"),
		},
		{
			name:        "Multiple favorites for the same article and user",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        &model.User{Model: gorm.Model{ID: 1}},
			dbCount:     3,
			dbError:     nil,
			expected:    true,
			expectedErr: nil,
		},
		{
			name:        "Large number of favorites in the database",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        &model.User{Model: gorm.Model{ID: 1}},
			dbCount:     1000000,
			dbError:     nil,
			expected:    true,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{
				countError: tt.dbError,
				count:      tt.dbCount,
			}

			store := &MockArticleStore{
				db: mockDB,
			}

			result, err := store.IsFavorited(tt.article, tt.user)

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}

			if (err != nil && tt.expectedErr == nil) || (err == nil && tt.expectedErr != nil) || (err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error()) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
