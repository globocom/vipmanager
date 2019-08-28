package model

import (
	"io"

	"github.com/globocom/vipmanager/http"
)

var HTTP http.HTTP

type Req interface {
	Encode(r io.Reader) error
}

type UrlBuilder interface {
	Build(base string) (string, error)
}
