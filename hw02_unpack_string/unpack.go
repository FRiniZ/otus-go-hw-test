package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

var zeroRune = rune(0)

const (
	SM2BEGIN = iota
	SM2READ
	SM2KEEP
	SM2PRINT
	SM2REPEAT
	SM2END
)

var (
	lastRune      = zeroRune
	b             strings.Builder
	lastRuneEsc   = false
	lastRuneDigit = false
	lastRuneKeep  = false
)

func stateMachine2ReadRune(r rune) error {
	var err error

	switch {
	case lastRuneEsc:
		lastRuneEsc = false
		if r != '\\' && !unicode.IsDigit(r) {
			err = ErrInvalidString
			break
		}
		_, err = stateMachine2(r, SM2KEEP)
	case unicode.IsDigit(r):
		_, err = stateMachine2(r, SM2REPEAT)
	case r == '\\':
		lastRuneEsc = true
	default:
		_, _ = stateMachine2(r, SM2KEEP)
	}

	return err
}

func stateMachine2(r rune, state int) (string, error) {
	var repeat int
	var err error
	str := ""

	switch state {
	case SM2BEGIN:
		lastRune = zeroRune
		lastRuneEsc = false
		lastRuneDigit = false
		lastRuneKeep = false
		b.Reset()

	case SM2READ:
		err = stateMachine2ReadRune(r)
	case SM2KEEP:
		if lastRuneKeep {
			_, _ = stateMachine2(r, SM2PRINT)
		}
		lastRune = r
		lastRuneDigit = false
		lastRuneKeep = true
	case SM2PRINT:
		b.WriteString(string(lastRune))
		lastRuneKeep = false
	case SM2REPEAT:
		// In case the first digit
		if lastRune == zeroRune {
			err = ErrInvalidString
			break
		}
		// In case two digit
		if lastRuneDigit {
			err = ErrInvalidString
			break
		}
		lastRuneKeep = false
		lastRuneDigit = true
		repeat, err = strconv.Atoi(string(r))
		if err != nil {
			err = ErrInvalidString
			break
		}
		b.WriteString(strings.Repeat(string(lastRune), repeat))
	case SM2END:
		if lastRuneKeep {
			_, _ = stateMachine2(r, SM2PRINT)
		}
		str = b.String()
	default:
		err = ErrInvalidString
	}

	return str, err
}

func Unpack(str string) (string, error) {
	result := ""
	var err error

	_, err = stateMachine2(zeroRune, SM2BEGIN)

	for _, r := range str {
		_, err = stateMachine2(r, SM2READ)
		if err != nil {
			break
		}
	}

	if err == nil {
		result, err = stateMachine2(zeroRune, SM2END)
	}

	return result, err
}
