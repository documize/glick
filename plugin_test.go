package glick_test

import (
	"errors"
	"testing"
	"time"

	"github.com/documize/glick"

	"golang.org/x/net/context"
)

func TestAPI(t *testing.T) {
	l := glick.New(nil)
	if err := l.RegAPI("z", nil, nil, time.Second); err != glick.ErrNilAPI {
		t.Error("does not return nil api error")
	}
	var dummy int
	outGood := func() interface{} { var d int; return interface{}(&d) }
	if err := l.RegAPI("z", dummy, outGood, time.Second); err != nil {
		t.Error("1st reg API returns error")
	}
	if err := l.RegAPI("z", dummy, outGood, time.Second); err != glick.ErrDupAPI {
		t.Error("does not return duplicate api error")
	}
	if _, err := l.Run(nil, "z", "unknown", dummy); err != glick.ErrNoPlug {
		t.Error("does not return no plugin")
	}
	if _, err := l.Run(nil, "unknown", "unknown", dummy); err != glick.ErrNoAPI {
		t.Error("does not return unknown api error")
	}
}

func Simp(ctx context.Context, in interface{}) (out interface{}, err error) {
	r := in.(int)
	return &r, nil
}
func outSimp() interface{} { var i int; return interface{}(&i) }

func TestSimple(t *testing.T) {
	l := glick.New(nil)
	api := "S"
	var i int
	if err := l.RegPlugin("unknown", "Test", Simp); err != glick.ErrNoAPI {
		t.Error("register plugin does not give unknown API error")
	}
	if err := l.RegAPI(api, i, outSimp, time.Second); err != nil {
		t.Error(err)
		return
	}
	if er1 := l.RegPlugin(api, "Test", Simp); er1 != nil {
		t.Error("register gives error", er1)
	}
	if ret, err := l.Run(nil, api, "Test", 42); err != nil {
		t.Error(err)
	} else {
		if *(ret.(*int)) != 42 {
			t.Error("called plugin did not work")
		}
	}
}

func TestDup(t *testing.T) {
	l := glick.New(nil)
	var d struct{}
	if er0 := l.RegAPI("A", d,
		func() interface{} { var s struct{}; return interface{}(&s) },
		time.Second); er0 != nil {
		t.Error("register API gives error")
	}
	if er1 := l.RegPlugin("A", "B", Simp); er1 != nil {
		t.Error("first entry gives error")
	}
	er2 := l.RegPlugin("A", "B", Simp)
	if er2 == nil {
		t.Error("second entry does not give error")
	}
	if er2 != glick.ErrDupPlug {
		t.Error("wrong duplicate error:", er2)
	}
}

func Tov(ctx context.Context, in interface{}) (interface{}, error) {
	t := true
	return &t, nil
}

func outTov() interface{} {
	var t bool
	return interface{}(&t)
}

func Def(ctx context.Context, in interface{}) (interface{}, error) {
	t := false
	return &t, nil
}

func outDef() interface{} {
	var t bool
	return interface{}(&t)
}

func Forever(ctx context.Context, in interface{}) (interface{}, error) {
	t := false
	for {
	}
	return &t, nil // this line is unreachable
}
func outForever() interface{} {
	var t bool
	return interface{}(&t)
}

func JustBad(ctx context.Context, in interface{}) (interface{}, error) {
	return nil, errors.New("just bad, bad, bad")
}

func outJustBad() interface{} {
	var t bool
	return interface{}(&t)
}

func TestOverloader(t *testing.T) {
	hadOvStub := Tov
	l := glick.New(func(ctx context.Context, api, act string) (context.Context, glick.Plugger, error) {
		if api == "abc" && act == "meaning-of-life" {
			return ctx, hadOvStub, nil
		}
		if api == "abc" && act == "bad" {
			return ctx, nil, errors.New("you done a bad... bad... thing")
		}
		return ctx, nil, nil
	})
	var prototype int
	if err := l.RegAPI("abc", prototype,
		func() interface{} { var b bool; return interface{}(&b) },
		time.Second); err != nil {
		t.Error(err)
		return
	}
	if err := l.RegPlugin("abc", "default", Def); err != nil {
		t.Error(err)
		return
	}
	if ret, err := l.Run(nil, "abc", "default", 1); err != nil {
		t.Error(err)
	} else {
		if *(ret.(*bool)) {
			t.Error("Overloaded function called in error")
		}
	}
	if ret, err := l.Run(nil, "abc", "meaning-of-life", 1); err != nil {
		t.Error(err)
	} else {
		if !(*(ret.(*bool))) {
			t.Error("Overloaded function not called")
		}
	}
	if err := l.RegPlugin("abc", "bad", Def); err != nil {
		t.Error(err)
		return
	}
	if _, err := l.Run(nil, "abc", "bad", 1); err == nil {
		t.Error("overloader should have errored")
		return
	}
	if err := l.RegPlugin("abc", "forever", Forever); err != nil {
		t.Error(err)
		return
	}
	ctx, can := context.WithTimeout(context.Background(), time.Millisecond)
	defer can()
	if _, err := l.Run(ctx, "abc", "forever", 1); err == nil {
		t.Error("overloader should have errored")
		return
	}
	if err := l.RegPlugin("abc", "justBad", JustBad); err != nil {
		t.Error(err)
		return
	}
	ctx, can = context.WithTimeout(context.Background(), time.Millisecond)
	defer can()
	if _, err := l.Run(ctx, "abc", "justBad", 1); err == nil {
		t.Error("overloader should have errored")
		return
	}

}
