package glick_test

import (
	"fmt"
	"math/rand"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
	"testing"
	"time"

	"github.com/documize/glick"
	"github.com/documize/glick/test"
)

func TestRPC(t *testing.T) {
	tisOut := func() interface{} {
		return interface{}(&test.IntStr{})
	}

	// set up the server

	if err := rpc.Register(&test.CI{}); err != nil {
		t.Error(err.Error())
		return
	}

	for i := 0; i < 2; i++ {

		endPt := "127.0.0.1:"
		rand.Seed(time.Now().UnixNano())
		endPt += fmt.Sprintf("%d", rand.Intn(9000)+1000)

		var up sync.WaitGroup
		up.Add(1)

		// start the protocol server
		go func(ii int, ep string) {
			listener, err := net.Listen("tcp", ep)
			if err != nil {
				t.Error(err.Error())
				return
			}
			up.Done()
			for {
				conn, err := listener.Accept()
				if err != nil {
					t.Error(err.Error())
					return
				}
				if ii == 0 {
					go rpc.ServeConn(conn)
				} else {
					go jsonrpc.ServeConn(conn)
				}
			}
		}(i, endPt)

		up.Wait()
		// run the client code
		var useJSON bool
		if i > 0 {
			useJSON = true
		}

		l := glick.New(nil)

		api := "ab"
		act := "cdef"
		if err := l.RegAPI(api, test.IntStr{}, tisOut, 2*time.Second); err != nil {
			t.Error(err)
			return
		}

		if err := l.RegPlugin(api, act,
			glick.PluginRPC(useJSON, "CI.CopyIntX", endPt, tisOut)); err != nil {
			t.Error("unable to create JsonRPC " + err.Error())
			return
		}

		par := test.IntStr{I: 42}
		if ret, err := l.Run(nil, api, act, par); err != nil {
			t.Error("unable to run plugin " + err.Error())
		} else {
			if ret.(*test.IntStr).I != 42 {
				t.Error("RPC integer copy did not work")
			}
		}
		par.I = 4
		if _, err := l.Run(nil, api, act, par); err == nil {
			t.Error("over-long plugin did not timeout")
		}

		if err := l.RegPlugin(api, "bep",
			glick.PluginRPC(useJSON, "", "localhost:8080", tisOut)); err == nil {
			t.Error("able to create empty end-point method")
			return
		}

		if err := l.RegPlugin(api, "bep",
			glick.PluginRPC(useJSON, "CI.CopyIntX", "±!@£$%^&*() bad end point", tisOut)); err == nil {
			t.Error("able to create bad endpoint")
			return
		}

		if _, err := l.Run(nil, api, "bep", par); err == nil {
			t.Error("did not error on bad end-point")
		}
		if err := l.RegPlugin(api, "errEP",
			glick.PluginRPC(useJSON, "CI.CopyIntX", "localhost:9999", tisOut)); err != nil {
			t.Error(err)
			return
		}
		if _, err := l.Run(nil, api, "errEP", par); err == nil {
			t.Error("did not error on error end-point")
		}
	}
}
