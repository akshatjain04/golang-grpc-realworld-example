package store

import (
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
			name: "Scenario 1: Retrieve Articles Without Filters",
			mockSetup: func(m *MockDB) {
				m.On("Preload", "Author").Return(m.DB)
				m.On("Preload", "Tags").Return(m.DB)
				m.On("Preload", "Favorites").Return(m.DB)
				m.On("Find", mock.Anything).Return(m.DB)
			},
			expected: []model.Article{
				{Title: "Article 1"},
				{Title: "Article 2"},
			},
		},
		{
			name:     "Scenario 2: Filter Articles by Username",
			username: "testuser",
			mockSetup: func(m *MockDB) {
				m.On("Preload", "Author").Return(m.DB)
				m.On("Preload", "Tags").Return(m.DB)
				m.On("Preload", "Favorites").Return(m.DB)
				m.On("Joins", mock.Anything).Return(m.DB)
				m.On("Where", mock.Anything, mock.Anything).Return(m.DB)
				m.On("Find", mock.Anything).Return(m.DB)
			},
			expected: []model.Article{
				{Title: "User Article"},
			},
		},
		{
			name:    "Scenario 3: Filter Articles by Tag",
			tagName: "testtag",
			mockSetup: func(m *MockDB) {
				m.On("Preload", "Author").Return(m.DB)
				m.On("Preload", "Tags").Return(m.DB)
				m.On("Preload", "Favorites").Return(m.DB)
				m.On("Joins", mock.Anything).Return(m.DB)
				m.On("Where", mock.Anything, mock.Anything).Return(m.DB)
				m.On("Find", mock.Anything).Return(m.DB)
			},
			expected: []model.Article{
				{Title: "Tagged Article"},
			},
		},
		{
			name: "Scenario 4: Retrieve Favorited Articles",
			favoritedBy: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(m *MockDB) {
				m.On("Preload", "Author").Return(m.DB)
				m.On("Preload", "Tags").Return(m.DB)
				m.On("Preload", "Favorites").Return(m.DB)
				m.On("Joins", mock.Anything).Return(m.DB)
				m.On("Where", mock.Anything, mock.Anything).Return(m.DB)
				m.On("Find", mock.Anything).Return(m.DB)
			},
			expected: []model.Article{
				{Title: "Favorited Article"},
			},
		},
		{
			name:   "Scenario 5: Apply Limit and Offset",
			limit:  5,
			offset: 10,
			mockSetup: func(m *MockDB) {
				m.On("Preload", "Author").Return(m.DB)
				m.On("Preload", "Tags").Return(m.DB)
				m.On("Preload", "Favorites").Return(m.DB)
				m.On("Offset", int64(10)).Return(m.DB)
				m.On("Limit", int64(5)).Return(m.DB)
				m.On("Find", mock.Anything).Return(m.DB)
			},
			expected: []model.Article{
				{Title: "Paginated Article 1"},
				{Title: "Paginated Article 2"},
			},
		},
		{
			name:     "Scenario 6: Combine Multiple Filters",
			username: "testuser",
			tagName:  "testtag",
			favoritedBy: &model.User{
				Model: gorm.Model{ID: 1},
			},
			limit:  5,
			offset: 0,
			mockSetup: func(m *MockDB) {
				m.On("Preload", "Author").Return(m.DB)
				m.On("Preload", "Tags").Return(m.DB)
				m.On("Preload", "Favorites").Return(m.DB)
				m.On("Joins", mock.Anything).Return(m.DB)
				m.On("Where", mock.Anything, mock.Anything).Return(m.DB)
				m.On("Offset", int64(0)).Return(m.DB)
				m.On("Limit", int64(5)).Return(m.DB)
				m.On("Find", mock.Anything).Return(m.DB)
			},
			expected: []model.Article{
				{Title: "Filtered Article"},
			},
		},
		{
			name:     "Scenario 7: Handle Non-Existent Data",
			username: "nonexistentuser",
			mockSetup: func(m *MockDB) {
				m.On("Preload", "Author").Return(m.DB)
				m.On("Preload", "Tags").Return(m.DB)
				m.On("Preload", "Favorites").Return(m.DB)
				m.On("Joins", mock.Anything).Return(m.DB)
				m.On("Where", mock.Anything, mock.Anything).Return(m.DB)
				m.On("Find", mock.Anything).Return(m.DB)
			},
			expected: []model.Article{},
		},
		{
			name: "Scenario 8: Error Handling for Database Issues",
			mockSetup: func(m *MockDB) {
				m.On("Preload", "Author").Return(m.DB)
				m.On("Preload", "Tags").Return(m.DB)
				m.On("Preload", "Favorites").Return(m.DB)
				m.On("Find", mock.Anything).Return(gorm.ErrRecordNotFound)
			},
			expected:    []model.Article{},
			expectedErr: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{}
			setupMockDB(mockDB)
			tt.mockSetup(mockDB)

			store := &ArticleStore{db: mockDB.DB}

			articles, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, articles)
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
