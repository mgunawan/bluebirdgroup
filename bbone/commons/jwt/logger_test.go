package jwt

import (
	"testing"

	"bluebirdgroup/bbone/commons/logger"
)

func TestSetLevel(t *testing.T) {
	logger.Info("Test logger")
	logger.SetLevel(logger.WarnLevel)
	logger.Error("Test logger")
}
