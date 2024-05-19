package utils

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"transact-api/constants"
)

func TestGetLogger(t *testing.T) {
	logger := GetLogger(nil)
	assert.NotNil(t, logger)
	assert.Equal(t, "transact-api", logger.Data["Application"])

	ctx := context.Background()
	logger = GetLogger(ctx)
	assert.NotNil(t, logger)
	assert.Equal(t, "transact-api", logger.Data["Application"])

	expectedLogger := logrus.NewEntry(logrus.New()).WithField("Application", "transact-api")
	ctx = context.WithValue(ctx, constants.Logger, expectedLogger)
	logger = GetLogger(ctx)
	assert.NotNil(t, logger)
	assert.Equal(t, expectedLogger, logger)
}
