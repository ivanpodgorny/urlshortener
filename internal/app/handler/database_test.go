package handler

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

type PingerMock struct {
	mock.Mock
}

func (m *PingerMock) Ping(_ context.Context) bool {
	args := m.Called()

	return args.Bool(0)
}

func TestDatabase_Ping(t *testing.T) {
	successPinger := &PingerMock{}
	successPinger.On("Ping").Return(true).Once()
	failPinger := &PingerMock{}
	failPinger.On("Ping").Return(false).Once()

	tests := []struct {
		name           string
		wantStatusCode int
		pinger         *PingerMock
	}{
		{
			name:           "успешная проверка",
			wantStatusCode: http.StatusOK,
			pinger:         successPinger,
		},
		{
			name:           "ошибка при проверке соединения",
			wantStatusCode: http.StatusInternalServerError,
			pinger:         failPinger,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := Database{pinger: tt.pinger}
			result := sendTestRequest(http.MethodGet, "/", nil, handler.Ping)
			assert.Equal(t, tt.wantStatusCode, result.StatusCode)
			err := result.Body.Close()
			require.NoError(t, err)
		})
	}
	successPinger.AssertExpectations(t)
	failPinger.AssertExpectations(t)
}
