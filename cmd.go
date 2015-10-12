package glick

import (
	"bytes"
	"fmt"
	"os/exec"
	"reflect"
	"strings"
	"sync"

	"golang.org/x/net/context"
)

// PluginCmd only works with an api with the signature string/*string.
// it runs the given command using the input string as stdin and putting stdout into the output string.
func PluginCmd(cmdPath string, args []string) Plugger {
	var mtx sync.Mutex // ensure we only run one command at a time
	cmdPath, e := exec.LookPath(cmdPath)
	if e != nil {
		return nil
	}
	return func(ctx context.Context, in interface{}) (interface{}, error) {
		mtx.Lock()
		defer mtx.Unlock()
		ins, ok := in.(string)
		if !ok {
			return nil, ErrNotStr
		}
		ecmd := exec.Command(cmdPath, args...)
		ecmd.Stdin = strings.NewReader(ins)
		var outBuf, errBuf bytes.Buffer
		ecmd.Stdout = &outBuf
		ecmd.Stderr = &errBuf
		err := ecmd.Run()
		s := outBuf.String()
		return &s, err
	}
}

// ConfigCmd provides the Configurator for the CMD class of plugin
func ConfigCmd(lib *Library) error {
	return lib.AddConfigurator("CMD", func(l *Library, line int, cfg *Config) error {
		var proto string
		if !(l.apim[cfg.API].in.AssignableTo(reflect.TypeOf(proto)) &&
			l.apim[cfg.API].out.AssignableTo(reflect.TypeOf(&proto))) {
			return fmt.Errorf("entry %d API %s is not of simple type (string/*string) ",
				line, cfg.API)
		}
		pi := PluginCmd(cfg.Path, cfg.Args)
		if err := l.RegPlugin(cfg.API, cfg.Action, pi); err != nil {
			return fmt.Errorf("entry %d CMD register plugin error: %v",
				line, err)
		}
		return nil
	})
}
