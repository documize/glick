package glick

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"golang.org/x/net/context"
)

var (
	// ErrNilAPI means an API value is nil.
	ErrNilAPI = errors.New("nil api")
	// ErrDupAPI means that a duplicate name has been given for an API.
	ErrDupAPI = errors.New("duplicate api")
	// ErrNoAPI means that the name of the API was not found in the map.
	ErrNoAPI = errors.New("unknown api")
	// ErrNoPlug means that no handler function was found for the plugin.
	ErrNoPlug = errors.New("no plugin found")
)

// Plugger provides a way to call plugins,
// it has the same design as Endpoint in "github.com/go-kit/kit".
type Plugger func(ctx context.Context, in interface{}) (out interface{}, err error)

// ProtoPlugOut provides a prototype for the output of a Plugger
type ProtoPlugOut func() interface{}

type plugOut struct {
	out interface{}
	err error
}

// Overloader allows the standard system settings for an API
// to be overloaded, depending on the context passed in.
type Overloader func(ctx context.Context, api, action string, handler Plugger) (context.Context, Plugger, error)
type plugkey struct {
	api, action string
}
type plugmap map[plugkey]Plugger
type apidef struct {
	ppi     interface{}
	ppo     ProtoPlugOut
	timeout time.Duration
}
type apimap map[string]apidef
type cfgmap map[string]Configurator

// Library holds the registered API and plugin database.
type Library struct {
	pim  plugmap
	apim apimap
	cfgm cfgmap
	mtx  sync.RWMutex // mutex is for map access
	ovfn Overloader
}

// New returns an initialized Library.
func New(ov Overloader) *Library {
	lib := &Library{
		apim: make(apimap),
		pim:  make(plugmap),
		cfgm: make(cfgmap),
		ovfn: ov,
	}
	if err := ConfigCmd(lib); err != nil {
		panic(err)
	}
	if err := ConfigGetURL(lib); err != nil {
		panic(err)
	}
	if err := ConfigRPC(lib); err != nil {
		panic(err)
	}
	return lib
}

// RegAPI allows registration of a named API.
// The in/out prototype defines the type that must be passed in and out.
// The timeout gives the maximum time that a Plugin using this API may take to execute.
func (l *Library) RegAPI(api string, inPrototype interface{}, outPlugProto ProtoPlugOut, timeout time.Duration) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if inPrototype == nil || outPlugProto == nil || outPlugProto() == nil {
		return ErrNilAPI
	}
	if _, found := l.apim[api]; found {
		return ErrDupAPI
	}
	l.apim[api] = apidef{inPrototype, outPlugProto, timeout}
	return nil
}

// RegPlugin registers a Plugger to use for this action on an api.
// Duplicate actions simply overload what is there.
func (l *Library) RegPlugin(api, action string, handler Plugger) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if _, hasAPI := l.apim[api]; !hasAPI {
		return ErrNoAPI
	}
	if handler == nil {
		return ErrNoPlug
	}
	l.pim[plugkey{api, action}] = handler
	return nil
}

// Run a plugin for a given action on an API, passing data in/out.
// The library overloader function may decide from the context that a non-standard
// action should be run.
func (l *Library) Run(ctx context.Context, api, action string, in interface{}) (out interface{}, err error) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	// check api correct
	if in == nil {
		return nil, ErrNilAPI
	}
	def, ok := l.apim[api]
	if ok {
		// TODO optimise?
		if !reflect.TypeOf(in).AssignableTo(reflect.TypeOf(def.ppi)) {
			return nil, fmt.Errorf("bad api types - in: got %T want %T",
				in, def.ppi)
		}
	} else {
		return nil, ErrNoAPI
	}

	if ctx == nil || ctx == context.TODO() {
		ctx = context.Background()
	}

	handler, found := l.pim[plugkey{api, action}]

	// should this run call and overload function?
	if l.ovfn != nil {
		var ovHandler Plugger
		var ovErr error
		ctx, ovHandler, ovErr = l.ovfn(ctx, api, action, handler)
		if ovErr != nil {
			return nil, ovErr
		}
		if ovHandler != nil {
			handler = ovHandler
			found = true
		}
	}
	if !found || handler == nil {
		return nil, ErrNoPlug
	}
	reply := make(chan plugOut)
	ctxWT, cancel := context.WithTimeout(ctx, l.apim[api].timeout)
	go func() {
		defer cancel()
		var plo plugOut
		plo.out, plo.err = handler(ctxWT, in)
		reply <- plo
	}()
	select {
	case <-ctxWT.Done():
		return nil, ctxWT.Err()
	case plo := <-reply:
		if plo.err == nil && (plo.out == nil ||
			!reflect.TypeOf(plo.out).AssignableTo(
				reflect.TypeOf(def.ppo()))) { // TODO optimize?
			return nil, fmt.Errorf("bad api type - out: got %T want %T",
				plo.out, def.ppo())
		}
		return plo.out, plo.err
	}
}

// ProtoPlugOut provides the way to return a function to create the output for a plugin.
func (l *Library) ProtoPlugOut(api string) (ppo ProtoPlugOut, err error) {
	v, ok := l.apim[api]
	if !ok {
		err = errors.New("could not find api: " + api)
	}
	ppo = v.ppo
	return
}
