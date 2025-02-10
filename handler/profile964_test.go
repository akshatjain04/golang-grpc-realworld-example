package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
	"github.com/golang/mock/gomock"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/raahii/golang-grpc-realworld-example/model"
)





type MockUserStore struct {
	mock.Mock
}


/*
ROOST_METHOD_HASH=FollowUser_d8ce9732de
ROOST_METHOD_SIG_HASH=FollowUser_2fec5af0c0

FUNCTION_DEF=func (h *Handler) FollowUser(ctx context.Context, req *pb.FollowRequest) (*pb.ProfileResponse, error) // FollowUser follow a user


*/
func (m *MockUserStore) Follow(a *model.User, b *model.User) error {
	args := m.Called(a, b)
	return args.Error(0)
}

func (m *MockUserStore) GetByID(id uint) (*model.User, error) {
	args := m.Called(id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserStore) GetByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	return args.Get(0).(*model.User), args.Error(1)
}

func TestHandlerFollowUser(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*MockUserStore)
		req            *proto.FollowRequest
		expectedResult *proto.ProfileResponse
		expectedError  error
	}{
		{
			name: "Successful FollowUser operation",
			setupMocks: func(m *MockUserStore) {
				m.On("GetByID", uint(1)).Return(&model.User{ID: 1, Username: "user1"}, nil)
				m.On("GetByUsername", "user2").Return(&model.User{ID: 2, Username: "user2"}, nil)
				m.On("Follow", mock.Anything, mock.Anything).Return(nil)
			},
			req:            &proto.FollowRequest{Username: "user2"},
			expectedResult: &proto.ProfileResponse{Profile: &proto.Profile{Username: "user2", Following: true}},
			expectedError:  nil,
		},
		{
			name: "FollowUser with unauthenticated user",
			setupMocks: func(m *MockUserStore) {

			},
			req:            &proto.FollowRequest{Username: "user2"},
			expectedResult: nil,
			expectedError:  status.Errorf(codes.Unauthenticated, "unauthenticated"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore := new(MockUserStore)
			tt.setupMocks(mockUserStore)
			logger := zerolog.New(nil)
			h := &Handler{
				logger: &logger,
				us:     mockUserStore,
			}

			ctx := context.Background()

			result, err := h.FollowUser(ctx, tt.req)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			mockUserStore.AssertExpectations(t)
		})
	}
}


/*
ROOST_METHOD_HASH=ShowProfile_462fa811fd
ROOST_METHOD_SIG_HASH=ShowProfile_36c21bedd3

FUNCTION_DEF=func (h *Handler) ShowProfile(ctx context.Context, req *pb.ShowProfileRequest) (*pb.ProfileResponse, error) // ShowProfile gets a profile


*/
func TestHandlerShowProfile(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserStore := store.NewMockUserStore(mockCtrl)
	mockLogger := zerolog.New(nil)
	handler := &Handler{logger: &mockLogger, us: mockUserStore}

	testCases := []struct {
		name      string
		setupMock func()
		req       *proto.ShowProfileRequest
		wantErr   bool
		errCode   codes.Code
	}{
		{
			name: "Unauthenticated request",
			setupMock: func() {

			},
			req:     &proto.ShowProfileRequest{Username: "testuser"},
			wantErr: true,
			errCode: codes.Unauthenticated,
		},
		{
			name: "Current user not found in the system",
			setupMock: func() {
				mockUserStore.EXPECT().GetByID(gomock.Any()).Return(nil, fmt.Errorf("user not found"))
			},
			req:     &proto.ShowProfileRequest{Username: "testuser"},
			wantErr: true,
			errCode: codes.NotFound,
		},
		{
			name: "Requested user not found in the system",
			setupMock: func() {
				mockUserStore.EXPECT().GetByID(gomock.Any()).Return(&store.User{}, nil)
				mockUserStore.EXPECT().GetByUsername("nonexistent").Return(nil, fmt.Errorf("user not found"))
			},
			req:     &proto.ShowProfileRequest{Username: "nonexistent"},
			wantErr: true,
			errCode: codes.NotFound,
		},
		{
			name: "Error when checking following status",
			setupMock: func() {
				mockUserStore.EXPECT().GetByID(gomock.Any()).Return(&store.User{}, nil)
				mockUserStore.EXPECT().GetByUsername("testuser").Return(&store.User{}, nil)
				mockUserStore.EXPECT().IsFollowing(gomock.Any(), gomock.Any()).Return(false, fmt.Errorf("database error"))
			},
			req:     &proto.ShowProfileRequest{Username: "testuser"},
			wantErr: true,
			errCode: codes.Internal,
		},
		{
			name: "Successful retrieval of a user profile",
			setupMock: func() {
				mockUserStore.EXPECT().GetByID(gomock.Any()).Return(&store.User{}, nil)
				mockUserStore.EXPECT().GetByUsername("testuser").Return(&store.User{}, nil)
				mockUserStore.EXPECT().IsFollowing(gomock.Any(), gomock.Any()).Return(true, nil)
			},
			req:     &proto.ShowProfileRequest{Username: "testuser"},
			wantErr: false,
			errCode: codes.OK,
		},
		{
			name: "Following status between the same user",
			setupMock: func() {
				user := &store.User{}
				mockUserStore.EXPECT().GetByID(gomock.Any()).Return(user, nil)
				mockUserStore.EXPECT().GetByUsername("testuser").Return(user, nil)

			},
			req:     &proto.ShowProfileRequest{Username: "testuser"},
			wantErr: false,
			errCode: codes.OK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			ctx := context.Background()

			resp, err := handler.ShowProfile(ctx, tc.req)
			if tc.wantErr {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tc.errCode, st.Code())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=UnfollowUser_6114d3d2b5
ROOST_METHOD_SIG_HASH=UnfollowUser_44ddaba28f

FUNCTION_DEF=func (h *Handler) UnfollowUser(ctx context.Context, req *pb.UnfollowRequest) (*pb.ProfileResponse, error) // UnfollowUser unfollow a user


*/
func TestHandlerUnfollowUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStore := store.NewMockUserStore(ctrl)
	logger := zerolog.Nop()

	h := &Handler{
		logger: &logger,
		us:     mockUserStore,
	}

	tests := []struct {
		name          string
		setupMocks    func()
		context       context.Context
		request       *proto.UnfollowRequest
		wantProfile   *proto.ProfileResponse
		wantErr       bool
		expectedError error
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			profile, err := h.UnfollowUser(tt.context, tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.Equal(t, tt.expectedError, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantProfile, profile)
			}
		})
	}
}

