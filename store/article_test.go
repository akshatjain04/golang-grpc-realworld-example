package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
	*gorm.DB
}

/*
ROOST_METHOD_HASH=GetArticles_101b7250e8
ROOST_METHOD_SIG_HASH=GetArticles_91bc0a6760

FUNCTION_DEF=func (s *ArticleStore) GetArticles(tagName, username string, favoritedBy *model.User, limit, offset int64) ([ // GetArticles get global articles
]model.Article, error)
*/
func (m *MockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	return m.DB
}

func (m *MockDB) Joins(query string, args ...interface{}) *gorm.DB {
	return m.DB
}

func (m *MockDB) Limit(limit interface{}) *gorm.DB {
	return m.DB
}

func (m *MockDB) Offset(offset interface{}) *gorm.DB {
	return m.DB
}

func (m *MockDB) Preload(column string, conditions ...interface{}) *gorm.DB {
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
		mockSetup   func(*MockDB)
		expected    []model.Article
		expectedErr error
	}{
		{
			name:        "Scenario 1: Get Articles with No Filters",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockSetup: func(m *MockDB) {
				m.DB.Error = nil
			},
			expected:    []model.Article{},
			expectedErr: nil,
		},
		{
			name:        "Scenario 2: Get Articles by Tag Name",
			tagName:     "test-tag",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockSetup: func(m *MockDB) {
				m.DB.Error = nil
			},
			expected:    []model.Article{},
			expectedErr: nil,
		},
		{
			name:        "Scenario 3: Get Articles by Author Username",
			tagName:     "",
			username:    "test-user",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockSetup: func(m *MockDB) {
				m.DB.Error = nil
			},
			expected:    []model.Article{},
			expectedErr: nil,
		},
		{
			name:        "Scenario 4: Get Favorited Articles",
			tagName:     "",
			username:    "",
			favoritedBy: &model.User{Model: gorm.Model{ID: 1}},
			limit:       10,
			offset:      0,
			mockSetup: func(m *MockDB) {
				m.DB.Error = nil
			},
			expected:    []model.Article{},
			expectedErr: nil,
		},
		{
			name:        "Scenario 5: Combine Multiple Filters",
			tagName:     "test-tag",
			username:    "test-user",
			favoritedBy: &model.User{Model: gorm.Model{ID: 1}},
			limit:       10,
			offset:      0,
			mockSetup: func(m *MockDB) {
				m.DB.Error = nil
			},
			expected:    []model.Article{},
			expectedErr: nil,
		},
		{
			name:        "Scenario 6: Handle Empty Result Set",
			tagName:     "non-existent-tag",
			username:    "non-existent-user",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockSetup: func(m *MockDB) {
				m.DB.Error = nil
			},
			expected:    []model.Article{},
			expectedErr: nil,
		},
		{
			name:        "Scenario 7: Test Pagination Limits",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       1000,
			offset:      10000,
			mockSetup: func(m *MockDB) {
				m.DB.Error = nil
			},
			expected:    []model.Article{},
			expectedErr: nil,
		},
		{
			name:        "Scenario 8: Error Handling for Database Issues",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockSetup: func(m *MockDB) {
				m.DB.Error = errors.New("database error")
			},
			expected:    []model.Article{},
			expectedErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{}
			setupMockDB(mockDB)
			tt.mockSetup(mockDB)

			store := &ArticleStore{db: mockDB.DB}

			articles, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			assert.Equal(t, tt.expected, articles)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return m.DB
}

func setupMockDB(mockDB *MockDB) {
	db, _, _ := sqlmock.New()
	gdb, _ := gorm.Open("mysql", db)
	mockDB.DB = gdb
}
