package glick

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
)

// Config defines a line in the JSON configuration file for a glick Libarary.
type Config struct {
	Plugin  string   // name of the plugin server, used to configure URL ports.
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
	if l == nil {
		return ErrNilLib
	}
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
	if l == nil {
		return ErrNilLib
	}
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

// Port returns the first port number it comes across for
// a given Plugin name in a json config file, in the form: ":9999".
func Port(configJSONpath, pluginServerName string) (string, error) {

	b, err := ioutil.ReadFile(configJSONpath)
	if err != nil {
		return "", err
	}

	var m []Config
	if err := json.Unmarshal(b, &m); err != nil {
		return "", err
	}

	for _, e := range m {
		if e.Plugin == pluginServerName {
			url, err := url.Parse(e.Path)
			if err != nil {
				return "", err
			}
			bits := strings.Split(url.Host, ":")
			if len(bits) == 2 { // ignore if no ":" in Host
				_, err = strconv.Atoi(bits[1])
				if err != nil {
					return bits[1], err
				}
				return ":" + bits[1], nil
			}
			_, err = strconv.Atoi(url.Opaque) // port could be in Opaque
			if err == nil {
				return ":" + url.Opaque, nil
			}
			switch url.Scheme { // more to go here?
			case "http":
				return ":80", nil
			case "https":
				return ":443", nil
			}
		}
	}
	return "", ErrNoAPI
}
