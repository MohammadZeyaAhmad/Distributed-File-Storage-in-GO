package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(io.Reader) error
}

type GOBDecoder struct{}

func (dec GOBDecoder) Decode(r io.Reader, msg any) error {
	return gob.NewDecoder(r).Decode(msg)
}

type DefaultDecoder struct{}