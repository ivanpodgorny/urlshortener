package interceptor

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
}

type AuthenticatorMock struct {
	mock.Mock
}

func (m *AuthenticatorMock) Authenticate(ctx context.Context) context.Context {
	m.Called(ctx)
	_ = grpc.SendHeader(ctx, metadata.Pairs("identity", "1"))

	return ctx
}

func TestAuthenticate(t *testing.T) {
	var (
		lis = bufconn.Listen(1024 * 1024)
		a   = &AuthenticatorMock{}
		s   = grpc.NewServer(grpc.UnaryInterceptor(Authenticate(a)))
		m   = &MockGRPCServer{}
		ctx = context.Background()
	)

	a.On("Authenticate", mock.AnythingOfType("*context.valueCtx")).Return(ctx).Once()
	proto.RegisterShortenerServer(s, m)
	go func() {
		_ = s.Serve(lis)
	}()
	time.Sleep(100 * time.Millisecond)

	bufDialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
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
	var header metadata.MD
	_, _ = client.GetURL(ctx, &proto.GetURLRequest{}, grpc.Header(&header))
	assert.Equal(t, "1", header.Get("identity")[0])
	a.AssertExpectations(t)
}
