package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

const HTTP_URI_BASE = "/vmr/v0"

type httpError struct {
	error
	statusCode int
	ID         string `json:"error"`
	Desc       string `json:"desc"`
}

var (
	INVALID_VESSEL_ID  = httpError{statusCode: http.StatusBadRequest, ID: "invalid-vessel-id"}
	ENDPOINT_NOT_FOUND = httpError{statusCode: http.StatusBadRequest, ID: "not-found"}
	JSON_UNMARSHAL     = httpError{statusCode: http.StatusInternalServerError, ID: "server-failed-json-parse"}
	JSON_MARSHAL       = httpError{statusCode: http.StatusInternalServerError, ID: "server-failed-json-create"}
)

func (h httpError) StringFromChain(chain error) string {
	if chain != nil {
		h.Desc = fmt.Sprintf("%v", chain.Error())
	}
	if b, err := json.Marshal(h); err != nil {
		return ""
	} else {
		return string(b)
	}
}

func (h httpError) Errorf(format string, args ...interface{}) error {
	h.error = errors.Errorf(format, args...)
	return h
}

func newVoyage(ctx context.Context, body string) ([]byte, error) {
	type request struct {
		VesselID int `json:"vessel-id"`
	}
	var req request
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		return []byte{}, JSON_UNMARSHAL.Errorf("new voyage: %v with input %s", err, body)
	}
	if req.VesselID == 0 {
		return []byte{}, INVALID_VESSEL_ID.Errorf("new voyage")
	}
	if b, err := json.Marshal(struct {
		VoyageID int `json:"voyage-id"`
	}{
		VoyageID: 1,
	}); err != nil {
		return []byte{}, JSON_MARSHAL.Errorf("new voyage: %v", err)
	} else {
		return b, nil
	}
}

func routeURL(r *http.Request) ([]byte, error) {
	type routeEntry struct {
		method, uri string
		handler     func(ctx context.Context, body string) ([]byte, error)
	}
	routes := []routeEntry{
		{http.MethodPost, HTTP_URI_BASE + "/voyage", newVoyage},
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
	http.HandleFunc(HTTP_URI_BASE, httpHandler)
}

func main() {
	log.Fatal(http.ListenAndServe(":80", nil))
}
