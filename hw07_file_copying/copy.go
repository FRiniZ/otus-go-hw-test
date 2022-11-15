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
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	var err_last error = nil

	in, err := os.Open (fromPath)
	if err != nil {
		return err
	}

	fi, err := in.Stat()
	if err != nil {
		return err
	}

	defer in.Close()

	out, err := os.Create (toPath)
	if err != nil {
		return err
	}

	defer out.Close()

	if offset > 0 {
	    if (offset > fi.Size()) {
		return ErrOffsetExceedsFileSize
	    }
	    in.Seek (offset, 0)
	}

	limit_r := limit
	if limit_r == 0 || (fi.Size() > 0 &&limit_r > fi.Size() - offset) {
	    limit_r = fi.Size() - offset
	}

	in_limit := io.LimitReader (in, limit_r)

	bar := pb.Full.Start64(limit_r)

	barReader := bar.NewProxyReader(in_limit)

	for {
		_, err := io.CopyN(out, barReader, 1024)
		if err != nil {
		    if err == io.EOF {
			break
		    }
		    err_last = err	
		    break;
		}
	}

	bar.Finish()

	return err_last
}
