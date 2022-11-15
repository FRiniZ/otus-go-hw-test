package main

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	in, err := os.Open(fromPath)
	if err != nil {
		log.Println(err)
		return ErrUnsupportedFile
	}

	fi, err := in.Stat()
	if err != nil {
		return err
	}

	defer in.Close()

	out, err := os.Create(toPath)
	if err != nil {
		log.Println(err)
		return ErrUnsupportedFile
	}

	defer out.Close()

	if offset > 0 {
		if offset > fi.Size() {
			return ErrOffsetExceedsFileSize
		}
		in.Seek(offset, 0)
	}

	limitRead := limit
	if limitRead == 0 || (fi.Size() > 0 && limitRead > fi.Size()-offset) {
		limitRead = fi.Size() - offset
	}

	inLimitReader := io.LimitReader(in, limitRead)

	bar := pb.Full.Start64(limitRead)

	barReader := bar.NewProxyReader(inLimitReader)

	io.Copy(out, barReader)

	bar.Finish()

	return nil
}
