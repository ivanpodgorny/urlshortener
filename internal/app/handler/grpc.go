package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ivanpodgorny/urlshortener/internal/proto"
)

// ShortenerServer реализует интерфейс GRPC-сервера сервиса сокращения URL.
type ShortenerServer struct {
	proto.UnimplementedShortenerServer
	authenticator IdentityProvider
	shortener     Shortener
}

// NewShortenerGRPCServer возвращает указатель на новый экземпляр ShortenerServer.
func NewShortenerGRPCServer(a IdentityProvider, s Shortener) *ShortenerServer {
	return &ShortenerServer{
		authenticator: a,
		shortener:     s,
	}
}

// CreateLink обрабатывает запрос на создание сокращенного URL.
func (s *ShortenerServer) CreateLink(ctx context.Context, request *proto.CreateLinkRequest) (*proto.CreateLinkResponse, error) {
	userID, err := s.authenticator.UserIdentifier(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "cannot get user ID")
	}

	id, inserted, err := s.shortener.Shorten(ctx, request.GetUrl(), userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	resp := proto.CreateLinkResponse{Id: id}
	var statusErr error
	if !inserted {
		statusErr = status.Error(codes.AlreadyExists, "url has already been added")
	}

	return &resp, statusErr
}

// CreateLinkBatch обрабатывает запрос на создание нескольких сокращенных URL.
func (s *ShortenerServer) CreateLinkBatch(ctx context.Context, request *proto.CreateLinkBatchRequest) (*proto.CreateLinkBatchResponse, error) {
	userID, err := s.authenticator.UserIdentifier(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "cannot get user ID")
	}

	resp := proto.CreateLinkBatchResponse{Urls: make([]*proto.URLData, 0, len(request.GetUrls()))}
	for _, u := range request.GetUrls() {
		id, _, err := s.shortener.Shorten(ctx, u, userID)
		if err != nil {
			continue
		}

		resp.Urls = append(resp.Urls, &proto.URLData{
			Url: u,
			Id:  id,
		})
	}

	return &resp, nil
}

// GetURL обрабатывает запрос на получение оригинального URL по ID.
func (s *ShortenerServer) GetURL(ctx context.Context, request *proto.GetURLRequest) (*proto.GetURLResponse, error) {
	u, err := s.shortener.Get(ctx, request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", err)
	}

	return &proto.GetURLResponse{Url: u}, nil
}

// GetAllURL возвращает все сокращенные URL пользователя, выполнившего запрос.
func (s *ShortenerServer) GetAllURL(ctx context.Context, _ *proto.GetAllURLRequest) (*proto.GetAllURLResponse, error) {
	userID, err := s.authenticator.UserIdentifier(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "cannot get user ID")
	}

	urls := s.shortener.GetAllUser(ctx, userID)
	resp := proto.GetAllURLResponse{Urls: make([]*proto.URLData, 0, len(urls))}
	for id, u := range urls {
		resp.Urls = append(resp.Urls, &proto.URLData{
			Url: u,
			Id:  id,
		})
	}

	return &resp, nil
}

// DeleteURLBatch выполняет удаление URL по переданным ID.
func (s *ShortenerServer) DeleteURLBatch(ctx context.Context, request *proto.DeleteURLBatchRequest) (*proto.DeleteURLBatchResponse, error) {
	userID, err := s.authenticator.UserIdentifier(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "cannot get user ID")
	}

	if err = s.shortener.DeleteBatch(ctx, request.GetIds(), userID); err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &proto.DeleteURLBatchResponse{}, nil
}
