package store

import (
	"testing"

	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	model "github.com/raahii/golang-grpc-realworld-example/model"
)

/*
ROOST_METHOD_HASH=ArticleStore_DeleteFavorite_29c18a04a8
ROOST_METHOD_SIG_HASH=ArticleStore_DeleteFavorite_53deb5e792

FUNCTION_DEF=func (s *ArticleStore) DeleteFavorite(a *model.Article, u *model.User) error // DeleteFavorite unfavorite an article
*/

func TestArticleStoreDeleteFavorite(t *testing.T) {
	type input struct {
		article *model.Article
		user    *model.User
		mocks   func(mock gosqlmock.Sqlmock)
	}

	// additional code goes here...
}
