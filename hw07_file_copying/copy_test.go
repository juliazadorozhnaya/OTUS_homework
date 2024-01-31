package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testFolder  = "testdata"
	srcFileName = "input.txt"
	dstFileName = "out.txt"
)

var (
	srcPath = filepath.Join(testFolder, srcFileName)
	dstPath = filepath.Join(testFolder, dstFileName)
)

func TestCopy(t *testing.T) {
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
		t.Run(tc.name, func(t *testing.T) {
			err := Copy(srcPath, dstPath, tc.offset, tc.limit)
			require.NoError(t, err)

			err = os.Remove(dstPath)
			require.NoError(t, err)
		})
	}
}

func TestOffsetIsMoreOrEqualThanFileSize(t *testing.T) {
	srcFile, err := os.Open(srcPath)
	require.NoError(t, err)
	defer srcFile.Close()

	srcFileStat, err := srcFile.Stat()
	require.NoError(t, err)

	testCases := []struct {
		offset int64
		limit  int64
	}{
		{srcFileStat.Size() + 1, 0},
		{srcFileStat.Size(), 0},
	}

	for _, tc := range testCases {
		err := Copy(srcFile.Name(), dstPath, tc.offset, tc.limit)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	}
}

func TestPermissionDeniedForDstFileIfItAlreadyExists(t *testing.T) {
	_, err := os.Create(dstPath)
	require.NoError(t, err)
	defer os.Remove(dstPath)

	err = os.Chmod(dstPath, 0o444)
	require.NoError(t, err)

	err = Copy(srcPath, dstPath, 0, 0)
	require.ErrorIs(t, err, os.ErrPermission)

	err = os.Chmod(dstPath, 0o666)
	require.NoError(t, err)
}

func TestUnsupportedFile(t *testing.T) {
	srcPath := filepath.Join(testFolder, "test")
	err := os.Mkdir(srcPath, os.ModePerm)
	require.NoError(t, err)
	defer os.Remove(srcPath)

	err = Copy(srcPath, dstPath, 0, 0)
	require.ErrorIs(t, err, ErrUnsupportedFile)
}
