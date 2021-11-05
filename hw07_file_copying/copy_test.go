package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const testFolder = "testdata"

const srcFileName = "input.txt"

const dstFileName = "out.txt"

var srcPath = filepath.Join(testFolder, srcFileName)

var dstPath = filepath.Join(testFolder, dstFileName)

func TestCopy(t *testing.T) {
	srcFile, err := os.Open(srcPath)
	require.NoError(t, err)

	testCases := []struct {
		name   string
		offset int64
		limit  int64
	}{
		{"offset0limit0", 0, 0},
		{"offset0limit10", 0, 10},
		{"offset0limit1000", 0, 1000},
		{"offset0limit10000", 0, 10000},
		{"offset100limit100", 100, 100},
		{"offset100limit1000", 100, 1000},
		{"offset6000limit1000", 6000, 1000},
		{"offset-1limit-1", -1, -1},
		{"offset-1limit0", -1, 0},
		{"offset0limit-1", 0, -1},
		{"offset-1limit1", -1, 1},
		{"offset1limit-1", 1, -1},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err = Copy(srcFile.Name(), dstPath, tc.offset, tc.limit)
			require.NoError(t, err)

			err = os.Remove(dstPath)
			require.NoError(t, err)
		})
	}
	defer func() {
		err = srcFile.Close()
		require.NoError(t, err)
	}()
}

func TestOffsetIsMoreOrEqualThanFileSize(t *testing.T) {
	srcFile, err := os.Open(srcPath)
	require.NoError(t, err)

	srcFileStat, err := srcFile.Stat()
	require.NoError(t, err)

	err = Copy(srcFile.Name(), dstPath, srcFileStat.Size()+1, 0)
	require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	err = Copy(srcFile.Name(), dstPath, srcFileStat.Size(), 0)
	require.ErrorIs(t, err, ErrOffsetExceedsFileSize)

	defer func() {
		err = srcFile.Close()
		require.NoError(t, err)
	}()
}

func TestPermissionDeniedForDstFileIfItAlreadyExists(t *testing.T) {
	srcFile, err := os.Open(srcPath)
	require.NoError(t, err)

	readOnlyFile, err := os.Create(dstPath)
	require.NoError(t, err)
	readOnlyFileName := readOnlyFile.Name()
	err = os.Chmod(readOnlyFileName, 0o444)
	require.NoError(t, err)

	err = Copy(srcFile.Name(), readOnlyFile.Name(), 0, 0)
	require.ErrorIs(t, err, os.ErrPermission)

	defer func() {
		err = os.Chmod(readOnlyFileName, 0o222)
		require.NoError(t, err)
		err = os.Remove(readOnlyFile.Name())
		require.NoError(t, err)
	}()
	defer func() {
		err = srcFile.Close()
		require.NoError(t, err)
	}()
}

func TestUnsupportedFile(t *testing.T) {
	const srcFileName = "test"
	srcPath = filepath.Join(testFolder, srcFileName)

	err := os.Mkdir(srcPath, os.ModePerm)
	require.NoError(t, err)

	err = Copy(srcPath, dstPath, 0, 0)
	require.ErrorIs(t, err, ErrUnsupportedFile)

	defer func() {
		err = os.Remove(srcPath)
		require.NoError(t, err)
	}()
}
