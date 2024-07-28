package p2p

import (
	"encoding/gob"
	"io"
)

type Coder interface {
	Decode(io.Reader, *RPC) error
}

type GobCode struct {
}

func (c *GobCode) Decode(r io.Reader, v *RPC) error {
	return gob.NewDecoder(r).Decode(v)
}

func (c *GobCode) Encode(r io.Writer, v *RPC) error {
	return gob.NewEncoder(r).Encode(v)
}

type DefaultCode struct {
}

func (d *DefaultCode) Decode(r io.Reader, v *RPC) error {
	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}
	v.Payload = buf[:n]
	return err
}
