package github.com/raahii/golang-grpc-realworld-example/store

import (
	errors "errors"
	testing "testing"
	gorm "github.com/jinzhu/gorm"
	model "github.com/raahii/golang-grpc-realworld-example/model"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	sync "sync"
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

func TestArticleStoreAddFavorite(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(*mockDB)
		article       *model.Article
		user          *model.User
		expectedError error
		expectedCount int32
	}{
		{
			name: "Successfully Add Favorite",
			setupMock: func(m *mockDB) {
				tx := &gorm.DB{}
				m.On("Begin").Return(tx)
				m.On("Model", mock.Anything).Return(m)
				assoc := &mockAssociation{}
				assoc.On("Append", mock.Anything).Return(nil)
				m.On("Association", "FavoritedUsers").Return(assoc)
				m.On("Update", "favorites_count", gorm.Expr("favorites_count + ?", 1)).Return(m)
				m.On("Commit").Return(tx)
			},
			article:       &model.Article{FavoritesCount: 0},
			user:          &model.User{},
			expectedError: nil,
			expectedCount: 1,
		},
		{
			name: "Database Error During Association Append",
			setupMock: func(m *mockDB) {
				tx := &gorm.DB{}
				m.On("Begin").Return(tx)
				m.On("Model", mock.Anything).Return(m)
				assoc := &mockAssociation{}
				assoc.On("Append", mock.Anything).Return(errors.New("database error"))
				m.On("Association", "FavoritedUsers").Return(assoc)
				m.On("Rollback").Return(tx)
			},
			article:       &model.Article{FavoritesCount: 0},
			user:          &model.User{},
			expectedError: errors.New("database error"),
			expectedCount: 0,
		},
		{
			name: "Database Error During Favorites Count Update",
			setupMock: func(m *mockDB) {
				tx := &gorm.DB{}
				m.On("Begin").Return(tx)
				m.On("Model", mock.Anything).Return(m)
				assoc := &mockAssociation{}
				assoc.On("Append", mock.Anything).Return(nil)
				m.On("Association", "FavoritedUsers").Return(assoc)
				m.On("Update", "favorites_count", gorm.Expr("favorites_count + ?", 1)).Return(&gorm.DB{Error: errors.New("update error")})
				m.On("Rollback").Return(tx)
			},
			article:       &model.Article{FavoritesCount: 0},
			user:          &model.User{},
			expectedError: errors.New("update error"),
			expectedCount: 0,
		},
		{
			name:          "Add Favorite with Nil Article",
			setupMock:     func(m *mockDB) {},
			article:       nil,
			user:          &model.User{},
			expectedError: errors.New("article cannot be nil"),
			expectedCount: 0,
		},
		{
			name:          "Add Favorite with Nil User",
			setupMock:     func(m *mockDB) {},
			article:       &model.Article{FavoritesCount: 0},
			user:          nil,
			expectedError: errors.New("user cannot be nil"),
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(mockDB)
			tt.setupMock(mockDB)

			store := &ArticleStore{db: mockDB}
			err := store.AddFavorite(tt.article, tt.user)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if tt.article != nil {
				assert.Equal(t, tt.expectedCount, tt.article.FavoritesCount)
			}

			mockDB.AssertExpectations(t)
		})
	}
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
		setupMock      func(*mockDB)
		article        *model.Article
		user           *model.User
		expectedError  error
		expectedCount  int32
		concurrentTest bool
	}{
		{
			name: "Successfully Unfavorite an Article",
			setupMock: func(m *mockDB) {
				tx := &gorm.DB{}
				m.On("Begin").Return(tx)
				m.On("Model", mock.Anything).Return(tx)
				assoc := &mockAssociation{}
				assoc.On("Delete", mock.Anything).Return(&gorm.Association{})
				m.On("Association", "FavoritedUsers").Return(assoc)
				m.On("Update", "favorites_count", gorm.Expr("favorites_count - ?", 1)).Return(tx)
				m.On("Commit").Return(tx)
			},
			article:       &model.Article{FavoritesCount: 1},
			user:          &model.User{},
			expectedError: nil,
			expectedCount: 0,
		},
		{
			name: "Attempt to Unfavorite an Article That Wasn't Favorited",
			setupMock: func(m *mockDB) {
				tx := &gorm.DB{}
				m.On("Begin").Return(tx)
				m.On("Model", mock.Anything).Return(tx)
				assoc := &mockAssociation{}
				assoc.On("Delete", mock.Anything).Return(&gorm.Association{})
				m.On("Association", "FavoritedUsers").Return(assoc)
				m.On("Update", "favorites_count", gorm.Expr("favorites_count - ?", 1)).Return(tx)
				m.On("Commit").Return(tx)
			},
			article:       &model.Article{FavoritesCount: 0},
			user:          &model.User{},
			expectedError: nil,
			expectedCount: 0,
		},
		{
			name: "Database Error During Association Deletion",
			setupMock: func(m *mockDB) {
				tx := &gorm.DB{}
				m.On("Begin").Return(tx)
				m.On("Model", mock.Anything).Return(tx)
				assoc := &mockAssociation{}
				assoc.On("Delete", mock.Anything).Return(&gorm.Association{Error: errors.New("DB error")})
				m.On("Association", "FavoritedUsers").Return(assoc)
				m.On("Rollback").Return(tx)
			},
			article:       &model.Article{FavoritesCount: 1},
			user:          &model.User{},
			expectedError: errors.New("DB error"),
			expectedCount: 1,
		},
		{
			name: "Database Error During FavoritesCount Update",
			setupMock: func(m *mockDB) {
				tx := &gorm.DB{}
				m.On("Begin").Return(tx)
				m.On("Model", mock.Anything).Return(tx)
				assoc := &mockAssociation{}
				assoc.On("Delete", mock.Anything).Return(&gorm.Association{})
				m.On("Association", "FavoritedUsers").Return(assoc)
				m.On("Update", "favorites_count", gorm.Expr("favorites_count - ?", 1)).Return(&gorm.DB{Error: errors.New("Update error")})
				m.On("Rollback").Return(tx)
			},
			article:       &model.Article{FavoritesCount: 1},
			user:          &model.User{},
			expectedError: errors.New("Update error"),
			expectedCount: 1,
		},
		{
			name: "Concurrent Unfavoriting",
			setupMock: func(m *mockDB) {
				tx := &gorm.DB{}
				m.On("Begin").Return(tx)
				m.On("Model", mock.Anything).Return(tx)
				assoc := &mockAssociation{}
				assoc.On("Delete", mock.Anything).Return(&gorm.Association{})
				m.On("Association", "FavoritedUsers").Return(assoc)
				m.On("Update", "favorites_count", gorm.Expr("favorites_count - ?", 1)).Return(tx)
				m.On("Commit").Return(tx)
			},
			article:        &model.Article{FavoritesCount: 5},
			user:           &model.User{},
			expectedError:  nil,
			expectedCount:  0,
			concurrentTest: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(mockDB)
			tt.setupMock(mockDB)

			store := &ArticleStore{
				db: mockDB,
			}

			if tt.concurrentTest {
				var wg sync.WaitGroup
				for i := 0; i < 5; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						err := store.DeleteFavorite(tt.article, tt.user)
						assert.NoError(t, err)
					}()
				}
				wg.Wait()
			} else {
				err := store.DeleteFavorite(tt.article, tt.user)
				assert.Equal(t, tt.expectedError, err)
			}

			assert.Equal(t, tt.expectedCount, tt.article.FavoritesCount)
			mockDB.AssertExpectations(t)
		})
	}
}

func (m *mockAssociation) Delete(values ...interface{}) *gorm.Association {
	args := m.Called(values...)
	return args.Get(0).(*gorm.Association)
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
			id:   ^uint(0),
			mockFindFunc: func(out interface{}, where ...interface{}) *gorm.DB {
				*(out.(*model.Comment)) = model.Comment{
					Model: gorm.Model{ID: ^uint(0)},
					Body:  "Large ID comment",
				}
				return &gorm.DB{Error: nil}
			},
			expectedError: nil,
			expectedComment: &model.Comment{
				Model: gorm.Model{ID: ^uint(0)},
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

