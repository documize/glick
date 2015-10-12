Versatile plugin framework for Go

This repository contains the "glick" plug-in framework, which is a work-in-progress.

Why "glick"? Well the framework is written in "go" and intended to be as easy to build with as lego bricks which "click" together, hence "glick".

TODO...

Besides the standard packages, "glick" relies on:
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"

Additionally, "glick/glpie" relies on:
	"github.com/natefinch/pie"

The tests in "glick/glgrpc" rely on:
	"google.golang.org/grpc"

The package "glick/glkit" relies on:
	"github.com/go-kit/kit"
	"gopkg.in/logfmt.v0"
	"gopkg.in/stack.v1"

In order to run the tests, server counterpart executables need to be built. "test/build_tests.sh" provides a bash script for doing this, it must be run from the directory it is in.

Examples of usage still need to be written, for now the tests give the best indication.

The code has only been tested on OSX.
