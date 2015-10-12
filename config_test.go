package glick_test

import (
	"testing"

	"github.com/documize/glick"
	"github.com/documize/glick/test"
)

func TestConfig(t *testing.T) {
	l := glick.New(nil)
	if err := l.Config([]byte("±§~`-=_+")); err == nil {
		t.Error("did not error on rubbish")
	}
	if err := l.Config([]byte("[]")); err != nil {
		t.Error(err)
	}
	if err := l.Config([]byte(`[
{"API":"missing"}
		]`)); err == nil {
		t.Error("missing API not an error")
	}
	protoString := ""
	outProtoString := func() interface{} { var s string; return interface{}(&s) }
	outProtoInt := func() interface{} { var i int; return interface{}(&i) }

	if err := l.RegAPI("string/*string", protoString, outProtoString, 0); err != nil {
		t.Error(err)
	}
	if err := l.Config([]byte(`[
{"API":"string/*string","Action":"doIt"}
		]`)); err == nil {
		t.Error("missing Type not an error")
	}
	if err := l.Config([]byte(`[
{"API":"string/*string","Action":"pwd","Type":"CMD","Path":"pwd"}
		]`)); err != nil {
		t.Error(err)
	}
	if err := l.Config([]byte(`[
{"API":"string/*string","Action":"pwd","Type":"CMD","Path":"pwd"}
		]`)); err == nil {
		t.Error("duplicate action did not error")
	}
	if err := l.Config([]byte(`[
{"API":"string/*string","Action":"garbage","Type":"CMD","Path":"garbage"}
		]`)); err == nil {
		t.Error("garbage cmd path did not error")
	}
	if err := l.RegAPI("int/*string", 42, outProtoString, 0); err != nil {
		t.Error(err)
	}
	if err := l.Config([]byte(`[
{"API":"int/*string","Action":"badAPI","Type":"CMD","Path":"pwd"}
		]`)); err == nil {
		t.Error("unsuited API for cmd did not error")
	}
	var is test.IntStr
	if err := l.RegAPI("test", is, outProtoInt, 0); err != nil {
		t.Error(err)
	}
	if err := l.Config([]byte(`[
{"API":"test","Action":"intStr3","Type":"RPC","Path":"localhost:4242","Method":"foo.bar"}
		]`)); err != nil {
		t.Error(err)
	}
	if err := l.Config([]byte(`[
{"API":"test","Action":"intStr4","Type":"RPC","Path":"foo;;:4242"}
		]`)); err == nil {
		t.Error("unsuited endPoint not spotted")
	}
	if err := l.Config([]byte(`[
{"API":"string/*string","Action":"goodURL","Type":"URL","Path":"http://golang.org","Static":true}
		]`)); err != nil {
		t.Error(err)
	}
	if err := l.Config([]byte(`[
{"API":"string/*string","Action":"badURL","Type":"URL","Path":"","Static":true}
		]`)); err == nil {
		t.Error("unsuited URL not spotted")
	}

	if err := l.AddConfigurator("zombie", nil); err == nil {
		t.Error("nil configurator not spotted")
	}
	if err := glick.ConfigGetURL(l); err == nil {
		t.Error("duplicate configurator not spotted")
	}
}
