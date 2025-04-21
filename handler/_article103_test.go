package handler

import (
	context "context"
	fmt "fmt"
	os "os"
	strconv "strconv"
	testing "testing"
	github_com_raahii_golang_grpc_realworld_example_handler "github.com/raahii/golang-grpc-realworld-example/handler"
	model "github.com/raahii/golang-grpc-realworld-example/model"
	proto "github.com/raahii/golang-grpc-realworld-example/proto"
	store "github.com/raahii/golang-grpc-realworld-example/store"
	zerolog "github.com/rs/zerolog"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	time "time"
	auth "github.com/raahii/golang-grpc-realworld-example/auth"
)








/*
ROOST_METHOD_HASH=Handler_GetArticle_223c0423fb
ROOST_METHOD_SIG_HASH=Handler_GetArticle_c243313093

FUNCTION_DEF=func (h *Handler) GetArticle(ctx context.Context, req *pb.GetArticleRequest) (*pb.ArticleResponse, error) // GetArticle gets a article


*/
func TestHandler_GetArticle(t *testing.T) {
	mockUserStore := new(store.MockUserStore)
	mockArticleStore := new(store.MockArticleStore)
	articleLogger := zerolog.New(os.Stderr)
	github_com_raahii_golang_grpc_realworld_example_handler := &github_com_raahii_golang_grpc_realworld_example_handler.Handler{logger: &articleLogger, us: mockUserStore, as: mockArticleStore}

	articleTestCases := []struct {
		testName       string
		input          string
		setupMocks     func()
		expectedResult *proto.ArticleResponse
		expectedError  error
	}{
		{
			testName: "Success - Retrieve Existing Article",
			input:    "1",
			setupMocks: func() {
				mockId, _ := strconv.Atoi("1")
				mockArticle := &model.Article{ID: uint(mockId)}
				mockArticleStore.On("GetByID", uint(mockId)).Return(mockArticle, nil)
			},
			expectedResult: &proto.ArticleResponse{Article: &proto.Article{Slug: "1"}},
			expectedError:  nil,
		},
		{
			testName:      "Failure - Unparsable Slug",
			input:         "unparsable_slug",
			expectedError: status.Error(codes.InvalidArgument, "invalid article id"),
		},
		{
			testName: "Failure - Nonexistent Article",
			input:    "10000",
			setupMocks: func() {
				mockId, _ := strconv.Atoi("10000")
				mockArticleStore.On("GetByID", uint(mockId)).Return(nil, model.ErrNoRecord)
			},
			expectedError: status.Error(codes.NotFound, "article not found"),
		},
		{
			testName: "Failure - Database Error",
			input:    "199",
			setupMocks: func() {
				mockId, _ := strconv.Atoi("199")
				mockArticleStore.On("GetByID", uint(mockId)).Return(nil, fmt.Errorf("database error"))
			},
			expectedError: fmt.Errorf("database error"),
		},
	}

	for _, tt := range articleTestCases {
		t.Run(tt.testName, func(t *testing.T) {
			if tt.setupMocks != nil {
				tt.setupMocks()
			}

			ctx := context.Background()
			resp, err := github_com_raahii_golang_grpc_realworld_example_handler.GetArticle(ctx, &proto.GetArticleRequest{Slug: tt.input})

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResult, resp)
				mockArticleStore.AssertExpectations(t)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=Handler_CreateArticle_e5cc3b252e
ROOST_METHOD_SIG_HASH=Handler_CreateArticle_ce1c125740

FUNCTION_DEF=func (h *Handler) CreateArticle(ctx context.Context, req *pb.CreateAritcleRequest) (*pb.ArticleResponse, error) 

*/
func TestCreateArticle_WithOtherScenarios(t *testing.T) {
	h, cleaner := setUp(t)
	defer cleaner(t)
	barUser := model.User{Username: "bar", Email: "bar@example.com", Password: "secret"}
	nonExistingUser := model.User{ID: 99999, Username: "nonExist", Email: "nonExist@example.com", Password: "secret"}
	for _, u := range []*model.User{&barUser} {
		if err := h.us.Create(u); err != nil {
			t.Fatalf("failed to create initial user record: %v", err)
		}
	}
	tests := []struct {
		title    string
		reqUser  *model.User
		req      *proto.CreateAritcleRequest
		wantCode codes.Code
	}{
		{"unauthenticated user attempting to create an article", nil, &proto.CreateAritcleRequest{Article: &proto.CreateAritcleRequest_Article{Title: "Unauthorized post", Description: "Unauthorized description", Body: "Unauthorized content", TagList: []string{"foo", "bar", "piyo"}}}, codes.Unauthenticated},
		{"attempt to create an article with a non-existing user", &nonExistingUser, &proto.CreateAritcleRequest{Article: &proto.CreateAritcleRequest_Article{Title: "Nonexistent post", Description: "Nonexistent description", Body: "Nonexistent content", TagList: []string{"foo", "bar", "piyo"}}}, codes.NotFound},
		{"attempt to create an article with invalid fields", &barUser, &proto.CreateAritcleRequest{Article: &proto.CreateAritcleRequest_Article{Title: "Invalid", Description: "Invalid", Body: "", TagList: []string{}}}, codes.InvalidArgument},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			requestTime := time.Now().Unix() - 1
			ctx := context.Background()
			if tt.reqUser != nil {
				token, err := auth.GenerateToken(tt.reqUser.ID)
				if err != nil {
					t.Error(err)
				}
				ctx = ctxWithToken(ctx, token)
			}
			_, err := h.CreateArticle(ctx, tt.req)
			st, ok := status.FromError(err)
			if ok {
				assert.Equal(t, tt.wantCode, st.Code())
			} else {
				t.Error("Returned error code does not match")
			}
			if err == nil {
				t.Errorf("%q expected to fail, but succeeded.", tt.title)
				t.FailNow()
			}
		})
	}
}


/*
ROOST_METHOD_HASH=Handler_GetArticles_ab2117cfa1
ROOST_METHOD_SIG_HASH=Handler_GetArticles_5d9fe7bf44

FUNCTION_DEF=func (h *Handler) GetArticles(ctx context.Context, req *pb.GetArticlesRequest) (*pb.ArticlesResponse, error) 

*/
func (m *mockArticleStore) GetArticles(tagName, username string, favoritedBy *model.User, limit, offset int64) ([]model.Article, error) {
	args := m.Called(tagName, username, favoritedBy, limit, offset)
	return args.Get(0).([]model.Article), args.Error(1)
}

func (m *mockArticleStore) IsFavorited(a *model.Article, u *model.User) (bool, error) {
	args := m.Called(a, u)
	return args.Bool(0), args.Error(1)
}

func (m *mockUserStore) GetByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *mockUserStore) GetByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *mockUserStore) IsFollowing(a *model.User, b *model.User) (bool, error) {
	args := m.Called(a, b)
	return args.Bool(0), args.Error(1)
}

