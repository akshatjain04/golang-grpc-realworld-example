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
type mockDB struct {
	countResult int
	countError  error
	*gorm.DB
}

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

/*
ROOST_METHOD_HASH=IsFavorited_799826fee5
ROOST_METHOD_SIG_HASH=IsFavorited_f6d5e67492

FUNCTION_DEF=func (s *ArticleStore) IsFavorited(a *model.Article, u *model.User) (bool, error) // IsFavorited returns whether the article is favorited by the user
*/
func (m *mockDB) Count(value interface{}) *gorm.DB {
	*(value.(*int)) = m.countResult
	return &gorm.DB{Error: m.countError}
}

func (m *mockDB) Table(name string) *gorm.DB {
	return m.DB
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
			countError:  nil,
			want:        true,
			wantErr:     false,
		},
		{
			name:        "Article is not favorited by the user",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        &model.User{Model: gorm.Model{ID: 1}},
			countResult: 0,
			countError:  nil,
			want:        false,
			wantErr:     false,
		},
		{
			name:        "Nil Article parameter",
			article:     nil,
			user:        &model.User{Model: gorm.Model{ID: 1}},
			countResult: 0,
			countError:  nil,
			want:        false,
			wantErr:     false,
		},
		{
			name:        "Nil User parameter",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        nil,
			countResult: 0,
			countError:  nil,
			want:        false,
			wantErr:     false,
		},
		{
			name:        "Database error",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        &model.User{Model: gorm.Model{ID: 1}},
			countResult: 0,
			countError:  errors.New("database error"),
			want:        false,
			wantErr:     true,
		},
		{
			name:        "Multiple favorites for the same article and user",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        &model.User{Model: gorm.Model{ID: 1}},
			countResult: 2,
			countError:  nil,
			want:        true,
			wantErr:     false,
		},
		{
			name:        "Zero count but no error from database",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        &model.User{Model: gorm.Model{ID: 1}},
			countResult: 0,
			countError:  nil,
			want:        false,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{
				DB:          &gorm.DB{},
				countResult: tt.countResult,
				countError:  tt.countError,
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
	return m.DB
}
