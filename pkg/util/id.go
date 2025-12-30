package util

import (
	"github.com/google/uuid"
	"github.com/rs/xid"
)

func NewXID() string {
	return xid.New().String()
}

func MustNewXID() string {
	v, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return v.String()
}
