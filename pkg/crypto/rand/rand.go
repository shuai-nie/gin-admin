package rand

import (
	"bytes"
	"errors"
	"math/rand"
)

const (
	Ldigit = 1 << iota
	LlowerCase
	LupperCase
	LlowerAndUpperCase  = LlowerCase | LupperCase
	LdingitAndLowerCase = Ldigit | LlowerCase
	LdingitAndUpperCase = Ldigit | LupperCase
	LdingitAndLetter    = Ldigit | LlowerCase | LupperCase
)

var (
	digitd            = []byte("0123456789")
	lowerdCaseLetters = []byte("abcdefghijklmnopqrstuvwxyz")
	upperdCaseLetters = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

var (
	// ErrInvalidFlag 是一个错误变量，用于表示无效的标志错误
	// 当程序遇到无效的标志时会返回此错误
	ErrInvalidFlag = errors.New("invalid flag")
)

func Random(length, flag int) (string, error) {
	if length < 1 {
		length = 6
	}

	source, err := getFlagSource(flag)
	if err != nil {
		return "", err
	}

	b, err := randomBytesMod(length, byte(len(source)))
	if err != nil {
		return "", err
	}

	b, err := randomBytesMod(length, byte(len(source)))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	for _, c := range b {
		buf.WriteByte(source[c])
	}

	return buf.String(), nil
}

func getFlagSource(flag int) ([]byte, error) {
	var source []byte
	if flag&Ldigit > 0 {
		source = append(source, digitd...)
	}

	if flag&LlowerCase > 0 {
		source = append(source, lowerdCaseLetters...)
	}

	if flag&LupperCase > 0 {
		source = append(source, upperdCaseLetters...)
	}

	sourceLen := len(source)
	if sourceLen == 0 {
		return nil, ErrInvalidFlag
	}
	return source, nil
}

func randomBytesMod(length int, mod byte) ([]byte, error) {
	b := make([]byte, length)
	max := 255 - 255%mod
	i := 0
LROOT:
	for {
		r, err := randomBytes(length + length/4)
		if err != nil {
			return nil, err
		}
		for _, c := range r {
			if c >= max {
				continue
			}

			b[i] = c % mod
			i++
			if i == length {
				break LROOT
			}
		}
	}
	return b, nil
}

func randomBytes(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
