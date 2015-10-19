package glick

import (
	"fmt"
	"net/rpc"
	"net/rpc/jsonrpc"
	"net/url"
	"reflect"

	"golang.org/x/net/context"
)

// PluginRPC returns a type which implements the Plugger interface for making an RPC.
// The return type of this class of plugin must be a pointer.
func PluginRPC(useJSON bool, serviceMethod, endPoint string, ppo ProtoPlugOut) Plugin {
	if endPoint == "" || serviceMethod == "" {
		return nil
	}
	if _, err := url.Parse(endPoint); err != nil {
		return nil
	}
	if reflect.TypeOf(ppo()).Kind() != reflect.Ptr {
		return nil
	}
	var client *rpc.Client
	var errDial error
	if useJSON {
		client, errDial = jsonrpc.Dial("tcp", endPoint)
	} else {
		client, errDial = rpc.Dial("tcp", endPoint)
	}
	if errDial != nil {
		return nil
	}
	return func(ctx context.Context, in interface{}) (out interface{}, err error) {
		out = ppo()
		err = client.Call(serviceMethod, in, out)
		return
	}
}

// ConfigRPC provides the Configurator for the RPC class of plugin.
func ConfigRPC(lib *Library) error {
	return lib.AddConfigurator("RPC", func(l *Library, line int, cfg *Config) error {
		ppo := l.apim[cfg.API].ppo
		pi := PluginRPC(cfg.JSON, cfg.Method, cfg.Path, ppo)
		if err := l.RegPlugin(cfg.API, cfg.Action, pi); err != nil {
			return fmt.Errorf("entry %d RPC register plugin error: %v",
				line, err)
		}
		return nil
	})
}
