package model

import (
	"io"
)

type Req interface {
	Encode(r io.Reader) error
}

type UrlBuilder interface {
	Build(base string) (string, error)
}
