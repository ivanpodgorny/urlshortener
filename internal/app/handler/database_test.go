package handler

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
		pinger         *PingerMock
		name           string
		wantStatusCode int
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
