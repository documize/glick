package glick

import (
	"fmt"
	"net/rpc"
	"net/rpc/jsonrpc"
	"reflect"

	"golang.org/x/net/context"
)

// PluginRPC returns a type which implements the Plugger interface for making an RPC.
// The return type of this class of plugin must be a pointer.
// The plugin creates a client per call to allow services to go up-and-down between calls.
func PluginRPC(useJSON bool, serviceMethod, endPoint string, ppo ProtoPlugOut) Plugin {
	if endPoint == "" || serviceMethod == "" || endPoint == "" ||
		reflect.TypeOf(ppo()).Kind() != reflect.Ptr {
		return nil
	}
	return func(ctx context.Context, in interface{}) (out interface{}, err error) {
		var client *rpc.Client
		var errDial error
		if useJSON {
			client, errDial = jsonrpc.Dial("tcp", endPoint)
		} else {
			client, errDial = rpc.Dial("tcp", endPoint)
		}
		if errDial != nil {
			return nil, errDial
		}
		out = ppo()
		err = client.Call(serviceMethod, in, out)
		err2 := client.Close()
		if err == nil {
			err = err2
		}
		return
	}
}

// ConfigRPC provides the Configurator for the RPC class of plugin.
func ConfigRPC(lib *Library) error {
	if lib == nil {
		return ErrNilLib
	}
	return lib.AddConfigurator("RPC", func(l *Library, line int, cfg *Config) error {
		ppo := l.apim[cfg.API].ppo
		pi := PluginRPC(cfg.JSON, cfg.Method, cfg.Path, ppo)
		for _, action := range cfg.Actions {
			if err := l.RegPlugin(cfg.API, action, pi, cfg); err != nil {
				return fmt.Errorf("entry %d RPC register plugin error: %v",
					line, err)
			}
		}
		return nil
	})
}
