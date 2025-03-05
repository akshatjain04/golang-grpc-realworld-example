package store

import (
	sql "database/sql"
	errors "errors"
	testing "testing"

	gorm "github.com/jinzhu/gorm"
	model "github.com/raahii/golang-grpc-realworld-example/model"
	assert "github.com/stretchr/testify/assert"
)

type mockDB struct {
	*gorm.DB
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
		mockDB      func() *gorm.DB
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
			mockDB: func() *gorm.DB {
				db := &mockDB{&gorm.DB{}}
				db.DB.Error = nil
				return db.DB
			},
			want:    []model.Article{{Title: "Article 1"}, {Title: "Article 2"}},
			wantErr: false,
		},
		{
			name:        "Get Articles by Tag Name",
			tagName:     "golang",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockDB: func() *gorm.DB {
				db := &mockDB{&gorm.DB{}}
				db.DB.Error = nil
				return db.DB
			},
			want:    []model.Article{{Title: "Golang Article"}},
			wantErr: false,
		},
		{
			name:        "Get Articles by Author Username",
			tagName:     "",
			username:    "johndoe",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockDB: func() *gorm.DB {
				db := &mockDB{&gorm.DB{}}
				db.DB.Error = nil
				return db.DB
			},
			want:    []model.Article{{Title: "John's Article"}},
			wantErr: false,
		},
		{
			name:        "Get Favorited Articles",
			tagName:     "",
			username:    "",
			favoritedBy: &model.User{Model: gorm.Model{ID: 1}},
			limit:       10,
			offset:      0,
			mockDB: func() *gorm.DB {
				db := &mockDB{&gorm.DB{}}
				db.DB.Error = nil
				return db.DB
			},
			want:    []model.Article{{Title: "Favorited Article"}},
			wantErr: false,
		},
		{
			name:        "Combine Multiple Filters",
			tagName:     "golang",
			username:    "johndoe",
			favoritedBy: &model.User{Model: gorm.Model{ID: 1}},
			limit:       10,
			offset:      0,
			mockDB: func() *gorm.DB {
				db := &mockDB{&gorm.DB{}}
				db.DB.Error = nil
				return db.DB
			},
			want:    []model.Article{{Title: "John's Favorited Golang Article"}},
			wantErr: false,
		},
		{
			name:        "Handle Empty Result Set",
			tagName:     "nonexistent",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockDB: func() *gorm.DB {
				db := &mockDB{&gorm.DB{}}
				db.DB.Error = nil
				return db.DB
			},
			want:    []model.Article{},
			wantErr: false,
		},
		{
			name:        "Test Pagination Limits",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       1000,
			offset:      9999,
			mockDB: func() *gorm.DB {
				db := &mockDB{&gorm.DB{}}
				db.DB.Error = nil
				return db.DB
			},
			want:    []model.Article{},
			wantErr: false,
		},
		{
			name:        "Error Handling for Database Failures",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			mockDB: func() *gorm.DB {
				db := &mockDB{&gorm.DB{}}
				db.DB.Error = errors.New("database error")
				return db.DB
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ArticleStore{
				db: tt.mockDB(),
			}
			got, err := s.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
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

func (m *mockDB) Offset(offset interface{}) *gorm.DB {
	return m.DB
}

func (m *mockDB) Preload(column string, conditions ...interface{}) *gorm.DB {
	return m.DB
}

func (m *mockDB) Rows() (*sql.Rows, error) {
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
