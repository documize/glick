Versatile plugin framework for Go

This repository contains the "glick" plug-in framework, which is a work-in-progress.

Why "glick"? Well the framework is written in "go" and intended to be as easy to build with as lego bricks which "click" together, hence "glick".

For a description of the package see: doc.go

For a simple example see: examples_test.go

## Dependencies

Besides the standard packages, "glick" relies on:
	"golang.org/x/net/context" and
	"golang.org/x/net/context/ctxhttp"

Additionally, "glick/glpie" provides an interface to and relies upon:
	"github.com/natefinch/pie"

The tests in "glick/glgrpc" provide example code to interface with:
	"google.golang.org/grpc"

The package "glick/glkit" provides an interface to and relies upon:
	"github.com/go-kit/kit" 
(in addition go-kit/kit requires: "gopkg.in/logfmt.v0" and "gopkg.in/stack.v1")

## Testing

In order to run the tests, server counterpart executables need to be built. "test/build_tests.sh" provides a bash script for doing this, it must be run from the directory it is in.

The code has only been tested on OSX.
