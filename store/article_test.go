package store

import (
	errors "errors"
	fmt "fmt"
	debug "runtime/debug"
	testing "testing"

	"github.com/DATA-DOG/go-sqlmock"
	gorm "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	model "github.com/raahii/golang-grpc-realworld-example/model"
)

/*
ROOST_METHOD_HASH=ArticleStore_AddFavorite_9460fca478
ROOST_METHOD_SIG_HASH=ArticleStore_AddFavorite_c13a109f91

FUNCTION_DEF=func (s *ArticleStore) AddFavorite(a *model.Article, u *model.User) error // AddFavorite favorite an article
*/
func TestArticleStoreAddFavorite(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock)
		article       *model.Article
		user          *model.User
		expectedError bool
	}{
		{
			name: "Successfully Adding a Favorite to an Article",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "favorite_articles"`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`UPDATE "articles" SET "favorites_count"`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			article: &model.Article{Model: gorm.Model{ID: 1}, FavoritesCount: 0},
			user:    &model.User{Model: gorm.Model{ID: 1}},
		},
		{
			name: "Adding a Favorite When the User Already Favorited the Article",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "favorite_articles"`).
					WillReturnError(errors.New("duplicate entry"))
				mock.ExpectRollback()
			},
			article:       &model.Article{Model: gorm.Model{ID: 1}, FavoritesCount: 1},
			user:          &model.User{Model: gorm.Model{ID: 1}},
			expectedError: true,
		},
		{
			name: "Adding a Favorite When the Database Transaction Fails",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "favorite_articles"`).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			article:       &model.Article{Model: gorm.Model{ID: 1}, FavoritesCount: 0},
			user:          &model.User{Model: gorm.Model{ID: 1}},
			expectedError: true,
		},
		{
			name: "Adding a Favorite When Updating favorites_count Fails",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "favorite_articles"`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`UPDATE "articles" SET "favorites_count"`).
					WillReturnError(errors.New("update error"))
				mock.ExpectRollback()
			},
			article:       &model.Article{Model: gorm.Model{ID: 1}, FavoritesCount: 0},
			user:          &model.User{Model: gorm.Model{ID: 1}},
			expectedError: true,
		},
		{
			name:          "Adding a Favorite with a Nil Article",
			setupMock:     func(mock sqlmock.Sqlmock) {},
			article:       nil,
			user:          &model.User{Model: gorm.Model{ID: 1}},
			expectedError: true,
		},
		{
			name:          "Adding a Favorite with a Nil User",
			setupMock:     func(mock sqlmock.Sqlmock) {},
			article:       &model.Article{Model: gorm.Model{ID: 1}, FavoritesCount: 0},
			user:          nil,
			expectedError: true,
		},
		{
			name: "Adding a Favorite When the Database Connection is Closed",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("connection closed"))
			},
			article:       &model.Article{Model: gorm.Model{ID: 1}, FavoritesCount: 0},
			user:          &model.User{Model: gorm.Model{ID: 1}},
			expectedError: true,
		},
		{
			name: "Concurrently Adding Favorites to the Same Article",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "favorite_articles"`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`UPDATE "articles" SET "favorites_count"`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			article: &model.Article{Model: gorm.Model{ID: 1}, FavoritesCount: 0},
			user:    &model.User{Model: gorm.Model{ID: 1}},
		},
		{
			name: "Adding a Favorite When the Article Has No FavoritedUsers Association",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "favorite_articles"`).
					WillReturnError(errors.New("missing association"))
				mock.ExpectRollback()
			},
			article:       &model.Article{Model: gorm.Model{ID: 1}, FavoritesCount: 0},
			user:          &model.User{Model: gorm.Model{ID: 1}},
			expectedError: true,
		},
		{
			name: "Adding a Favorite When the Article Already Has Maximum Integer favorites_count",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "favorite_articles"`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`UPDATE "articles" SET "favorites_count"`).
					WillReturnError(errors.New("integer overflow"))
				mock.ExpectRollback()
			},
			article:       &model.Article{Model: gorm.Model{ID: 1}, FavoritesCount: int32(^uint32(0) >> 1)},
			user:          &model.User{Model: gorm.Model{ID: 1}},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock database: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("sqlite3", db)
			if err != nil {
				t.Fatalf("Failed to open gorm database: %v", err)
			}
			defer gormDB.Close()

			store := &ArticleStore{db: gormDB}

			tt.setupMock(mock)

			err = store.AddFavorite(tt.article, tt.user)

			if (err != nil) != tt.expectedError {
				t.Errorf("Expected error: %v, got: %v", tt.expectedError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=ArticleStore_GetArticles_101b7250e8
ROOST_METHOD_SIG_HASH=ArticleStore_GetArticles_91bc0a6760

FUNCTION_DEF=func (s *ArticleStore) GetArticles(tagName, username string, favoritedBy *model.User, limit, offset int64) ([ // GetArticles get global articles
]model.Article, error)
*/
func TestArticleStoreGetArticles(t *testing.T) {
	tests := []struct {
		name          string
		tagName       string
		username      string
		favoritedBy   *model.User
		limit         int64
		offset        int64
		setupMock     func(mock sqlmock.Sqlmock)
		expectedCount int
		expectError   bool
	}{
		{
			name:        "Retrieve Articles Without Filters",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).
						AddRow(1, "Article 1").
						AddRow(2, "Article 2"))
			},
			expectedCount: 2,
			expectError:   false,
		},
		{
			name:        "Retrieve Articles by Specific Username",
			tagName:     "",
			username:    "testuser",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles" JOIN users ON articles.user_id = users.id WHERE users.username = \?`).
					WithArgs("testuser").
					WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).
						AddRow(1, "User's Article"))
			},
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:        "Retrieve Articles by Tag Name",
			tagName:     "tech",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles" JOIN article_tags ON articles.id = article_tags.article_id JOIN tags ON tags.id = article_tags.tag_id WHERE tags.name = \?`).
					WithArgs("tech").
					WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).
						AddRow(1, "Tech Article"))
			},
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:        "Retrieve Articles Favorited by a Specific User",
			tagName:     "",
			username:    "",
			favoritedBy: &model.User{Model: gorm.Model{ID: 1}},
			limit:       10,
			offset:      0,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT article_id FROM "favorite_articles" WHERE user_id = \?`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"article_id"}).
						AddRow(1).
						AddRow(2))
				mock.ExpectQuery(`SELECT \* FROM "articles" WHERE id in \(\?, \?\)`).
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).
						AddRow(1, "Favorited Article 1").
						AddRow(2, "Favorited Article 2"))
			},
			expectedCount: 2,
			expectError:   false,
		},
		{
			name:        "Retrieve Articles When No Matching Records Exist",
			tagName:     "nonexistent",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles" JOIN article_tags ON articles.id = article_tags.article_id JOIN tags ON tags.id = article_tags.tag_id WHERE tags.name = \?`).
					WithArgs("nonexistent").
					WillReturnRows(sqlmock.NewRows([]string{"id", "title"}))
			},
			expectedCount: 0,
			expectError:   false,
		},
		{
			name:        "Handle Database Query Failure",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WillReturnError(fmt.Errorf("database error"))
			},
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock database: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("sqlite3", db)
			if err != nil {
				t.Fatalf("Failed to open gorm database: %v", err)
			}

			store := &ArticleStore{db: gormDB}

			tt.setupMock(mock)

			articles, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}

			if len(articles) != tt.expectedCount {
				t.Errorf("Expected %d articles, got %d", tt.expectedCount, len(articles))
			}
		})
	}
}

/*
ROOST_METHOD_HASH=ArticleStore_IsFavorited_799826fee5
ROOST_METHOD_SIG_HASH=ArticleStore_IsFavorited_f6d5e67492

FUNCTION_DEF=func (s *ArticleStore) IsFavorited(a *model.Article, u *model.User) (bool, error) // IsFavorited returns whether the article is favorited by the user
*/
func TestArticleStoreIsFavorited(t *testing.T) {
	tests := []struct {
		name          string
		article       *model.Article
		user          *model.User
		mockBehavior  func(mock sqlmock.Sqlmock)
		expected      bool
		expectedError error
	}{
		{
			name:    "Article is favorited by the user",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			user:    &model.User{Model: gorm.Model{ID: 1}},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count(.+) FROM "favorite_articles" WHERE article_id = \? AND user_id = \?`).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expected:      true,
			expectedError: nil,
		},
		{
			name:    "Article is not favorited by the user",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			user:    &model.User{Model: gorm.Model{ID: 1}},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count(.+) FROM "favorite_articles" WHERE article_id = \? AND user_id = \?`).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expected:      false,
			expectedError: nil,
		},
		{
			name:          "Article is nil",
			article:       nil,
			user:          &model.User{Model: gorm.Model{ID: 1}},
			mockBehavior:  func(mock sqlmock.Sqlmock) {},
			expected:      false,
			expectedError: nil,
		},
		{
			name:          "User is nil",
			article:       &model.Article{Model: gorm.Model{ID: 1}},
			user:          nil,
			mockBehavior:  func(mock sqlmock.Sqlmock) {},
			expected:      false,
			expectedError: nil,
		},
		{
			name:          "Both article and user are nil",
			article:       nil,
			user:          nil,
			mockBehavior:  func(mock sqlmock.Sqlmock) {},
			expected:      false,
			expectedError: nil,
		},
		{
			name:    "Database query fails",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			user:    &model.User{Model: gorm.Model{ID: 1}},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count(.+) FROM "favorite_articles" WHERE article_id = \? AND user_id = \?`).
					WithArgs(1, 1).
					WillReturnError(fmt.Errorf("database error"))
			},
			expected:      false,
			expectedError: fmt.Errorf("database error"),
		},
		{
			name:    "User has favorited multiple articles",
			article: &model.Article{Model: gorm.Model{ID: 2}},
			user:    &model.User{Model: gorm.Model{ID: 1}},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count(.+) FROM "favorite_articles" WHERE article_id = \? AND user_id = \?`).
					WithArgs(2, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expected:      true,
			expectedError: nil,
		},
		{
			name:    "Multiple users have favorited the same article",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			user:    &model.User{Model: gorm.Model{ID: 2}},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count(.+) FROM "favorite_articles" WHERE article_id = \? AND user_id = \?`).
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expected:      true,
			expectedError: nil,
		},
		{
			name:    "User has favorited other articles but not the given one",
			article: &model.Article{Model: gorm.Model{ID: 3}},
			user:    &model.User{Model: gorm.Model{ID: 1}},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count(.+) FROM "favorite_articles" WHERE article_id = \? AND user_id = \?`).
					WithArgs(3, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expected:      false,
			expectedError: nil,
		},
		{
			name:    "Database returns zero count",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			user:    &model.User{Model: gorm.Model{ID: 1}},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count(.+) FROM "favorite_articles" WHERE article_id = \? AND user_id = \?`).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expected:      false,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock database: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("Failed to open gorm database: %v", err)
			}
			defer gormDB.Close()

			store := &ArticleStore{db: gormDB}

			tt.mockBehavior(mock)

			result, err := store.IsFavorited(tt.article, tt.user)

			if result != tt.expected || (err != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("Test %s failed: expected (%v, %v), got (%v, %v)", tt.name, tt.expected, tt.expectedError, result, err)
			} else {
				t.Logf("Test %s passed successfully", tt.name)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}
