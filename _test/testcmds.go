package test

import (
	"errors"
	"time"
)

// IntStr is a testing structure
type IntStr struct {
	I int
}

// CI is a testing type
type CI struct{}

// CopyIntX is a testing method
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
