package security

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"

	"github.com/ivanpodgorny/urlshortener/internal/proto"
)

type MockGRPCServer struct {
	proto.UnimplementedShortenerServer
	authenticator *GRPCAuthenticator
	mock.Mock
}

func (s *MockGRPCServer) GetURL(ctx context.Context, _ *proto.GetURLRequest) (*proto.GetURLResponse, error) {
	id, _ := s.authenticator.UserIdentifier(ctx)
	s.Called(id)

	return &proto.GetURLResponse{}, nil
}

func TestGRPCAuthenticator(t *testing.T) {
	authenticator := NewGRPCAuthenticator(
		NewHMACTokenCreatorParser(""),
		NewGRPCContextUserProvider(),
	)
	ctx := context.Background()
	_, err := authenticator.UserIdentifier(ctx)
	assert.Error(t, err, "неаутентифицированный пользователь")
	ctx = authenticator.Authenticate(ctx)
	id, err := authenticator.UserIdentifier(ctx)
	assert.NoError(t, err, "создание нового токена для пользователя")
	ctx = authenticator.Authenticate(ctx)
	existingID, err := authenticator.UserIdentifier(ctx)
	assert.NoError(t, err, "получение существующего токена пользователя")
	assert.Equal(t, id, existingID, "получение существующего токена пользователя")
}

func TestGRPCAuthenticator_Authenticate(t *testing.T) {
	var (
		lis = bufconn.Listen(1024 * 1024)
		s   = grpc.NewServer()
		m   = &MockGRPCServer{
			authenticator: NewGRPCAuthenticator(
				NewHMACTokenCreatorParser(""),
				NewGRPCContextUserProvider(),
			),
		}
		userID = "userID"
	)
	m.On("GetURL", userID).Once()
	proto.RegisterShortenerServer(s, m)
	go func() {
		_ = s.Serve(lis)
	}()
	time.Sleep(100 * time.Millisecond)

	bufDialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	client := proto.NewShortenerClient(conn)
	ctx = metadata.AppendToOutgoingContext(ctx, userIDCookie, userID)
	_, err = client.GetURL(ctx, &proto.GetURLRequest{})
	require.NoError(t, err)
	m.AssertExpectations(t)
}
