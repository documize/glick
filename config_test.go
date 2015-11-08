package glick_test

import (
	"testing"

	"github.com/documize/glick"
	test "github.com/documize/glick/_test"
)

func TestConfig(t *testing.T) {
	l, ne := glick.New(nil)
	if ne != nil {
		t.Error(ne)
	}
	if err := l.Configure([]byte("±§~`-=_+")); err == nil {
		t.Error("did not error on rubbish")
	}
	if err := l.Configure([]byte("[]")); err != nil {
		t.Error(err)
	}
	if err := l.Configure([]byte(`[
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
	if err := l.Configure([]byte(`[
{"API":"string/*string","Actions":["doIt"]}
		]`)); err == nil {
		t.Error("missing Type not an error")
	}
	if err := l.Configure([]byte(`[
{"API":"string/*string","Actions":["pwd"],"Type":"CMD","Cmd":"pwd"}
		]`)); err != nil {
		t.Error(err)
	}
	if err := l.Configure([]byte(`[
{"API":"string/*string","Actions":["garbage"],"Type":"CMD","Cmd":"garbage"}
		]`)); err == nil {
		t.Error("garbage cmd path did not error")
	}
	if err := l.RegAPI("int/*string", 42, outProtoString, 0); err != nil {
		t.Error(err)
	}
	if err := l.Configure([]byte(`[
{"API":"int/*string","Actions":["badAPI"],"Type":"CMD","Cmd":"pwd"}
		]`)); err == nil {
		t.Error("unsuited API for cmd did not error")
	}
	if err := l.Configure([]byte(`[
{"API":"int/*string","Actions":["badAPI"],"Type":"URL","Path":"pwd"}
		]`)); err == nil {
		t.Error("unsuited API for URL did not error")
	}
	var is test.IntStr
	if err := l.RegAPI("test", is, outProtoInt, 0); err != nil {
		t.Error(err)
	}
	if err := l.Configure([]byte(`[
{"API":"test","Actions":["intStr3"],"Type":"RPC","Path":"localhost:4242","Method":"foo.bar","Token":"ABC"}
		]`)); err != nil {
		t.Error(err)
	}
	if l.Token("test", "intStr3") != "ABC" {
		t.Error("Token value not set and retrieved")
	}
	if err := l.Configure([]byte(`[
{"API":"test","Actions":["intStr4"],"Type":"RPC","Path":"foo;;:4242"}
		]`)); err == nil {
		t.Error("unsuited endPoint not spotted")
	}
	if err := l.Configure([]byte(`[
{"API":"string/*string","Actions":["goodURL"],"Type":"URL","Path":"http://golang.org","Static":true}
		]`)); err != nil {
		t.Error(err)
	}
	if err := l.Configure([]byte(`[
{"API":"string/*string","Actions":["badURL"],"Type":"URL","Path":"","Static":true}
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
