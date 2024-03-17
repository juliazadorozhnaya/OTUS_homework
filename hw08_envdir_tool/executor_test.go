package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	envs := Environment{
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
	t.Run("read", func(t *testing.T) {
		env, err := ReadDir("./testdata/env")
		require.NoError(t, err)
		require.Equal(t, envs, env)
	})

	t.Run("path not found", func(t *testing.T) {
		var expected Environment

		env, err := ReadDir("./testdata/no_file")
		require.Error(t, err)
		require.Equal(t, expected, env)
	})
}
