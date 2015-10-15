package main

import (
	"fmt"
	"log"
	"net/rpc/jsonrpc"
	"os"

	"github.com/natefinch/pie"

	test "github.com/documize/glick/_test"
)

func main() {
	log.SetPrefix("[plugin bad (json) log] ")

	p := pie.NewProvider()
	if err := p.Register(&test.CI{}); err != nil {
		log.Fatalf("failed to register Plugin: %s", err)
	}
	fmt.Fprintln(os.Stdout, "bad test")
	p.ServeCodec(jsonrpc.NewServerCodec)
}
