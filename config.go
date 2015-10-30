package glick

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Config defines a line in the JSON configuration file for a glick Libarary.
type Config struct {
	API     string   // must already exist.
	Actions []string // these must be unique within the API.
	Token   string   // authorisation string to pass in the API, if it contains a Token field.
	Type    string   // the type of plugin, e.g. "RPC","URL","CMD"...
	JSON    bool     // should the plugin use JSON rather than GOB encoding, if relavent.
	Method  string   // the service method to use in the plugin, if relavent.
	Static  bool     // only used by "URL" to signal a static address.
	Path    string   // path to the end-point for "RPC" or local command for "URL".
	Args    []string // only used by "CMD", command line arguments.
	Comment string   // a place to put comments about the entry.
}

// Configurator is a type of function that allows plug-in fuctionality to the Config process.
type Configurator func(lib *Library, line int, cfg *Config) error

// AddConfigurator adds a type of configuration to the library.
func (l *Library) AddConfigurator(name string, cfg Configurator) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if _, exists := l.cfgm[name]; exists {
		return errors.New("duplicate configurator")
	}
	if cfg == nil {
		return errors.New("nil configurator")
	}
	l.cfgm[name] = cfg
	return nil
}

// Configure takes a JSON-encoded byte slice and configures the plugins for a library from it.
// NOTE: duplicate actions overload earlier versions.
func (l *Library) Configure(b []byte) error {
	var m []Config
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	for line, cfg := range m {
		if _, ok := l.apim[cfg.API]; !ok {
			return fmt.Errorf("entry %d unknown api %s ", line+1, cfg.API)
		}
		if cfgfn, ok := l.cfgm[cfg.Type]; ok {
			if err := cfgfn(l, line+1, &cfg); err != nil {
				return err
			}
		} else {
			validTypes := ""
			for t := range l.cfgm {
				validTypes += " '" + t + "',"
			}
			return fmt.Errorf("entry %d unknown config type %s (expected one of:%s)",
				line+1, cfg.Type, validTypes)
		}
	}
	return nil
}
