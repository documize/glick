package glkit_test

import (
	"testing"

	"github.com/documize/glick"

	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// example below modified from https://github.com/go-kit/kit/blob/master/examples/stringsvc1/main.go

// StringService provides operations on strings.
type StringService interface {
	Uppercase(string) (string, error)
	Count(string) int
}

type stringService struct{}

func (stringService) Uppercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func (stringService) Count(s string) int {
	return len(s)
}

func servermain() {
	ctx := context.Background()
	svc := stringService{}

	uppercaseHandler := httptransport.NewServer(
		ctx,
		makeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
	)

	countHandler := httptransport.NewServer(
		ctx,
		makeCountEndpoint(svc),
		decodeCountRequest,
		encodeResponse,
	)

	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(uppercaseRequest)
		v, err := svc.Uppercase(req.S)
		if err != nil {
			return uppercaseResponse{v, err.Error()}, nil
		}
		return uppercaseResponse{v, ""}, nil
	}
}

func makeCountEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(countRequest)
		v := svc.Count(req.S)
		return countResponse{v}, nil
	}
}

func decodeUppercaseRequest(r *http.Request) (interface{}, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCountRequest(r *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

type uppercaseRequest struct {
	S string `json:"s"`
}

type uppercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}

type countRequest struct {
	S string `json:"s"`
}

type countResponse struct {
	V int `json:"v"`
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")

func TestGoKitStringsvc1(t *testing.T) {
	go servermain()

	r, err := http.Post("http://localhost:8080/uppercase", "application/json",
		strings.NewReader(`{"s":"hello, world"}`))
	if err != nil {
		t.Error(err)
	}
	b, err2 := ioutil.ReadAll(r.Body)
	if err2 != nil {
		t.Error(err2)
	}
	r.Body.Close()
	if string(b) != `{"v":"HELLO, WORLD"}`+"\n" {
		t.Error("/uppercase did not work: " + string(b))
	}
	r, err = http.Post("http://localhost:8080/count", "application/json",
		strings.NewReader(`{"s":"hello, world"}`))
	if err != nil {
		t.Error(err)
	}
	b, err2 = ioutil.ReadAll(r.Body)
	if err2 != nil {
		t.Error(err2)
	}
	r.Body.Close()
	if string(b) != `{"v":12}`+"\n" {
		t.Error("/count did not work: " + string(b))
	}
}

func TestAssignFn(t *testing.T) {
	var glp glick.Plugger
	var kep endpoint.Endpoint

	x := func(c context.Context, i interface{}) (interface{}, error) {
		return nil, nil
	}

	glp = x
	kep = x
	glp = glick.Plugger(kep)
	kep = endpoint.Endpoint(glp)
}
