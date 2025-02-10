package handler


import (
	"context"
	"fmt"
	"testing"
	"time"
	"github.com/golang/mock/gomock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"errors"
)








/*
ROOST_METHOD_HASH=CurrentUser_b5d1eae50d
ROOST_METHOD_SIG_HASH=CurrentUser_81bc996773

FUNCTION_DEF=func (h *Handler) CurrentUser(ctx context.Context, req *pb.Empty) (*pb.UserResponse, error) // CurrentUser gets a current user


*/
func TestHandlerCurrentUser(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStore := store.NewMockUserStore(ctrl)
	mockLogger := zerolog.NewMockLogger(ctrl)

	testHandler := &Handler{
		logger: &mockLogger,
		us:     mockUserStore,
	}

	testCases := []struct {
		description    string
		setupMocks     func()
		expectedResult *proto.UserResponse
		expectedError  error
	}{
		{
			description: "Successful retrieval of current user details",
			setupMocks: func() {
				mockUserStore.EXPECT().GetByID(gomock.Any()).Return(&model.User{
					Email:    "test@example.com",
					Username: "testuser",
					Bio:      "A test user",
					Image:    "https://example.com/testuser.jpg",
				}, nil).Times(1)
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 1, nil
				}
				auth.GenerateToken = func(id uint) (string, error) {
					return "valid_token_string", nil
				}
			},
			expectedResult: &proto.UserResponse{
				User: &proto.User{
					Email:    "test@example.com",
					Token:    "valid_token_string",
					Username: "testuser",
					Bio:      "A test user",
					Image:    "https://example.com/testuser.jpg",
				},
			},
			expectedError: nil,
		},
		{
			description: "User ID not found in the context",
			setupMocks: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 0, status.Error(codes.Unauthenticated, "No user ID found in the context")
				}
			},
			expectedResult: nil,
			expectedError:  status.Error(codes.Unauthenticated, "unauthenticated"),
		},
		{
			description: "User not found in the user store",
			setupMocks: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 1, nil
				}
				mockUserStore.EXPECT().GetByID(gomock.Any()).Return(nil, status.Error(codes.NotFound, "User not found")).Times(1)
			},
			expectedResult: nil,
			expectedError:  status.Error(codes.NotFound, "user not found"),
		},
		{
			description: "Token generation failure",
			setupMocks: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 1, nil
				}
				mockUserStore.EXPECT().GetByID(gomock.Any()).Return(&model.User{}, nil).Times(1)
				auth.GenerateToken = func(id uint) (string, error) {
					return "", fmt.Errorf("Token generation failed")
				}
			},
			expectedResult: nil,
			expectedError:  status.Error(codes.Aborted, "internal server error"),
		},
		{
			description: "Expiry of user's token",
			setupMocks: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 0, status.Error(codes.Unauthenticated, "Token expired")
				}
			},
			expectedResult: nil,
			expectedError:  status.Error(codes.Unauthenticated, "unauthenticated"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {

			tc.setupMocks()

			response, err := testHandler.CurrentUser(context.Background(), &proto.Empty{})

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, response)
			}

			if err != nil {
				t.Logf("Test '%s' failed: %v", tc.description, err)
			} else {
				t.Logf("Test '%s' succeeded.", tc.description)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=LoginUser_1f3a5a4bc8
ROOST_METHOD_SIG_HASH=LoginUser_a5eec41841

FUNCTION_DEF=func (h *Handler) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.UserResponse, error) // LoginUser is existing user login


*/
func TestHandlerLoginUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserStore := store.NewMockUserStore(mockCtrl)
	mockLogger := zerolog.NewMockLogger(mockCtrl)
	mockAuth := auth.NewMockAuth(mockCtrl)

	tests := []struct {
		name          string
		setupMock     func()
		req           *proto.LoginUserRequest
		expectedError error
		expectedResp  *proto.UserResponse
	}{
		{
			name: "Successful login with correct credentials",
			setupMock: func() {
				mockUser := &model.User{
					Email:    "test@example.com",
					Password: "$2a$10$examplehashedpassword",
				}
				mockUserStore.EXPECT().GetByEmail("test@example.com").Return(mockUser, nil).Times(1)
				mockAuth.EXPECT().GenerateToken(mockUser.ID).Return("dummyToken", nil).Times(1)
			},
			req: &proto.LoginUserRequest{
				User: &proto.LoginUserRequest_User{
					Email:    "test@example.com",
					Password: "password",
				},
			},
			expectedError: nil,
			expectedResp: &proto.UserResponse{
				User: &proto.User{
					Email:    "test@example.com",
					Token:    "dummyToken",
					Bio:      "",
					Image:    "",
					Username: "",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()
			h := &Handler{
				logger: mockLogger,
				us:     mockUserStore,
				as:     nil,
			}
			resp, err := h.LoginUser(context.Background(), tc.req)

			if tc.expectedError != nil {
				assert.Error(t, err)
				st, _ := status.FromError(err)
				assert.Equal(t, tc.expectedError.(*status.Status).Code(), st.Code())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResp, resp)
			}
		})
	}
}

