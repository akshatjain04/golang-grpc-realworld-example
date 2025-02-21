package store

import (
	"errors"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDB struct {
	mock.Mock
	*gorm.DB
	countResult int
	countError  error
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
				assoc := &mockDB{}
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
				assoc := &mockDB{}
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
				assoc := &mockDB{}
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
				db: mockDB.DB,
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
			store := &ArticleStore{
				db: mockDB.DB,
			}
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

func (m *mockDB) Append(values ...interface{}) error {
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

func (m *mockDB) Count(value interface{}) *gorm.DB {
	*value.(*int) = m.countResult
	return &gorm.DB{Error: m.countError}
}

func (m *mockDB) Model(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *mockDB) Rollback() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *mockDB) Table(name string) *gorm.DB {
	return &gorm.DB{Error: m.countError}
}

func (m *mockDB) Update(column string, value interface{}) *gorm.DB {
	args := m.Called(column, value)
	return args.Get(0).(*gorm.DB)
}
