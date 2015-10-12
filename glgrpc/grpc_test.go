package glgrpc_test

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/documize/glick"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func ConfigGRPChw(lib *glick.Library) error {
	return lib.AddConfigurator("gRPChw", func(l *glick.Library, line int, cfg *glick.Config) error {
		if err := l.RegPlugin(cfg.API, cfg.Action,
			func(ctx context.Context, in interface{}) (out interface{}, err error) {
				ins, ok := in.(*pb.HelloRequest)
				if !ok {
					return nil, errors.New("not *pb.HelloRequest")
				}
				out = interface{}(&pb.HelloReply{})
				outsp := out.(*pb.HelloReply)
				conn, err := grpc.Dial(address, grpc.WithInsecure())
				if err != nil {
					return nil, err
				}
				defer conn.Close()
				c := pb.NewGreeterClient(conn)

				r, err := c.SayHello(context.Background(), ins)
				if err != nil {
					return nil, err
				}
				*outsp = *r
				return out, nil
			}); err != nil {
			return fmt.Errorf("entry %d GRPChw register plugin error: %v",
				line, err)
		}
		return nil
	})
}

func TestGRPChw(t *testing.T) {
	gopath, ok := os.LookupEnv("GOPATH")
	if !ok {
		t.Error("can't find GOPATH")
		return
	}
	exe := gopath + "/src/google.golang.org/grpc/examples/helloworld/greeter_server/greeter_server"
	ecmd := exec.Command(exe)
	go func() {
		err := ecmd.Run()
		if err != nil {
			t.Error(err)
		}
	}()
	l := glick.New(nil)
	var req pb.HelloRequest
	var err error
	if err = l.RegAPI("hw", &req, func() interface{} { var hr pb.HelloReply; return interface{}(&hr) }, 2*time.Second); err != nil {
		t.Error(err)
	}
	ConfigGRPChw(l)
	if err := l.Config([]byte(`[
{"API":"hw","Action":"hwAct","Type":"gRPChw","Path":"` + address + `"}
		]`)); err != nil {
		t.Error(err)
	}
	req.Name = "gRPC"
	repI, err := l.Run(nil, "hw", "hwAct", &req)
	if err != nil {
		t.Error(err)
		return
	}
	if repI.(*pb.HelloReply).Message != "Hello gRPC" {
		t.Error("gRPC call did not work")
	}
}
