package main

import (
	"log"
	"net/rpc/jsonrpc"

	"github.com/natefinch/pie"

	"documize/pif/test"
)

func main() {
	log.SetPrefix("[plugin json log] ")

	p := pie.NewProvider()
	if err := p.Register(&test.CI{}); err != nil {
		log.Fatalf("failed to register Plugin: %s", err)
	}
	p.ServeCodec(jsonrpc.NewServerCodec)
}
