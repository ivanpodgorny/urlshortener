package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ivanpodgorny/urlshortener/internal/proto"
)

func TestShortenerServer_CreateLink(t *testing.T) {
	var (
		userID        = "userID"
		url           = "url"
		dupURL        = "dupURL"
		errURL        = "errURL"
		id            = "id"
		dupID         = "dupID"
		ctx           = context.Background()
		authenticator = &AuthenticatorMock{}
		shortener     = &ShortenerMock{}
	)
	authenticator.On("UserIdentifier").Return(userID, nil).Times(3)
	shortener.On("Shorten", url, userID).Return(id, true, nil).Once()
	shortener.On("Shorten", dupURL, userID).Return(dupID, false, nil).Once()
	shortener.On("Shorten", errURL, userID).Return("", false, errors.New("")).Once()
	server := ShortenerServer{
		authenticator: authenticator,
		shortener:     shortener,
	}

	resp, err := server.CreateLink(ctx, &proto.CreateLinkRequest{Url: url})
	assert.NoError(t, err)
	assert.Equal(t, id, resp.GetId())
	resp, err = server.CreateLink(ctx, &proto.CreateLinkRequest{Url: dupURL})
	testGRPCErrorCode(t, err, codes.AlreadyExists)
	assert.Equal(t, dupID, resp.GetId())
	_, err = server.CreateLink(ctx, &proto.CreateLinkRequest{Url: errURL})
	testGRPCErrorCode(t, err, codes.Internal)
	authenticator.AssertExpectations(t)
	shortener.AssertExpectations(t)
}

func TestShortenerServer_CreateLinkBatch(t *testing.T) {
	var (
		userID        = "userID"
		url           = "url"
		dupURL        = "dupURL"
		errURL        = "errURL"
		id            = "id"
		dupID         = "dupID"
		authenticator = &AuthenticatorMock{}
		shortener     = &ShortenerMock{}
	)
	authenticator.On("UserIdentifier").Return(userID, nil).Once()
	shortener.On("Shorten", url, userID).Return(id, true, nil).Once()
	shortener.On("Shorten", dupURL, userID).Return(dupID, false, nil).Once()
	shortener.On("Shorten", errURL, userID).Return("", false, errors.New("")).Once()
	server := ShortenerServer{
		authenticator: authenticator,
		shortener:     shortener,
	}

	resp, err := server.CreateLinkBatch(context.Background(), &proto.CreateLinkBatchRequest{Urls: []string{url, dupURL, errURL}})
	assert.NoError(t, err)
	assert.ElementsMatch(t, []*proto.URLData{
		{
			Url: url,
			Id:  id,
		},
		{
			Url: dupURL,
			Id:  dupID,
		},
	}, resp.GetUrls())
	authenticator.AssertExpectations(t)
	shortener.AssertExpectations(t)
}

func TestShortenerServer_DeleteURLBatch(t *testing.T) {
	var (
		userID        = "userID"
		id            = "id"
		errID         = "errID"
		ctx           = context.Background()
		authenticator = &AuthenticatorMock{}
		shortener     = &ShortenerMock{}
	)
	authenticator.On("UserIdentifier").Return(userID, nil).Twice()
	shortener.On("DeleteBatch", []string{id}, userID).Return(nil).Once()
	shortener.On("DeleteBatch", []string{errID}, userID).Return(errors.New("")).Once()
	server := ShortenerServer{
		authenticator: authenticator,
		shortener:     shortener,
	}

	_, err := server.DeleteURLBatch(ctx, &proto.DeleteURLBatchRequest{Ids: []string{id}})
	assert.NoError(t, err)
	_, err = server.DeleteURLBatch(ctx, &proto.DeleteURLBatchRequest{Ids: []string{errID}})
	testGRPCErrorCode(t, err, codes.Internal)
	authenticator.AssertExpectations(t)
	shortener.AssertExpectations(t)
}

func TestShortenerServer_GetAllURL(t *testing.T) {
	var (
		userID        = "userID"
		url           = "url"
		secURL        = "secURL"
		id            = "id"
		secID         = "secID"
		authenticator = &AuthenticatorMock{}
		shortener     = &ShortenerMock{}
	)
	authenticator.On("UserIdentifier").Return(userID, nil).Once()
	shortener.On("GetAllUser", userID).Return(map[string]string{id: url, secID: secURL}).Once()
	server := ShortenerServer{
		authenticator: authenticator,
		shortener:     shortener,
	}

	resp, err := server.GetAllURL(context.Background(), &proto.GetAllURLRequest{})
	assert.NoError(t, err)
	assert.ElementsMatch(t, []*proto.URLData{
		{
			Url: url,
			Id:  id,
		},
		{
			Url: secURL,
			Id:  secID,
		},
	}, resp.GetUrls())
	authenticator.AssertExpectations(t)
	shortener.AssertExpectations(t)
}

func TestShortenerServer_GetURL(t *testing.T) {
	var (
		url       = "url"
		id        = "id"
		errID     = "errID"
		ctx       = context.Background()
		shortener = &ShortenerMock{}
	)
	shortener.On("Get", id).Return(url, nil).Once()
	shortener.On("Get", errID).Return("", errors.New("")).Once()
	server := ShortenerServer{
		shortener: shortener,
	}

	resp, err := server.GetURL(ctx, &proto.GetURLRequest{Id: id})
	assert.NoError(t, err)
	assert.Equal(t, url, resp.GetUrl())
	_, err = server.GetURL(ctx, &proto.GetURLRequest{Id: errID})
	testGRPCErrorCode(t, err, codes.NotFound)
	shortener.AssertExpectations(t)
}

func TestGRPCUserAuthenticationErrors(t *testing.T) {
	var (
		ctx           = context.Background()
		authenticator = &AuthenticatorMock{}
	)
	authenticator.On("UserIdentifier").Return("", errors.New("")).Times(4)
	server := ShortenerServer{
		authenticator: authenticator,
	}

	_, err := server.CreateLink(ctx, &proto.CreateLinkRequest{})
	testGRPCErrorCode(t, err, codes.PermissionDenied)
	_, err = server.CreateLinkBatch(ctx, &proto.CreateLinkBatchRequest{})
	testGRPCErrorCode(t, err, codes.PermissionDenied)
	_, err = server.GetAllURL(ctx, &proto.GetAllURLRequest{})
	testGRPCErrorCode(t, err, codes.PermissionDenied)
	_, err = server.DeleteURLBatch(ctx, &proto.DeleteURLBatchRequest{})
	testGRPCErrorCode(t, err, codes.PermissionDenied)

	authenticator.AssertExpectations(t)
}

func testGRPCErrorCode(t *testing.T, err error, code codes.Code) {
	s, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, code, s.Code())
}
