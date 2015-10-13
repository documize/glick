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
func PluginRPC(useJSON bool, serviceMethod, endPoint string, ppo ProtoPlugOut) Plugger {
	if endPoint == "" || serviceMethod == "" {
		return nil
	}
	if _, err := url.Parse(endPoint); err != nil {
		return nil
	}
	if reflect.TypeOf(ppo()).Kind() != reflect.Ptr {
		return nil
	}
	return func(ctx context.Context, in interface{}) (out interface{}, err error) {
		out = ppo()
		var client *rpc.Client
		if useJSON {
			client, err = jsonrpc.Dial("tcp", endPoint)
		} else {
			client, err = rpc.Dial("tcp", endPoint)
		}
		if err != nil {
			return nil, err
		}
		defer client.Close()
		err = client.Call(serviceMethod, in, out)
		return
	}
}

// ConfigRPC provides the Configurator for the RPC class of plugin
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
