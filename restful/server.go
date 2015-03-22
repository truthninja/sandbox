package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// Server
const (
	GET  = "GET"
	POST = "POST"
)

type HandlerFunc func(w http.ResponseWriter, req *http.Request)

type Resource interface {
	Get(values url.Values) (int, interface{})
	Post(values url.Values) (int, interface{})
}

type Post405 struct{}

func (r Post405) Post(values url.Values) (int, interface{}) {
	return http.StatusNotImplemented, fmt.Sprintf("POST not implemented, your request was %v", values)
}

type APIServer struct {
	port string
}

func (a *APIServer) RESTHandler(resource Resource) HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		err := req.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusNoContent)
			fmt.Fprintf(w, "Unable to parse request, no content to server")
			return
		}

		var code = http.StatusNotFound
		var data interface{}
		switch req.Method {
		case GET:
			code, data = resource.Get(req.Form)
		case POST:
			code, data = resource.Post(req.Form)
		default:
			a.Abort(w)
			return
		}

		content, err := json.Marshal(data)
		if err != nil {
			a.Abort(w)
			return
		}
		w.WriteHeader(code)
		w.Write(content)
	}
}

func (a *APIServer) AddHandler(endpoint string, resource Resource) {
	http.HandleFunc(endpoint, a.RESTHandler(resource))
}

func (a *APIServer) Start() {
	log.Fatal(http.ListenAndServe(a.port, nil))
}

func (a *APIServer) Abort(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

// Client Resources
type Sleeby struct {
	Post405
}

func (s Sleeby) Get(values url.Values) (int, interface{}) {
	return http.StatusOK, map[string]string{"hella": "sleeby"}
}

func main() {
	appy := APIServer{":8080"}
	appy.AddHandler("/", Sleeby{})
	appy.Start()
}
