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

type mockDB struct {
	countResult int
	countError  error
}
type MockDB struct {
	mock.Mock
	*gorm.DB
}

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

/*
ROOST_METHOD_HASH=GetArticles_101b7250e8
ROOST_METHOD_SIG_HASH=GetArticles_91bc0a6760

FUNCTION_DEF=func (s *ArticleStore) GetArticles(tagName, username string, favoritedBy *model.User, limit, offset int64) ([ // GetArticles get global articles
]model.Article, error)
*/
func TestArticleStoreGetArticles(t *testing.T) {
	tests := []struct {
		name        string
		tagName     string
		username    string
		favoritedBy *model.User
		limit       int64
		offset      int64
		mockSetup   func(*MockDB)
		want        []model.Article
		wantErr     bool
	}{
		{
			name:        "Get Articles with No Filters",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockSetup: func(m *MockDB) {

			},
			want:    []model.Article{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{}
			setupMockDB(mockDB)
			tt.mockSetup(mockDB)

			store := &ArticleStore{db: mockDB.DB}

			got, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func setupMockDB(mockDB *MockDB) {
	db, _, _ := sqlmock.New()
	gdb, _ := gorm.Open("mysql", db)
	mockDB.DB = gdb
}
