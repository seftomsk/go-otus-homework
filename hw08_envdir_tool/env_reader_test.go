package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const dirPath = "testdata/env"

func getEnv() Environment {
	env := make(Environment)
	env["BAR"] = EnvValue{"bar", false}
	env["EMPTY"] = EnvValue{"", false}
	env["FOO"] = EnvValue{"   foo\nwith new line", false}
	env["HELLO"] = EnvValue{"\"hello\"", false}
	env["UNSET"] = EnvValue{"", true}

	return env
}

func TestReadDir(t *testing.T) {
	env := getEnv()
	actual, err := ReadDir(dirPath)
	require.NoError(t, err)
	require.Equal(t, env, actual)
}

func TestInvalidDir(t *testing.T) {
	t.Run("empty path", func(t *testing.T) {
		_, err := ReadDir("")
		require.ErrorIs(t, err, ErrEmptyDir)
	})

	t.Run("invalid path", func(t *testing.T) {
		_, err := ReadDir("path")
		require.Error(t, err)
	})
}

func TestSkippingFileWithEqualInName(t *testing.T) {
	file, err := os.CreateTemp(dirPath, "test=")
	require.NoError(t, err)
	defer func() {
		err = file.Close()
		require.NoError(t, err)
		err = os.Remove(file.Name())
		require.NoError(t, err)
	}()
	env := getEnv()
	actual, err := ReadDir(dirPath)
	require.NoError(t, err)
	require.Equal(t, env, actual)
}
