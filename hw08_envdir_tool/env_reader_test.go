package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	dirName := "./testdata/env"
	expectedEnvs := Environment{
		"BAR": EnvValue{
			Value:      "bar",
			NeedRemove: false,
		},
		"EMPTY": EnvValue{
			Value:      "",
			NeedRemove: false,
		},
		"FOO": EnvValue{
			Value:      "   foo\nwith new line",
			NeedRemove: false,
		},
		"HELLO": EnvValue{
			Value:      "\"hello\"",
			NeedRemove: false,
		},
		"UNSET": EnvValue{
			Value:      "",
			NeedRemove: true,
		},
	}
	t.Run("positive tests", func(t *testing.T) {
		t.Run("read envs from dir", func(t *testing.T) {
			envs, err := ReadDir(dirName)

			require.Equal(t, expectedEnvs, envs)
			require.NoError(t, err)
		})

		t.Run("empty dir", func(t *testing.T) {
			dir := "testdata/empty"

			err := os.Mkdir(dir, os.FileMode(0o755))
			require.Nil(t, err)
			defer func() {
				err = os.Remove(dir)
				require.NoError(t, err)
			}()

			env, err := ReadDir(dir)
			require.NoError(t, err)
			require.Equal(t, Environment{}, env)
		})
	})
	t.Run("negative tests", func(t *testing.T) {
		t.Run("not existing dir", func(t *testing.T) {
			_, err := ReadDir("not_existing_dir")
			require.Error(t, err)
		})

		t.Run("invalid filename", func(t *testing.T) {
			dir := "testdata/invalid_env"

			err := os.Mkdir(dir, os.FileMode(0o755))
			require.NoError(t, err)
			defer func() {
				err = os.RemoveAll(dir)
				require.NoError(t, err)
			}()

			f, err := os.Create(path.Join(dir, "NAME=INVALID"))
			require.Nil(t, err)
			defer func() {
				err = f.Close()
				require.NoError(t, err)
			}()

			env, err := ReadDir(dir)
			require.Zero(t, env)
			require.Equal(t, ErrInvalidFilename, err)
		})
	})

}
