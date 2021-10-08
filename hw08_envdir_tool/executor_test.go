package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	const dirPath = "testdata/env"
	env, err := ReadDir(dirPath)
	require.NoError(t, err)

	cmd := []string{"/bin/bash", "testdata/echo.sh", "arg1=1", "arg2=2"}
	exitCode := RunCmd(cmd, env)
	require.Equal(t, 0, exitCode)

	cmd = []string{"/bin/bash", "testdata/invalid"}
	exitCode = RunCmd(cmd, env)
	require.Equal(t, 127, exitCode)
}
