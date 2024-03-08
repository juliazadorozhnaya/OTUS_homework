package main

import (
	"bufio"
	"errors"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

var ErrInvalidFilename = errors.New("invalid filename")

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dirContent, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.New("failed to read dir")
	}

	env := make(Environment)
	for _, fileInfo := range dirContent {
		if fileInfo.IsDir() {
			continue
		}

		if strings.Contains(fileInfo.Name(), "=") {
			return nil, ErrInvalidFilename
		}

		var info fs.FileInfo
		info, err = fileInfo.Info()
		if err != nil {
			return nil, err
		}
		if info.Size() == 0 {
			env[fileInfo.Name()] = EnvValue{
				Value:      "",
				NeedRemove: true,
			}
			continue
		}

		val, err := getFirstValueFromFile(path.Join(dir, fileInfo.Name()))
		if err != nil {
			return nil, err
		}

		env[fileInfo.Name()] = EnvValue{
			Value:      val,
			NeedRemove: false,
		}
	}

	return env, nil
}

func getFirstValueFromFile(fileName string) (string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return "", errors.New("failed to open file")
	}
	defer func() {
		err = f.Close()
		if err != nil {
			log.Println("failed to close file: %w", err)
		}
	}()

	buf := bufio.NewScanner(f)
	buf.Scan()
	firstLine := buf.Text()

	if err := buf.Err(); err != nil {
		return "", errors.New("err in scan")
	}

	val := strings.ReplaceAll(firstLine, "\x00", "\n")
	val = strings.TrimRight(val, " \t\n")

	return val, nil
}
