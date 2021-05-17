package main

import (
	"errors"
	"go/token"
	"io"
)

var (
	ErrorsIoReadAtByTokenPositionRatIsNil = errors.New("rat is nil")
)

func IoReadAtByTokenPos(rat io.ReaderAt, begin, end token.Pos) (string, error) {
	if rat == nil {
		return "", ErrorsIoReadAtByTokenPositionRatIsNil
	}

	bs := make([]byte, end-begin)
	if _, err := rat.ReadAt(bs, int64(begin-1)); err != nil {
		return "", err
	}

	return string(bs), nil
}
