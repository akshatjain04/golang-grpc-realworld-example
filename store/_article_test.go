package github.com/raahii/golang-grpc-realworld-example/store

import (
	errors "errors"
	sync "sync"
	testing "testing"
	gorm "github.com/jinzhu/gorm"
	model "github.com/raahii/golang-grpc-realworld-example/model"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
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
		expectedCommit bool
	}{
		{
			name: "Successfully Delete a Favorite Article",
			setupMock: func(db *mockDB, assoc *mockAssociation) {
				db.On("Begin").Return(db)
				db.On("Model", mock.Anything).Return(db)
				db.On("Association", "FavoritedUsers").Return(assoc)
				assoc.On("Delete", mock.Anything).Return(assoc)
				assoc.On("Error").Return(nil)
				db.On("Update", "favorites_count", gorm.Expr("favorites_count - ?", 1)).Return(db)
				db.On("Commit").Return(db)
			},
			article:        &model.Article{FavoritesCount: 2},
			user:           &model.User{},
			expectedError:  nil,
			expectedCount:  1,
			expectedCommit: true,
		},
		{
			name: "Delete Favorite for Non-Existent Association",
			setupMock: func(db *mockDB, assoc *mockAssociation) {
				db.On("Begin").Return(db)
				db.On("Model", mock.Anything).Return(db)
				db.On("Association", "FavoritedUsers").Return(assoc)
				assoc.On("Delete", mock.Anything).Return(assoc)
				assoc.On("Error").Return(nil)
				db.On("Update", "favorites_count", gorm.Expr("favorites_count - ?", 1)).Return(db)
				db.On("Commit").Return(db)
			},
			article:        &model.Article{FavoritesCount: 0},
			user:           &model.User{},
			expectedError:  nil,
			expectedCount:  0,
			expectedCommit: true,
		},
		{
			name: "Database Transaction Rollback on Association Deletion Error",
			setupMock: func(db *mockDB, assoc *mockAssociation) {
				db.On("Begin").Return(db)
				db.On("Model", mock.Anything).Return(db)
				db.On("Association", "FavoritedUsers").Return(assoc)
				assoc.On("Delete", mock.Anything).Return(assoc)
				assoc.On("Error").Return(errors.New("association deletion error"))
				db.On("Rollback").Return(db)
			},
			article:        &model.Article{FavoritesCount: 1},
			user:           &model.User{},
			expectedError:  errors.New("association deletion error"),
			expectedCount:  1,
			expectedCommit: false,
		},
		{
			name: "Database Transaction Rollback on Update Error",
			setupMock: func(db *mockDB, assoc *mockAssociation) {
				db.On("Begin").Return(db)
				db.On("Model", mock.Anything).Return(db)
				db.On("Association", "FavoritedUsers").Return(assoc)
				assoc.On("Delete", mock.Anything).Return(assoc)
				assoc.On("Error").Return(nil)
				db.On("Update", "favorites_count", gorm.Expr("favorites_count - ?", 1)).Return(db)
				db.On("Error").Return(errors.New("update error"))
				db.On("Rollback").Return(db)
			},
			article:        &model.Article{FavoritesCount: 1},
			user:           &model.User{},
			expectedError:  errors.New("update error"),
			expectedCount:  1,
			expectedCommit: false,
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

			if tt.expectedCommit {
				mockDB.AssertCalled(t, "Commit")
			} else {
				mockDB.AssertCalled(t, "Rollback")
			}
		})
	}
}

func TestArticleStoreDeleteFavoriteConcurrent(t *testing.T) {

	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	db.AutoMigrate(&model.Article{}, &model.User{})

	store := &ArticleStore{db: db}

	article := &model.Article{
		Title:          "Test Article",
		FavoritesCount: 5,
	}
	db.Create(article)

	users := make([]*model.User, 5)
	for i := 0; i < 5; i++ {
		users[i] = &model.User{Username: "user" + string(i)}
		db.Create(users[i])
		db.Model(article).Association("FavoritedUsers").Append(users[i])
	}

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(user *model.User) {
			defer wg.Done()
			err := store.DeleteFavorite(article, user)
			assert.NoError(t, err)
		}(users[i])
	}
	wg.Wait()

	var finalArticle model.Article
	db.First(&finalArticle, article.ID)

	assert.Equal(t, int32(0), finalArticle.FavoritesCount)
	assert.Equal(t, 0, db.Model(&finalArticle).Association("FavoritedUsers").Count())
}

func (m *mockAssociation) Delete(values ...interface{}) *gorm.Association {
	args := m.Called(values...)
	return args.Get(0).(*gorm.Association)
}

func (m *mockAssociation) Error() error {
	args := m.Called()
	return args.Error(0)
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
					Model: gorm.Model{ID: 1},
					Body:  "Test comment",
				}
				return &gorm.DB{Error: nil}
			},
			expectedError: nil,
			expectedComment: &model.Comment{
				Model: gorm.Model{ID: 1},
				Body:  "Test comment",
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
			name: "Retrieve a comment with the minimum possible ID",
			id:   1,
			mockFindFunc: func(out interface{}, where ...interface{}) *gorm.DB {
				*(out.(*model.Comment)) = model.Comment{
					Model: gorm.Model{ID: 1},
					Body:  "Minimum ID comment",
				}
				return &gorm.DB{Error: nil}
			},
			expectedError: nil,
			expectedComment: &model.Comment{
				Model: gorm.Model{ID: 1},
				Body:  "Minimum ID comment",
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
			name: "Retrieve a comment with a very large ID",
			id:   4294967295,
			mockFindFunc: func(out interface{}, where ...interface{}) *gorm.DB {
				*(out.(*model.Comment)) = model.Comment{
					Model: gorm.Model{ID: 4294967295},
					Body:  "Large ID comment",
				}
				return &gorm.DB{Error: nil}
			},
			expectedError: nil,
			expectedComment: &model.Comment{
				Model: gorm.Model{ID: 4294967295},
				Body:  "Large ID comment",
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

func (m *mockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	return m.findFunc(out, where...)
}

