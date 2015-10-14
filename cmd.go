package glick

import (
	"bytes"
	"fmt"
	"os/exec"
	"sync"

	"golang.org/x/net/context"
)

// PluginCmd only works with an api with a simple Text/Text signature.
// it runs the given operating system command using the input string
// as stdin and putting stdout into the output string.
func PluginCmd(cmdPath string, args []string, model interface{}) Plugger {
	var mtx sync.Mutex // ensure we only run one command at a time
	cmdPath, e := exec.LookPath(cmdPath)
	if e != nil {
		return nil
	}
	return func(ctx context.Context, in interface{}) (interface{}, error) {
		var err error
		mtx.Lock()
		defer mtx.Unlock()
		ecmd := exec.Command(cmdPath, args...)
		ecmd.Stdin, err = TextReader(in)
		if err != nil {
			return nil, err
		}
		var outBuf, errBuf bytes.Buffer
		ecmd.Stdout = &outBuf
		ecmd.Stderr = &errBuf
		err = ecmd.Run()
		if err != nil {
			return nil, err
		}
		return TextConvert(outBuf.Bytes(), model)
	}
}

// ConfigCmd provides the Configurator for plugins that run operating system commands.
func ConfigCmd(lib *Library) error {
	return lib.AddConfigurator("CMD", func(l *Library, line int, cfg *Config) error {
		if !(IsText(l.apim[cfg.API].ppi) && IsText(l.apim[cfg.API].ppo())) {
			return fmt.Errorf("entry %d API %s is not of simple type (string/*string) ",
				line, cfg.API)
		}
		pi := PluginCmd(cfg.Path, cfg.Args, l.apim[cfg.API].ppo())
		if err := l.RegPlugin(cfg.API, cfg.Action, pi); err != nil {
			return fmt.Errorf("entry %d CMD register plugin error: %v",
				line, err)
		}
		return nil
	})
}
