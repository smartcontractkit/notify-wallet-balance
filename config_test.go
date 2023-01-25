package main

import (
	"testing"

	"github.com/test-go/testify/require"
)

func TestConfigRun(t *testing.T) {
	err := loadConfig()
	require.NoError(t, err, "Error loading config")
}
