# build the test executables
cd bad
go build
cd ../gob
go build
cd ../json
go build
cd $GOPATH/src/google.golang.org/grpc/examples/helloworld/greeter_server
go build
