package store

import (
	sql "database/sql"
	testing "testing"

	gorm "github.com/jinzhu/gorm"
	model "github.com/raahii/golang-grpc-realworld-example/model"
	assert "github.com/stretchr/testify/assert"
)

type mockDB struct {
	*gorm.DB
}

/*
ROOST_METHOD_HASH=ArticleStore_GetArticles_101b7250e8
ROOST_METHOD_SIG_HASH=ArticleStore_GetArticles_91bc0a6760

FUNCTION_DEF=func (s *ArticleStore) GetArticles(tagName, username string, favoritedBy *model.User, limit, offset int64) ([ // GetArticles get global articles
]model.Article, error)
*/
func TestArticleStoreGetArticles(t *testing.T) {
	tests := []struct {
		name           string
		tagName        string
		username       string
		favoritedBy    *model.User
		limit          int64
		offset         int64
		mockSetup      func(*mockDB)
		expectedResult []model.Article
		expectedError  error
	}{
		{
			name:        "Scenario 1: Get Articles with No Filters",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockSetup: func(m *mockDB) {
				m.DB = &gorm.DB{Error: nil}
			},
			expectedResult: []model.Article{{Title: "Article 1"}, {Title: "Article 2"}},
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{}
			tt.mockSetup(mockDB)

			store := &ArticleStore{
				db: mockDB.DB,
			}

			result, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

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

func (m *mockDB) Rows() (*sql.Rows, error) {
	return nil, nil
}

func (m *mockDB) Select(query interface{}, args ...interface{}) *gorm.DB {
	return m.DB
}

func (m *mockDB) Table(name string) *gorm.DB {
	return m.DB
}

func (m *mockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return m.DB
}
