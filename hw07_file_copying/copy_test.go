package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	tests := []struct {
		name   string
		from   string
		to     string
		offset int64
		limit  int64
		err    error
		msg    string
	}{
		{
			name:   "copy 100Mb data",
			from:   "/dev/random",
			to:     "/dev/null",
			offset: 0,
			limit:  1024 * 1024 * 100,
			err:    nil,
			msg:    "err should be: \"nil\", but err: %q",
		},
		{
			name:   "check ErrOffsetExceedsFileSize",
			from:   "/dev/random",
			to:     "/dev/null",
			offset: 1000,
			limit:  1000,
			err:    ErrOffsetExceedsFileSize,
			msg:    "err should be: \"offset exceeds file size\", but err: %q",
		},
		{
			name:   "check ErrUnsupportedFile",
			from:   "/dev/random",
			to:     "/1.txt",
			offset: 0,
			limit:  1000,
			err:    ErrUnsupportedFile,
			msg:    "err should be: \"unsupported file\", but err: %q",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := Copy(tc.from, tc.to, tc.offset, tc.limit)
			require.Truef(t, errors.Is(err, tc.err), tc.msg, err)
		})
	}
}
