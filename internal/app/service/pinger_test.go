package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPinger_Ping(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)
	pinger := Pinger{db: db}
	mock.ExpectPing().WillDelayFor(2 * time.Second)
	assert.Equal(t, false, pinger.Ping(ctx), "превышен таймаут")
	mock.ExpectPing().WillReturnError(errors.New(""))
	assert.Equal(t, false, pinger.Ping(ctx), "ping с ошибкой")
	mock.ExpectPing().WillReturnError(nil)
	assert.Equal(t, true, pinger.Ping(ctx), "успешная проверка")
	assert.NoError(t, mock.ExpectationsWereMet())
}
