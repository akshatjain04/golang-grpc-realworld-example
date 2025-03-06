package store

import (
	errors "errors"
	testing "testing"
	time "time"

	gorm "github.com/jinzhu/gorm"
	model "github.com/raahii/golang-grpc-realworld-example/model"
	assert "github.com/stretchr/testify/assert"
)

type mockDB struct {
	beginCalled    bool
	commitCalled   bool
	rollbackCalled bool
	appendError    error
	updateError    error
	DB             *gorm.DB
	findFunc       func(out interface{}) *gorm.DB
	countError     error
	count          int
}

func TestArticleStoreAddFavorite(t *testing.T) {
	tests := []struct {
		name           string
		article        *model.Article
		user           *model.User
		appendError    error
		updateError    error
		expectedError  error
		expectedCount  int32
		expectRollback bool
	}{
		{
			name:          "Successfully Add Favorite",
			article:       &model.Article{FavoritesCount: 0},
			user:          &model.User{},
			expectedCount: 1,
		},
		{
			name:           "Database Error on Association",
			article:        &model.Article{FavoritesCount: 0},
			user:           &model.User{},
			appendError:    errors.New("association error"),
			expectedError:  errors.New("association error"),
			expectRollback: true,
		},
		{
			name:           "Database Error on Update",
			article:        &model.Article{FavoritesCount: 0},
			user:           &model.User{},
			updateError:    errors.New("update error"),
			expectedError:  errors.New("update error"),
			expectRollback: true,
		},
		{
			name:          "Add Favorite with Nil Article",
			article:       nil,
			user:          &model.User{},
			expectedError: errors.New("article or user is nil"),
		},
		{
			name:          "Add Favorite with Nil User",
			article:       &model.Article{FavoritesCount: 0},
			user:          nil,
			expectedError: errors.New("article or user is nil"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{
				appendError: tt.appendError,
				updateError: tt.updateError,
				DB:          &gorm.DB{},
			}

			store := &ArticleStore{
				db: mockDB.DB,
			}

			err := store.AddFavorite(tt.article, tt.user)

			if tt.expectedError != nil {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.article != nil && tt.expectedError == nil {
				if tt.article.FavoritesCount != tt.expectedCount {
					t.Errorf("expected favorites count %d, got %d", tt.expectedCount, tt.article.FavoritesCount)
				}
			}

			if tt.expectRollback && !mockDB.rollbackCalled {
				t.Error("expected rollback to be called, but it wasn't")
			}

			if !tt.expectRollback && !mockDB.commitCalled {
				t.Error("expected commit to be called, but it wasn't")
			}
		})
	}
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
			name:        "Get Articles Without Any Filters",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockSetup: func(m *mockDB) {
				m.DB.Error = nil
			},
			expected:    []model.Article{},
			expectedErr: nil,
		},
		{
			name:        "Get Articles by Username",
			tagName:     "",
			username:    "testuser",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockSetup: func(m *mockDB) {
				m.DB.Error = nil
			},
			expected:    []model.Article{},
			expectedErr: nil,
		},
		{
			name:        "Get Articles by Tag",
			tagName:     "testtag",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
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
			name:        "Combine Multiple Filters",
			tagName:     "testtag",
			username:    "testuser",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockSetup: func(m *mockDB) {
				m.DB.Error = nil
			},
			expected:    []model.Article{},
			expectedErr: nil,
		},
		{
			name:        "Handle Empty Result Set",
			tagName:     "nonexistenttag",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockSetup: func(m *mockDB) {
				m.DB.Error = nil
			},
			expected:    []model.Article{},
			expectedErr: nil,
		},
		{
			name:        "Test Pagination with Limit and Offset",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       5,
			offset:      5,
			mockSetup: func(m *mockDB) {
				m.DB.Error = nil
			},
			expected:    []model.Article{},
			expectedErr: nil,
		},
		{
			name:        "Handle Database Error",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockSetup: func(m *mockDB) {
				m.DB.Error = errors.New("database error")
			},
			expected:    []model.Article{},
			expectedErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{DB: &gorm.DB{}}
			tt.mockSetup(mockDB)

			store := &ArticleStore{db: mockDB.DB}
			articles, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			assert.Equal(t, tt.expected, articles)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestArticleStoreGetFeedArticles(t *testing.T) {
	tests := []struct {
		name     string
		userIDs  []uint
		limit    int64
		offset   int64
		mockFind func(out interface{}) *gorm.DB
		want     []model.Article
		wantErr  bool
	}{
		{
			name:    "Successful Retrieval of Feed Articles",
			userIDs: []uint{1, 2},
			limit:   2,
			offset:  0,
			mockFind: func(out interface{}) *gorm.DB {
				articles := out.(*[]model.Article)
				*articles = []model.Article{
					{
						Model:       gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
						Title:       "Article 1",
						Description: "Description 1",
						Body:        "Body 1",
						UserID:      1,
						Author:      model.User{Model: gorm.Model{ID: 1}, Username: "user1"},
					},
					{
						Model:       gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
						Title:       "Article 2",
						Description: "Description 2",
						Body:        "Body 2",
						UserID:      2,
						Author:      model.User{Model: gorm.Model{ID: 2}, Username: "user2"},
					},
				}
				return &gorm.DB{Error: nil}
			},
			want: []model.Article{
				{
					Model:       gorm.Model{ID: 1},
					Title:       "Article 1",
					Description: "Description 1",
					Body:        "Body 1",
					UserID:      1,
					Author:      model.User{Model: gorm.Model{ID: 1}, Username: "user1"},
				},
				{
					Model:       gorm.Model{ID: 2},
					Title:       "Article 2",
					Description: "Description 2",
					Body:        "Body 2",
					UserID:      2,
					Author:      model.User{Model: gorm.Model{ID: 2}, Username: "user2"},
				},
			},
			wantErr: false,
		},
		{
			name:    "Empty Result Set",
			userIDs: []uint{3, 4},
			limit:   10,
			offset:  0,
			mockFind: func(out interface{}) *gorm.DB {
				return &gorm.DB{Error: nil}
			},
			want:    []model.Article{},
			wantErr: false,
		},
		{
			name:    "Database Error Handling",
			userIDs: []uint{1, 2},
			limit:   10,
			offset:  0,
			mockFind: func(out interface{}) *gorm.DB {
				return &gorm.DB{Error: errors.New("database error")}
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Limit Exceeds Available Articles",
			userIDs: []uint{1},
			limit:   5,
			offset:  0,
			mockFind: func(out interface{}) *gorm.DB {
				articles := out.(*[]model.Article)
				*articles = []model.Article{
					{
						Model:       gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
						Title:       "Article 1",
						Description: "Description 1",
						Body:        "Body 1",
						UserID:      1,
						Author:      model.User{Model: gorm.Model{ID: 1}, Username: "user1"},
					},
				}
				return &gorm.DB{Error: nil}
			},
			want: []model.Article{
				{
					Model:       gorm.Model{ID: 1},
					Title:       "Article 1",
					Description: "Description 1",
					Body:        "Body 1",
					UserID:      1,
					Author:      model.User{Model: gorm.Model{ID: 1}, Username: "user1"},
				},
			},
			wantErr: false,
		},
		{
			name:    "Correct Application of Offset",
			userIDs: []uint{1, 2},
			limit:   2,
			offset:  1,
			mockFind: func(out interface{}) *gorm.DB {
				articles := out.(*[]model.Article)
				*articles = []model.Article{
					{
						Model:       gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
						Title:       "Article 2",
						Description: "Description 2",
						Body:        "Body 2",
						UserID:      2,
						Author:      model.User{Model: gorm.Model{ID: 2}, Username: "user2"},
					},
					{
						Model:       gorm.Model{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
						Title:       "Article 3",
						Description: "Description 3",
						Body:        "Body 3",
						UserID:      1,
						Author:      model.User{Model: gorm.Model{ID: 1}, Username: "user1"},
					},
				}
				return &gorm.DB{Error: nil}
			},
			want: []model.Article{
				{
					Model:       gorm.Model{ID: 2},
					Title:       "Article 2",
					Description: "Description 2",
					Body:        "Body 2",
					UserID:      2,
					Author:      model.User{Model: gorm.Model{ID: 2}, Username: "user2"},
				},
				{
					Model:       gorm.Model{ID: 3},
					Title:       "Article 3",
					Description: "Description 3",
					Body:        "Body 3",
					UserID:      1,
					Author:      model.User{Model: gorm.Model{ID: 1}, Username: "user1"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{
				findFunc: tt.mockFind,
			}
			s := &ArticleStore{
				db: mockDB.DB,
			}

			got, err := s.GetFeedArticles(tt.userIDs, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)

			if len(got) > 0 {
				for _, article := range got {
					assert.NotEmpty(t, article.Author)
				}
			}
		})
	}
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
			name:        "Nil article parameter",
			article:     nil,
			user:        &model.User{Model: gorm.Model{ID: 1}},
			dbCount:     0,
			dbError:     nil,
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "Nil user parameter",
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
			name:        "Multiple favorites for the same article-user pair",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        &model.User{Model: gorm.Model{ID: 1}},
			dbCount:     5,
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
			mockDB := &mockDB{
				countError: tt.dbError,
				count:      tt.dbCount,
				DB:         &gorm.DB{},
			}

			store := &ArticleStore{
				db: mockDB.DB,
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

func (m *mockDB) Association(column string) *gorm.Association {
	return &gorm.Association{}
}

func (m *mockDB) Begin() *gorm.DB {
	m.beginCalled = true
	return &gorm.DB{}
}

func (m *mockDB) Commit() *gorm.DB {
	m.commitCalled = true
	return &gorm.DB{}
}

func (m *mockDB) Count(value interface{}) *gorm.DB {
	*value.(*int) = m.count
	return &gorm.DB{Error: m.countError}
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

func (m *mockDB) Model(value interface{}) *gorm.DB {
	return &gorm.DB{}
}

func (m *mockDB) Offset(offset interface{}) *gorm.DB {
	return m.DB
}

func (m *mockDB) Preload(column string, conditions ...interface{}) *gorm.DB {
	return m.DB
}

func (m *mockDB) Rollback() *gorm.DB {
	m.rollbackCalled = true
	return &gorm.DB{}
}

func (m *mockDB) Rows() (*gorm.RowsQueryResult, error) {
	return nil, nil
}

func (m *mockDB) Select(query interface{}, args ...interface{}) *gorm.DB {
	return m.DB
}

func (m *mockDB) Table(name string) *gorm.DB {
	return m.DB
}

func (m *mockDB) Update(column string, value interface{}) *gorm.DB {
	return &gorm.DB{Error: m.updateError}
}

func (m *mockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return m.DB
}
