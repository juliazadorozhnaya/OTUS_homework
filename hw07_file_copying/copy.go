package main

import (
	"errors"
	"github.com/cheggaaa/pb/v3"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrUnexpectedEOF         = errors.New("unexpected EOF")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 {
		offset = 0
	}

	if limit < 0 {
		limit = 0
	}

	sourceFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	stat, err := sourceFile.Stat()
	if err != nil {
		return err
	}
	if stat.IsDir() {
		return ErrUnsupportedFile
	}

	fileSize := stat.Size()

	if offset >= fileSize {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 || offset+limit > fileSize {
		limit = fileSize - offset
	}

	destFile, err := os.Create(toPath)
	if err != nil {
		return err
	}

	_, err = sourceFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	bar := pb.Start64(limit)
	barReader := bar.NewProxyReader(sourceFile)
	bar.Start()

	_, err = io.CopyN(destFile, barReader, limit)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return ErrUnexpectedEOF
		}
		return err
	}
	bar.Finish()

	return nil
}
