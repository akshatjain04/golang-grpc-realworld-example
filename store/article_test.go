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
			name:        "Retrieve Articles Without Filters",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockSetup: func(m *MockDB) {
				m.On("Preload", "Author").Return(m)
				m.On("Offset", int64(0)).Return(m)
				m.On("Limit", int64(10)).Return(m)
				m.On("Find", mock.AnythingOfType("*[]model.Article"), mock.Anything).Return(m).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*[]model.Article)
					*arg = []model.Article{{Title: "Article 1"}, {Title: "Article 2"}}
				})
			},
			expected: []model.Article{{Title: "Article 1"}, {Title: "Article 2"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			setupMockDB(mockDB)
			tt.mockSetup(mockDB)

			store := &ArticleStore{db: mockDB.DB}

			articles, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
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
