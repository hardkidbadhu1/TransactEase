package utils

import (
	"context"
	"github.com/sirupsen/logrus"
	"transact-api/constants"
)

func GetLogger(c context.Context) *logrus.Entry {
	if c == nil {
		return logrus.NewEntry(logrus.New()).WithContext(c).WithField("Application", "transact-api")
	}

	if logger, ok := c.Value(constants.Logger).(*logrus.Entry); ok {
		return logger
	}

	return logrus.NewEntry(logrus.New()).WithContext(c).WithField("Application", "transact-api")
}
