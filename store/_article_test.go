package github.com/raahii/golang-grpc-realworld-example/store

import (
	errors "errors"
	sync "sync"
	testing "testing"
	gorm "github.com/jinzhu/gorm"
	model "github.com/raahii/golang-grpc-realworld-example/model"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	math "math"
	reflect "reflect"
)





type mockAssociation struct {
	mock.Mock
}
type mockDB struct {
	mock.Mock
}


/*
ROOST_METHOD_HASH=ArticleStore_AddFavorite_9460fca478
ROOST_METHOD_SIG_HASH=ArticleStore_AddFavorite_c13a109f91

FUNCTION_DEF=func (s *ArticleStore) AddFavorite(a *model.Article, u *model.User) error // AddFavorite favorite an article


*/
func (m *MockAssociation) Append(values ...interface{}) error {
	args := m.Called(values...)
	return args.Error(0)
}

func (m *MockDB) Association(column string) *gorm.Association {
	args := m.Called(column)
	return args.Get(0).(*gorm.Association)
}

func (m *MockDB) Begin() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Commit() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Model(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Rollback() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Update(column string, value interface{}) *gorm.DB {
	args := m.Called(column, value)
	return args.Get(0).(*gorm.DB)
}


/*
ROOST_METHOD_HASH=ArticleStore_DeleteFavorite_29c18a04a8
ROOST_METHOD_SIG_HASH=ArticleStore_DeleteFavorite_53deb5e792

FUNCTION_DEF=func (s *ArticleStore) DeleteFavorite(a *model.Article, u *model.User) error // DeleteFavorite unfavorite an article


*/
func TestArticleStoreDeleteFavorite(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*mockDB, *mockAssociation)
		article        *model.Article
		user           *model.User
		expectedError  error
		expectedCount  int32
		checkFavorited bool
	}{
		{
			name: "Successfully Delete a Favorite Article",
			setupMock: func(db *mockDB, assoc *mockAssociation) {
				tx := &gorm.DB{}
				db.On("Begin").Return(tx)
				db.On("Model", mock.Anything).Return(db)
				db.On("Association", "FavoritedUsers").Return(assoc)
				assoc.On("Delete", mock.Anything).Return(assoc)
				assoc.On("Error").Return(nil)
				db.On("Update", "favorites_count", gorm.Expr("favorites_count - ?", 1)).Return(db)
				db.On("Commit").Return(db)
			},
			article:        &model.Article{FavoritesCount: 5},
			user:           &model.User{},
			expectedError:  nil,
			expectedCount:  4,
			checkFavorited: true,
		},
		{
			name: "Delete Favorite for Non-Existent Association",
			setupMock: func(db *mockDB, assoc *mockAssociation) {
				tx := &gorm.DB{}
				db.On("Begin").Return(tx)
				db.On("Model", mock.Anything).Return(db)
				db.On("Association", "FavoritedUsers").Return(assoc)
				assoc.On("Delete", mock.Anything).Return(assoc)
				assoc.On("Error").Return(nil)
				db.On("Update", "favorites_count", gorm.Expr("favorites_count - ?", 1)).Return(db)
				db.On("Commit").Return(db)
			},
			article:       &model.Article{FavoritesCount: 0},
			user:          &model.User{},
			expectedError: nil,
			expectedCount: 0,
		},
		{
			name: "Database Transaction Rollback on Association Deletion Error",
			setupMock: func(db *mockDB, assoc *mockAssociation) {
				tx := &gorm.DB{}
				db.On("Begin").Return(tx)
				db.On("Model", mock.Anything).Return(db)
				db.On("Association", "FavoritedUsers").Return(assoc)
				assoc.On("Delete", mock.Anything).Return(assoc)
				assoc.On("Error").Return(errors.New("association deletion error"))
				db.On("Rollback").Return(db)
			},
			article:       &model.Article{FavoritesCount: 5},
			user:          &model.User{},
			expectedError: errors.New("association deletion error"),
			expectedCount: 5,
		},
		{
			name: "Database Transaction Rollback on Update Error",
			setupMock: func(db *mockDB, assoc *mockAssociation) {
				tx := &gorm.DB{}
				db.On("Begin").Return(tx)
				db.On("Model", mock.Anything).Return(db)
				db.On("Association", "FavoritedUsers").Return(assoc)
				assoc.On("Delete", mock.Anything).Return(assoc)
				assoc.On("Error").Return(nil)
				db.On("Update", "favorites_count", gorm.Expr("favorites_count - ?", 1)).Return(&gorm.DB{Error: errors.New("update error")})
				db.On("Rollback").Return(db)
			},
			article:       &model.Article{FavoritesCount: 5},
			user:          &model.User{},
			expectedError: errors.New("update error"),
			expectedCount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(mockDB)
			mockAssoc := new(mockAssociation)
			tt.setupMock(mockDB, mockAssoc)

			store := &ArticleStore{db: mockDB}
			err := store.DeleteFavorite(tt.article, tt.user)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedCount, tt.article.FavoritesCount)

			mockDB.AssertExpectations(t)
			mockAssoc.AssertExpectations(t)
		})
	}
}

func TestArticleStoreDeleteFavoriteConcurrent(t *testing.T) {

	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	db.AutoMigrate(&model.User{}, &model.Article{})

	store := &ArticleStore{db: db}

	article := &model.Article{FavoritesCount: 5}
	db.Create(article)

	users := make([]model.User, 5)
	for i := range users {
		db.Create(&users[i])
		db.Model(article).Association("FavoritedUsers").Append(&users[i])
	}

	var wg sync.WaitGroup
	for i := range users {
		wg.Add(1)
		go func(user *model.User) {
			defer wg.Done()
			err := store.DeleteFavorite(article, user)
			assert.NoError(t, err)
		}(&users[i])
	}
	wg.Wait()

	db.First(article, article.ID)
	assert.Equal(t, int32(0), article.FavoritesCount)

	var favoritedUsers []model.User
	db.Model(article).Association("FavoritedUsers").Find(&favoritedUsers)
	assert.Empty(t, favoritedUsers)
}

func (m *mockAssociation) Delete(values ...interface{}) *gorm.Association {
	args := m.Called(values...)
	return args.Get(0).(*gorm.Association)
}

func (m *mockAssociation) Error() error {
	args := m.Called()
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
func TestArticleStoreGetArticles(t *testing.T) {
	tests := []struct {
		name        string
		tagName     string
		username    string
		favoritedBy *model.User
		limit       int64
		offset      int64
		mockDB      *mockDB
		expected    []model.Article
		expectedErr error
	}{
		{
			name:     "Get Articles Without Filters",
			tagName:  "",
			username: "",
			limit:    10,
			offset:   0,
			mockDB: &mockDB{
				findFunc: func(out interface{}) *gorm.DB {
					*(out.(*[]model.Article)) = []model.Article{
						{Title: "Article 1"},
						{Title: "Article 2"},
					}
					return &gorm.DB{Error: nil}
				},
			},
			expected: []model.Article{
				{Title: "Article 1"},
				{Title: "Article 2"},
			},
			expectedErr: nil,
		},
		{
			name:     "Get Articles by Username",
			username: "testuser",
			limit:    10,
			offset:   0,
			mockDB: &mockDB{
				findFunc: func(out interface{}) *gorm.DB {
					*(out.(*[]model.Article)) = []model.Article{
						{Title: "User Article"},
					}
					return &gorm.DB{Error: nil}
				},
			},
			expected: []model.Article{
				{Title: "User Article"},
			},
			expectedErr: nil,
		},
		{
			name:    "Get Articles by Tag",
			tagName: "testtag",
			limit:   10,
			offset:  0,
			mockDB: &mockDB{
				findFunc: func(out interface{}) *gorm.DB {
					*(out.(*[]model.Article)) = []model.Article{
						{Title: "Tagged Article"},
					}
					return &gorm.DB{Error: nil}
				},
			},
			expected: []model.Article{
				{Title: "Tagged Article"},
			},
			expectedErr: nil,
		},
		{
			name:        "Get Favorited Articles",
			favoritedBy: &model.User{Model: gorm.Model{ID: 1}},
			limit:       10,
			offset:      0,
			mockDB: &mockDB{
				rowsFunc: func() (*gorm.Rows, error) {
					return nil, nil
				},
				findFunc: func(out interface{}) *gorm.DB {
					*(out.(*[]model.Article)) = []model.Article{
						{Title: "Favorited Article"},
					}
					return &gorm.DB{Error: nil}
				},
			},
			expected: []model.Article{
				{Title: "Favorited Article"},
			},
			expectedErr: nil,
		},
		{
			name:     "Combine Multiple Filters",
			username: "testuser",
			tagName:  "testtag",
			limit:    10,
			offset:   0,
			mockDB: &mockDB{
				findFunc: func(out interface{}) *gorm.DB {
					*(out.(*[]model.Article)) = []model.Article{
						{Title: "Combined Filter Article"},
					}
					return &gorm.DB{Error: nil}
				},
			},
			expected: []model.Article{
				{Title: "Combined Filter Article"},
			},
			expectedErr: nil,
		},
		{
			name:   "Handle Empty Result Set",
			limit:  10,
			offset: 0,
			mockDB: &mockDB{
				findFunc: func(out interface{}) *gorm.DB {
					*(out.(*[]model.Article)) = []model.Article{}
					return &gorm.DB{Error: nil}
				},
			},
			expected:    []model.Article{},
			expectedErr: nil,
		},
		{
			name:   "Test Pagination Limits",
			limit:  1000,
			offset: 10000,
			mockDB: &mockDB{
				findFunc: func(out interface{}) *gorm.DB {
					*(out.(*[]model.Article)) = []model.Article{}
					return &gorm.DB{Error: nil}
				},
			},
			expected:    []model.Article{},
			expectedErr: nil,
		},
		{
			name:   "Error Handling for Database Issues",
			limit:  10,
			offset: 0,
			mockDB: &mockDB{
				findFunc: func(out interface{}) *gorm.DB {
					return &gorm.DB{Error: errors.New("database error")}
				},
			},
			expected:    []model.Article{},
			expectedErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &ArticleStore{db: tt.mockDB}
			articles, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			assert.Equal(t, tt.expected, articles)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func (m *mockDB) Find(out interface{}) *gorm.DB {
	return m.findFunc(out)
}

func (m *mockDB) Joins(query string, args ...interface{}) *gorm.DB {
	return m
}

func (m *mockDB) Limit(limit interface{}) *gorm.DB {
	return m
}

func (m *mockDB) Offset(offset interface{}) *gorm.DB {
	return m
}

func (m *mockDB) Preload(column string, conditions ...interface{}) *gorm.DB {
	return m
}

func (m *mockDB) Rows() (*gorm.Rows, error) {
	return m.rowsFunc()
}

func (m *mockDB) Select(query interface{}, args ...interface{}) *gorm.DB {
	return m
}

func (m *mockDB) Table(name string) *gorm.DB {
	return m
}

func (m *mockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return m
}


/*
ROOST_METHOD_HASH=ArticleStore_GetCommentByID_7ecaa81f20
ROOST_METHOD_SIG_HASH=ArticleStore_GetCommentByID_f6f8a51973

FUNCTION_DEF=func (s *ArticleStore) GetCommentByID(id uint) (*model.Comment, error) // GetCommentByID finds an comment from id


*/
func TestArticleStoreGetCommentById(t *testing.T) {
	tests := []struct {
		name            string
		id              uint
		mockFindFunc    func(out interface{}, where ...interface{}) *gorm.DB
		expectedError   error
		expectedComment *model.Comment
	}{
		{
			name: "Successfully retrieve an existing comment",
			id:   1,
			mockFindFunc: func(out interface{}, where ...interface{}) *gorm.DB {
				*(out.(*model.Comment)) = model.Comment{
					Model:     gorm.Model{ID: 1},
					Body:      "Test comment",
					UserID:    1,
					ArticleID: 1,
				}
				return &gorm.DB{Error: nil}
			},
			expectedError: nil,
			expectedComment: &model.Comment{
				Model:     gorm.Model{ID: 1},
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
		},
		{
			name: "Attempt to retrieve a non-existent comment",
			id:   999,
			mockFindFunc: func(out interface{}, where ...interface{}) *gorm.DB {
				return &gorm.DB{Error: gorm.ErrRecordNotFound}
			},
			expectedError:   gorm.ErrRecordNotFound,
			expectedComment: nil,
		},
		{
			name: "Handle database connection error",
			id:   1,
			mockFindFunc: func(out interface{}, where ...interface{}) *gorm.DB {
				return &gorm.DB{Error: errors.New("database connection error")}
			},
			expectedError:   errors.New("database connection error"),
			expectedComment: nil,
		},
		{
			name: "Retrieve a comment with maximum uint ID",
			id:   math.MaxUint32,
			mockFindFunc: func(out interface{}, where ...interface{}) *gorm.DB {
				*(out.(*model.Comment)) = model.Comment{
					Model:     gorm.Model{ID: math.MaxUint32},
					Body:      "Max ID comment",
					UserID:    1,
					ArticleID: 1,
				}
				return &gorm.DB{Error: nil}
			},
			expectedError: nil,
			expectedComment: &model.Comment{
				Model:     gorm.Model{ID: math.MaxUint32},
				Body:      "Max ID comment",
				UserID:    1,
				ArticleID: 1,
			},
		},
		{
			name: "Attempt to retrieve a comment with ID 0",
			id:   0,
			mockFindFunc: func(out interface{}, where ...interface{}) *gorm.DB {
				return &gorm.DB{Error: gorm.ErrRecordNotFound}
			},
			expectedError:   gorm.ErrRecordNotFound,
			expectedComment: nil,
		},
		{
			name: "Verify all fields of the retrieved comment",
			id:   1,
			mockFindFunc: func(out interface{}, where ...interface{}) *gorm.DB {
				*(out.(*model.Comment)) = model.Comment{
					Model:     gorm.Model{ID: 1, CreatedAt: testTime, UpdatedAt: testTime},
					Body:      "Full comment",
					UserID:    2,
					Author:    model.User{Model: gorm.Model{ID: 2}, Username: "testuser"},
					ArticleID: 3,
					Article:   model.Article{Model: gorm.Model{ID: 3}, Title: "Test Article"},
				}
				return &gorm.DB{Error: nil}
			},
			expectedError: nil,
			expectedComment: &model.Comment{
				Model:     gorm.Model{ID: 1, CreatedAt: testTime, UpdatedAt: testTime},
				Body:      "Full comment",
				UserID:    2,
				Author:    model.User{Model: gorm.Model{ID: 2}, Username: "testuser"},
				ArticleID: 3,
				Article:   model.Article{Model: gorm.Model{ID: 3}, Title: "Test Article"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{findFunc: tt.mockFindFunc}
			store := &ArticleStore{db: mockDB}

			comment, err := store.GetCommentByID(tt.id)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedComment, comment)
		})
	}
}


/*
ROOST_METHOD_HASH=NewArticleStore_85784abca5
ROOST_METHOD_SIG_HASH=NewArticleStore_436ae9c986

FUNCTION_DEF=func NewArticleStore(db *gorm.DB) *ArticleStore // NewArticleStore returns a new ArticleStore


*/
func TestNewArticleStore(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want *ArticleStore
	}{
		{
			name: "Create ArticleStore with valid gorm.DB",
			args: args{db: &gorm.DB{}},
			want: &ArticleStore{db: &gorm.DB{}},
		},
		{
			name: "Create ArticleStore with nil gorm.DB",
			args: args{db: nil},
			want: &ArticleStore{db: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewArticleStore(tt.args.db)
			if got == nil {
				t.Errorf("NewArticleStore() returned nil")
			}
			if !reflect.DeepEqual(got.db, tt.want.db) {
				t.Errorf("NewArticleStore() = %v, want %v", got.db, tt.want.db)
			}
		})
	}

	t.Run("Verify ArticleStore uniqueness for multiple calls", func(t *testing.T) {
		db := &gorm.DB{}
		store1 := NewArticleStore(db)
		store2 := NewArticleStore(db)

		if store1 == store2 {
			t.Errorf("NewArticleStore() returned the same instance for multiple calls")
		}
		if store1.db != store2.db {
			t.Errorf("NewArticleStore() returned instances with different db references")
		}
	})

	t.Run("Check ArticleStore with different gorm.DB instances", func(t *testing.T) {
		db1 := &gorm.DB{}
		db2 := &gorm.DB{}
		store1 := NewArticleStore(db1)
		store2 := NewArticleStore(db2)

		if store1 == store2 {
			t.Errorf("NewArticleStore() returned the same instance for different db instances")
		}
		if store1.db == store2.db {
			t.Errorf("NewArticleStore() returned instances with the same db reference for different db instances")
		}
	})

	t.Run("Verify ArticleStore creation in concurrent environment", func(t *testing.T) {
		db := &gorm.DB{}
		var wg sync.WaitGroup
		stores := make([]*ArticleStore, 100)

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				stores[index] = NewArticleStore(db)
			}(i)
		}

		wg.Wait()

		for i := 1; i < len(stores); i++ {
			if stores[i] == stores[i-1] {
				t.Errorf("NewArticleStore() returned the same instance in concurrent calls")
			}
			if stores[i].db != db {
				t.Errorf("NewArticleStore() returned instance with incorrect db reference in concurrent calls")
			}
		}
	})
}

