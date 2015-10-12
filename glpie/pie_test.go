package glpie_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/documize/glick"
	"github.com/documize/glick/glpie"
	"github.com/documize/glick/test"
)

func pieSwitchTest(t *testing.T, useJSON bool) {
	l := glick.New(nil)
	if err := glpie.ConfigPIE(l); err != nil {
		t.Error(err)
		return
	}
	var proto string
	protoOut := func() interface{} {
		s := ""
		return interface{}(&s)
	}
	if err := l.RegAPI("string/&string", proto, protoOut, 10*time.Second); err != nil {
		t.Error(err)
		return
	}
	if err := l.RegPlugin("string/&string", "cmdBad",
		glpie.PluginPie(useJSON, "dingbat", "doodah", nil, protoOut)); err == nil {
		t.Error("garbage pie plugin did not fail")
		return
	}
	if _, err := l.Run(nil, "string/&string", "cmdBad", proto); err == nil {
		t.Error("bad command did not fail")
		return
	}
	api := fmt.Sprintf("API%v", useJSON)
	act := fmt.Sprintf("ACT%v", useJSON)
	tisOut := func() interface{} {
		return interface{}(&test.IntStr{})
	}
	if err := l.RegAPI(api, test.IntStr{}, tisOut, 2*time.Second); err != nil {
		t.Error(err)
		return
	}
	cmdPath := "../test/gob/gob"
	if useJSON {
		cmdPath = "../test/json/json"
	}
	if err := l.RegPlugin(api, act,
		glpie.PluginPie(useJSON, "CI.CopyIntX", cmdPath, nil, tisOut)); err != nil {
		t.Error("unable to create " + err.Error())
		return
	}

	par := test.IntStr{I: 42}
	if ret, err := l.Run(nil, api, act, par); err != nil {
		t.Error("unable to run pie " + err.Error())
	} else {
		if ret.(*test.IntStr).I != 42 {
			t.Error("pie integer copy did not work")
		}
	}
	par.I = 4
	if _, err := l.Run(nil, api, act, par); err == nil {
		t.Error("over-long pie plugin did not timeout")
	}
	if err := l.RegPlugin(api, act+"bad",
		glpie.PluginPie(true, "CI.CopyIntX", "../test/bad/bad", nil, tisOut)); err != nil {
		t.Error("unable to create " + err.Error())
		return
	}
	par.I = 0
	if _, err := l.Run(nil, api, act+"bad", par); err == nil {
		t.Error("bad pie plugin did not error")
	}
	if err := l.RegPlugin(api, act+"badder",
		glpie.PluginPie(true, "CI.CopyIntX", "../test/bad/main.go", nil, tisOut)); err != nil {
		t.Error("unable to create " + err.Error())
		return
	}
	par.I = 0
	if _, err := l.Run(nil, api, act+"badder", par); err == nil {
		t.Error("non-runnable bad pie plugin did not error")
	}
	if err := l.Config([]byte(`[
{"API":"` + api + `","Action":"intStr1","Type":"PIE","Path":"../test/gob/gob","Method":"CI.CopyIntX"}
		]`)); err != nil {
		t.Error(err)
	}
	par.I = 42
	if _, err := l.Run(nil, api, "intStr1", par); err != nil {
		t.Error("unable to run intStr1 " + err.Error())
	}
	if err := l.Config([]byte(`[
{"API":"` + api + `","Action":"intStr2","Type":"PIE"}
		]`)); err == nil {
		t.Error("unsuited end pie exe not spotted")
	}

}

func TestPie(t *testing.T) {
	pieSwitchTest(t, true)
	pieSwitchTest(t, false)
}
