package handler

import (
	"context"
	"errors"
	"testing"
	"github.com/golang/mock/gomock"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)








/*
ROOST_METHOD_HASH=GetTags_f169b54062
ROOST_METHOD_SIG_HASH=GetTags_52f72598a3

FUNCTION_DEF=func (h *Handler) GetTags(ctx context.Context, req *pb.Empty) (*pb.TagsResponse, error) 

*/
func TestHandlerGetTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArticleStore := NewMockArticleStore(ctrl)
	mockLogger := zerolog.New(nil)

	h := &Handler{
		logger: &mockLogger,
		as:     mockArticleStore,
	}

	tests := []struct {
		name          string
		setupMock     func()
		expectedTags  []string
		expectedError error
	}{
		{
			name: "Successful retrieval of tags",
			setupMock: func() {
				mockArticleStore.EXPECT().GetTags().Return([]model.Tag{{Name: "tag1"}, {Name: "tag2"}}, nil)
			},
			expectedTags:  []string{"tag1", "tag2"},
			expectedError: nil,
		},
		{
			name: "ArticleStore returns an error",
			setupMock: func() {
				mockArticleStore.EXPECT().GetTags().Return(nil, errors.New("internal error"))
			},
			expectedTags:  nil,
			expectedError: status.Error(codes.Aborted, "internal server error"),
		},
		{
			name: "Zero tags exist in the store",
			setupMock: func() {
				mockArticleStore.EXPECT().GetTags().Return([]model.Tag{}, nil)
			},
			expectedTags:  []string{},
			expectedError: nil,
		},

		{
			name: "Logging behavior on successful retrieval",
			setupMock: func() {
				mockArticleStore.EXPECT().GetTags().Return([]model.Tag{{Name: "tag1"}}, nil)

			},
			expectedTags:  []string{"tag1"},
			expectedError: nil,
		},
		{
			name: "Logging behavior on error",
			setupMock: func() {
				mockArticleStore.EXPECT().GetTags().Return(nil, errors.New("internal error"))

			},
			expectedTags:  nil,
			expectedError: status.Error(codes.Aborted, "internal server error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			resp, err := h.GetTags(context.Background(), &proto.Empty{})

			if tt.expectedError != nil {
				assert.Error(t, err)
				if s, ok := status.FromError(err); ok {
					assert.Equal(t, s.Code(), tt.expectedError.(*status.Status).Code())
					assert.Equal(t, s.Message(), tt.expectedError.(*status.Status).Message())
				}
			} else {
				assert.NoError(t, err)
			}

			if resp != nil {
				assert.Equal(t, tt.expectedTags, resp.Tags)
			}
		})
	}
}

