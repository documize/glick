package glkit

import (
	"fmt"

	"github.com/documize/glick"

	"github.com/go-kit/kit/endpoint"
)

// Kit provides the code for running plugins which are provided via "github.com/go-kit/kit".
type Kit struct {
	API, Action string
	EP          endpoint.Endpoint
}

// ConfigKit provides the Configurator for the GoKit class of plugin
func ConfigKit(lib *glick.Library, kits []Kit) error {
	m := make(map[string]endpoint.Endpoint)
	for _, k := range kits {
		m[k.API+"::"+k.Action] = k.EP
	}
	return lib.AddConfigurator("GoKit", func(l *glick.Library, line int, cfg *glick.Config) error {
		ep, ok := m[cfg.API+"::"+cfg.Action]
		if !ok {
			return fmt.Errorf("entry %d GoKit plugin not found for api: %s action:%s",
				line, cfg.API, cfg.Action)
		}
		if err := l.RegPlugin(cfg.API, cfg.Action, glick.Plugger(ep)); err != nil {
			return fmt.Errorf("entry %d Kit register plugin error: %v",
				line, err)
		}
		return nil
	})
}
