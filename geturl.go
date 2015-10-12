package glick

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

// PluginGetURL only works with an api with the signature string/*string.
func PluginGetURL(static bool, uri string) Plugger {
	if static {
		if uri == "" {
			return nil
		}
	}
	return func(ctx context.Context, in interface{}) (out interface{}, err error) {
		ins, ok := in.(string)
		if !ok {
			return nil, ErrNotStr
		}
		if static {
			ins = uri
		}
		resp, err := ctxhttp.Get(ctx, http.DefaultClient, ins)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		byts, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err // unable to create a simple test case for this error
		}
		s := string(byts)
		return &s, nil
	}
}

// ConfigGetURL provides the Configurator for the URL class of plugin
func ConfigGetURL(lib *Library) error {
	return lib.AddConfigurator("URL", func(l *Library, line int, cfg *Config) error {
		var proto string
		if !(l.apim[cfg.API].in.AssignableTo(reflect.TypeOf(proto)) &&
			l.apim[cfg.API].out.AssignableTo(reflect.TypeOf(&proto))) {
			return fmt.Errorf("entry %d API %s is not of simple type (string/*string) ",
				line, cfg.API)
		}
		pi := PluginGetURL(cfg.Static, cfg.Path)
		if err := l.RegPlugin(cfg.API, cfg.Action, pi); err != nil {
			return fmt.Errorf("entry %d URL register plugin error: %v",
				line, err)
		}
		return nil
	})
}
