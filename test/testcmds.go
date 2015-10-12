package test

import (
	"errors"
	"time"
)

type IntStr struct {
	I int
}

type CI struct{}

func (c *CI) CopyIntX(in IntStr, out *IntStr) error {
	switch in.I {
	case 0:
		return errors.New("an error message")
	case 42:
		out.I = in.I
	default: // insert a delay
		<-time.After(time.Duration(in.I) * time.Second)
	}
	return nil
}
