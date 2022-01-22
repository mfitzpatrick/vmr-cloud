package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

const HTTP_URI_BASE = "/vmr/v0"

func routeURL(r *http.Request) ([]byte, error) {
	type routeEntry struct {
		method, uri string
		handler     func(ctx context.Context, body string) ([]byte, error)
	}
	routes := []routeEntry{
		{http.MethodPost, HTTP_URI_BASE + "/voyage", postVoyage},
		{http.MethodGet, HTTP_URI_BASE + "/voyage", getVoyage},
		{http.MethodGet, HTTP_URI_BASE + "/voyage/list", listVoyage},
		{http.MethodPost, HTTP_URI_BASE + "/risk", newRisk},
		{http.MethodPost, HTTP_URI_BASE + "/assist", postAssist},
		{http.MethodGet, HTTP_URI_BASE + "/assist", getAssist},
	}
	var body string
	if b, err := ioutil.ReadAll(r.Body); err == nil {
		body = string(b)
	} // else error ignored
	for _, route := range routes {
		if route.uri == r.URL.Path && route.method == r.Method {
			return route.handler(r.Context(), body)
		}
	}
	return []byte{}, ENDPOINT_NOT_FOUND.Errorf("endpoint %s %s not found", r.Method, r.URL.Path)
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if b, err := routeURL(r); err != nil {
		var herror httpError
		if ok := errors.As(err, &herror); ok {
			w.WriteHeader(herror.statusCode)
			w.Write([]byte(herror.StringFromChain(err)))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(
				"{\"error\":\"unexpected\",\"desc\":\"Unexpected server failure: %v\"}",
				err)))
		}
	} else {
		w.Write(b)
	}
}

func init() {
	http.HandleFunc(HTTP_URI_BASE+"/", httpHandler)
}

func main() {
	log.Fatal(http.ListenAndServe(":80", nil))
}
