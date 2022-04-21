package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/micro/go-micro/v2/registry/memory"
	"github.com/micro/go-micro/v2/server"
)

func TestHTTPServer(t *testing.T) {
	reg := memory.NewRegistry()

	// create server
	srv := NewServer(server.Registry(reg))

	// create server mux
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`hello world`))
	})

	// create handler
	hd := srv.NewHandler(mux)

	// register handler
	if err := srv.Handle(hd); err != nil {
		t.Error(err)
	}

	// start server
	if err := srv.Start(); err != nil {
		t.Error(err)
	}

	// lookup server
	service, err := reg.GetService(server.DefaultName)
	if err != nil {
		t.Error(err)
	}

	if len(service) != 1 {
		t.Errorf("Expected 1 service got %d: %+v", len(service), service)
	}

	if len(service[0].Nodes) != 1 {
		t.Errorf("Expected 1 node got %d: %+v", len(service[0].Nodes), service[0].Nodes)
	}

	// make request
	rsp, err := http.Get(fmt.Sprintf("http://%s", service[0].Nodes[0].Address))
	if err != nil {
		t.Error(err)
	}
	defer rsp.Body.Close()

	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		t.Error(err)
	}

	if s := string(b); s != "hello world" {
		t.Errorf("Expected response %s, got %s", "hello world", s)
	}

	// stop server
	if err := srv.Stop(); err != nil {
		t.Error(err)
	}
}