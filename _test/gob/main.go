package main

import (
	"log"

	"github.com/natefinch/pie"

	test "github.com/documize/glick/_test"
)

func main() {
	log.SetPrefix("[plugin gob log] ")

	p := pie.NewProvider()
	if err := p.Register(&test.CI{}); err != nil {
		log.Fatalf("failed to register Plugin: %s", err)
	}
	p.Serve()
}
