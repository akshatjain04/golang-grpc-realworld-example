package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockArticleStore struct {
	errOnGetByID       bool
	errOnCreateComment bool
}

func (m *mockArticleStore) GetByID(uint) (*model.Article, error) {
	if m.errOnGetByID {
		return nil, errors.New("test error")
	}
	return &model.Article{}, nil
}

func (m *mockArticleStore) CreateComment(*model.Comment) error {
	if m.errOnCreateComment {
		return errors.New("test error")
	}
	return nil
}

type mockUserStore struct {
	errOnGetByID bool
}

func (m *mockUserStore) GetByID(uint) (*model.User, error) {
	if m.errOnGetByID {
		return nil, errors.New("test error")
	}
	return &model.User{}, nil
}

func TestHandler_CreateComment(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	req := &pb.CreateCommentRequest{}
	tests := []struct {
		name        string
		setup       func() *Handler
		expectError bool
		expectCode  codes.Code
	}{
		{
			"Scenario 1: Comment Creation with Valid Data",
			func() *Handler {
				return &Handler{
					logger: &zerolog.Logger{},
					us:     &mockUserStore{},
					as:     &mockArticleStore{},
				}
			},
			false,
			codes.OK,
		},
		{
			"Scenario 2: User Not Authenticated",
			func() *Handler {
				return &Handler{
					logger: &zerolog.Logger{},
					us:     &mockUserStore{true},
					as:     &mockArticleStore{},
				}
			},
			true,
			codes.Unauthenticated,
		},
		{
			"Scenario 3: User Not Found",
			func() *Handler {
				return &Handler{
					logger: &zerolog.Logger{},
					us:     &mockUserStore{true},
					as:     &mockArticleStore{},
				}
			},
			true,
			codes.NotFound,
		},
		{
			"Scenario 4: Invalid Article ID",
			func() *Handler {
				return &Handler{
					logger: &zerolog.Logger{},
					us:     &mockUserStore{},
					as:     &mockArticleStore{true},
				}
			},
			true,
			codes.InvalidArgument,
		},
		{
			"Scenario 5: Failed Comment Creation",
			func() *Handler {
				return &Handler{
					logger: &zerolog.Logger{},
					us:     &mockUserStore{},
					as:     &mockArticleStore{false, true},
				}
			},
			true,
			codes.Aborted,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := test.setup()
			_, err := h.CreateComment(ctx, req)
			if test.expectError {
				assert.Error(err)
				s, _ := status.FromError(err)
				assert.Equal(test.expectCode, s.Code())
			} else {
				assert.NoError(err)
			}
		})
	}
}
