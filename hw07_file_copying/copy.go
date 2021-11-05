package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrWritingToFile         = errors.New("error while writing to a file")
	ErrUnexpectedEOF         = errors.New("unexpected EOF")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 {
		offset = 0
	}
	if limit < 0 {
		limit = 0
	}
	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return ErrUnsupportedFile
	}

	fileSize := fileInfo.Size()

	if offset >= fileSize {
		return ErrOffsetExceedsFileSize
	}

	targetLimit := fileSize
	if limit != 0 && limit < fileSize {
		targetLimit = limit
	}
	if offset != 0 {
		if limit != 0 && offset+limit >= fileSize {
			targetLimit = fileSize - offset
		}
		if limit == 0 {
			targetLimit = fileSize - offset
		}
	}

	fileFrom, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	if offset > 0 {
		_, err := fileFrom.Seek(offset, io.SeekStart)
		if err != nil {
			return err
		}
	}
	fileTo, err := os.Create(toPath)
	if err != nil {
		return err
	}

	defer func() {
		_ = fileFrom.Close()
	}()
	defer func() {
		_ = fileTo.Close()
	}()

	bar := pb.Start64(targetLimit)
	barReader := bar.NewProxyReader(fileFrom)
	bar.Start()

	n, err := io.CopyN(fileTo, barReader, targetLimit)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return ErrUnexpectedEOF
		}
		return err
	}
	bar.Finish()
	if n != targetLimit {
		return ErrWritingToFile
	}

	return nil
}
