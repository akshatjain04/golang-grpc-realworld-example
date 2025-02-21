package github.com/raahii/golang-grpc-realworld-example/store

import (
	"errors"
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"time"
)








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
		name           string
		setupMock      func(*mockDB)
		article        *model.Article
		user           *model.User
		expectedError  error
		expectedCount  int32
		expectedAppend bool
	}{
		{
			name: "Successfully Add Favorite",
			setupMock: func(m *mockDB) {
				tx := &gorm.DB{}
				m.On("Begin").Return(tx)
				m.On("Model", mock.Anything).Return(tx)
				assoc := &mockAssociation{}
				assoc.On("Append", mock.Anything).Return(nil)
				m.On("Association", "FavoritedUsers").Return(assoc)
				m.On("Update", "favorites_count", gorm.Expr("favorites_count + ?", 1)).Return(tx)
				m.On("Commit").Return(tx)
			},
			article:        &model.Article{FavoritesCount: 0},
			user:           &model.User{},
			expectedError:  nil,
			expectedCount:  1,
			expectedAppend: true,
		},
		{
			name: "Add Favorite with Database Error on Association",
			setupMock: func(m *mockDB) {
				tx := &gorm.DB{}
				m.On("Begin").Return(tx)
				m.On("Model", mock.Anything).Return(tx)
				assoc := &mockAssociation{}
				assoc.On("Append", mock.Anything).Return(errors.New("association error"))
				m.On("Association", "FavoritedUsers").Return(assoc)
				m.On("Rollback").Return(tx)
			},
			article:        &model.Article{FavoritesCount: 0},
			user:           &model.User{},
			expectedError:  errors.New("association error"),
			expectedCount:  0,
			expectedAppend: false,
		},
		{
			name: "Add Favorite with Database Error on Update",
			setupMock: func(m *mockDB) {
				tx := &gorm.DB{}
				m.On("Begin").Return(tx)
				m.On("Model", mock.Anything).Return(tx)
				assoc := &mockAssociation{}
				assoc.On("Append", mock.Anything).Return(nil)
				m.On("Association", "FavoritedUsers").Return(assoc)
				m.On("Update", "favorites_count", gorm.Expr("favorites_count + ?", 1)).Return(&gorm.DB{Error: errors.New("update error")})
				m.On("Rollback").Return(tx)
			},
			article:        &model.Article{FavoritesCount: 0},
			user:           &model.User{},
			expectedError:  errors.New("update error"),
			expectedCount:  0,
			expectedAppend: false,
		},
		{
			name: "Add Favorite with Nil Article",
			setupMock: func(m *mockDB) {

			},
			article:        nil,
			user:           &model.User{},
			expectedError:  errors.New("article is nil"),
			expectedCount:  0,
			expectedAppend: false,
		},
		{
			name: "Add Favorite with Nil User",
			setupMock: func(m *mockDB) {

			},
			article:        &model.Article{FavoritesCount: 0},
			user:           nil,
			expectedError:  errors.New("user is nil"),
			expectedCount:  0,
			expectedAppend: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(mockDB)
			if tt.setupMock != nil {
				tt.setupMock(mockDB)
			}

			store := &ArticleStore{
				db: mockDB,
			}

			err := store.AddFavorite(tt.article, tt.user)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
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
ROOST_METHOD_HASH=ArticleStore_GetArticles_101b7250e8
ROOST_METHOD_SIG_HASH=ArticleStore_GetArticles_91bc0a6760

FUNCTION_DEF=func (s *ArticleStore) GetArticles(tagName, username string, favoritedBy *model.User, limit, offset int64) ([ // GetArticles get global articles
]model.Article, error) 

*/
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

func (m *mockDB) Rows() (*gorm.Rows, error) {
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


/*
ROOST_METHOD_HASH=ArticleStore_GetFeedArticles_a37e1934b6
ROOST_METHOD_SIG_HASH=ArticleStore_GetFeedArticles_f5f09c020e

FUNCTION_DEF=func (s *ArticleStore) GetFeedArticles(userIDs [ // GetFeedArticles returns following users' articles
]uint, limit, offset int64) ([]model.Article, error) 

*/
func (m *mockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where)
	return args.Get(0).(*gorm.DB)
}

func (m *mockDB) Limit(limit interface{}) *gorm.DB {
	args := m.Called(limit)
	return args.Get(0).(*gorm.DB)
}

func (m *mockDB) Offset(offset interface{}) *gorm.DB {
	args := m.Called(offset)
	return args.Get(0).(*gorm.DB)
}

func (m *mockDB) Preload(column string, conditions ...interface{}) *gorm.DB {
	args := m.Called(column, conditions)
	return args.Get(0).(*gorm.DB)
}

func (m *mockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	callArgs := m.Called(query, args)
	return callArgs.Get(0).(*gorm.DB)
}


/*
ROOST_METHOD_HASH=ArticleStore_IsFavorited_799826fee5
ROOST_METHOD_SIG_HASH=ArticleStore_IsFavorited_f6d5e67492

FUNCTION_DEF=func (s *ArticleStore) IsFavorited(a *model.Article, u *model.User) (bool, error) // IsFavorited returns whether the article is favorited by the user


*/
func (m *mockDB) Count(value interface{}) *gorm.DB {
	*value.(*int) = m.countResult
	return &gorm.DB{Error: m.countError}
}

func (m *mockDB) Table(name string) *gorm.DB {
	return &gorm.DB{Error: m.countError}
}

func TestArticleStoreIsFavorited(t *testing.T) {
	tests := []struct {
		name            string
		article         *model.Article
		user            *model.User
		mockCountResult int
		mockCountError  error
		want            bool
		wantErr         bool
	}{
		{
			name:            "Article is favorited by the user",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 1,
			mockCountError:  nil,
			want:            true,
			wantErr:         false,
		},
		{
			name:            "Article is not favorited by the user",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 0,
			mockCountError:  nil,
			want:            false,
			wantErr:         false,
		},
		{
			name:            "Nil Article parameter",
			article:         nil,
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 0,
			mockCountError:  nil,
			want:            false,
			wantErr:         false,
		},
		{
			name:            "Nil User parameter",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            nil,
			mockCountResult: 0,
			mockCountError:  nil,
			want:            false,
			wantErr:         false,
		},
		{
			name:            "Database error occurs",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 0,
			mockCountError:  errors.New("database error"),
			want:            false,
			wantErr:         true,
		},
		{
			name:            "Empty database table",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 0,
			mockCountError:  nil,
			want:            false,
			wantErr:         false,
		},
		{
			name:            "Multiple favorites exist",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 5,
			mockCountError:  nil,
			want:            true,
			wantErr:         false,
		},
		{
			name:            "Article exists but not favorited by any user",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 0,
			mockCountError:  nil,
			want:            false,
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{
				countResult: tt.mockCountResult,
				countError:  tt.mockCountError,
			}
			s := &ArticleStore{
				db: mockDB,
			}
			got, err := s.IsFavorited(tt.article, tt.user)
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
	return &gorm.DB{Error: m.countError}
}

