package github.com/raahii/golang-grpc-realworld-example/store

import (
	"errors"
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"time"
)








/*
ROOST_METHOD_HASH=ArticleStore_AddFavorite_9460fca478
ROOST_METHOD_SIG_HASH=ArticleStore_AddFavorite_c13a109f91

FUNCTION_DEF=func (s *ArticleStore) AddFavorite(a *model.Article, u *model.User) error // AddFavorite favorite an article


*/
func (m *mockAssociation) Append(values ...interface{}) error {
	args := m.Called(values...)
	return args.Error(0)
}

func (m *mockDB) Association(column string) *gorm.Association {
	args := m.Called(column)
	return args.Get(0).(*gorm.Association)
}

func (m *mockDB) Begin() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *mockDB) Commit() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *mockDB) Model(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *mockDB) Rollback() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *mockDB) Update(column string, value interface{}) *gorm.DB {
	args := m.Called(column, value)
	return args.Get(0).(*gorm.DB)
}


/*
ROOST_METHOD_HASH=ArticleStore_GetArticles_101b7250e8
ROOST_METHOD_SIG_HASH=ArticleStore_GetArticles_91bc0a6760

FUNCTION_DEF=func (s *ArticleStore) GetArticles(tagName, username string, favoritedBy *model.User, limit, offset int64) ([ // GetArticles get global articles
]model.Article, error) 

*/
func (m *mockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	return m.DB
}

func (m *mockDB) Joins(query string, args ...interface{}) *gorm.DB {
	return m.DB
}

func (m *mockDB) Limit(limit interface{}) *gorm.DB {
	return m.DB
}

func (m *mockDB) Offset(offset interface{}) *gorm.DB {
	return m.DB
}

func (m *mockDB) Preload(column string, conditions ...interface{}) *gorm.DB {
	return m.DB
}

func (m *mockDB) Rows() (*gorm.Rows, error) {
	return nil, nil
}

func (m *mockDB) Select(query interface{}, args ...interface{}) *gorm.DB {
	return m.DB
}

func (m *mockDB) Table(name string) *gorm.DB {
	return m.DB
}

func TestArticleStoreGetArticles(t *testing.T) {
	tests := []struct {
		name        string
		tagName     string
		username    string
		favoritedBy *model.User
		limit       int64
		offset      int64
		mockSetup   func(*mockDB)
		expected    []model.Article
		expectedErr error
	}{
		{
			name:     "Get Articles with No Filters",
			tagName:  "",
			username: "",
			limit:    10,
			offset:   0,
			mockSetup: func(m *mockDB) {
				m.DB.Error = nil

			},
			expected:    []model.Article{},
			expectedErr: nil,
		},
		{
			name:     "Get Articles by Username",
			tagName:  "",
			username: "testuser",
			limit:    10,
			offset:   0,
			mockSetup: func(m *mockDB) {
				m.DB.Error = nil

			},
			expected:    []model.Article{},
			expectedErr: nil,
		},
		{
			name:     "Get Articles by Tag",
			tagName:  "testtag",
			username: "",
			limit:    10,
			offset:   0,
			mockSetup: func(m *mockDB) {
				m.DB.Error = nil

			},
			expected:    []model.Article{},
			expectedErr: nil,
		},
		{
			name:     "Get Favorited Articles",
			tagName:  "",
			username: "",
			favoritedBy: &model.User{
				Model: gorm.Model{ID: 1},
			},
			limit:  10,
			offset: 0,
			mockSetup: func(m *mockDB) {
				m.DB.Error = nil

			},
			expected:    []model.Article{},
			expectedErr: nil,
		},
		{
			name:     "Get Articles with Multiple Filters",
			tagName:  "testtag",
			username: "testuser",
			limit:    10,
			offset:   0,
			mockSetup: func(m *mockDB) {
				m.DB.Error = nil

			},
			expected:    []model.Article{},
			expectedErr: nil,
		},
		{
			name:     "Pagination Test",
			tagName:  "",
			username: "",
			limit:    5,
			offset:   5,
			mockSetup: func(m *mockDB) {
				m.DB.Error = nil

			},
			expected:    []model.Article{},
			expectedErr: nil,
		},
		{
			name:     "Error Handling for Invalid User in FavoritedBy",
			tagName:  "",
			username: "",
			favoritedBy: &model.User{
				Model: gorm.Model{ID: 999},
			},
			limit:  10,
			offset: 0,
			mockSetup: func(m *mockDB) {
				m.DB.Error = gorm.ErrRecordNotFound
			},
			expected:    []model.Article{},
			expectedErr: gorm.ErrRecordNotFound,
		},
		{
			name:     "Empty Result Set",
			tagName:  "nonexistenttag",
			username: "",
			limit:    10,
			offset:   0,
			mockSetup: func(m *mockDB) {
				m.DB.Error = nil

			},
			expected:    []model.Article{},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{&gorm.DB{}}
			tt.mockSetup(mockDB)

			store := &ArticleStore{db: mockDB}
			articles, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, articles)
		})
	}
}

func (m *mockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return m.DB
}


/*
ROOST_METHOD_HASH=ArticleStore_GetFeedArticles_a37e1934b6
ROOST_METHOD_SIG_HASH=ArticleStore_GetFeedArticles_f5f09c020e

FUNCTION_DEF=func (s *ArticleStore) GetFeedArticles(userIDs [ // GetFeedArticles returns following users' articles
]uint, limit, offset int64) ([]model.Article, error) 

*/
func (m *mockDB) Find(out interface{}) *gorm.DB {
	return m.findFunc(out)
}

func (m *mockDB) Limit(limit interface{}) *gorm.DB {
	return m
}

func (m *mockDB) Offset(offset interface{}) *gorm.DB {
	return m
}

func (m *mockDB) Preload(column string) *gorm.DB {
	return m
}

func (m *mockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return m
}


/*
ROOST_METHOD_HASH=ArticleStore_IsFavorited_799826fee5
ROOST_METHOD_SIG_HASH=ArticleStore_IsFavorited_f6d5e67492

FUNCTION_DEF=func (s *ArticleStore) IsFavorited(a *model.Article, u *model.User) (bool, error) // IsFavorited returns whether the article is favorited by the user


*/
func (m *mockDB) Count(value interface{}) *gorm.DB {
	*value.(*int) = m.countResult
	return &gorm.DB{Error: m.countError}
}

func (m *mockDB) Table(name string) *gorm.DB {
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
			name:            "Nil Article parameter",
			article:         nil,
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 0,
			mockCountError:  nil,
			expectedResult:  false,
			expectedError:   nil,
		},
		{
			name:            "Nil User parameter",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            nil,
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
			name:            "Zero count for non-favorited article",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 0,
			mockCountError:  nil,
			expectedResult:  false,
			expectedError:   nil,
		},
		{
			name:            "Multiple favorites for the same article",
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

			store := &ArticleStore{
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

func (m *mockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return &gorm.DB{Error: m.countError}
}

