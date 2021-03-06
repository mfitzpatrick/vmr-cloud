package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type httpError struct {
	error
	statusCode int
	ID         string `json:"error"`
	Desc       string `json:"desc"`
}

var (
	INVALID_VESSEL_ID  = httpError{statusCode: http.StatusBadRequest, ID: "invalid-vessel-id"}
	INVALID_VOYAGE_ID  = httpError{statusCode: http.StatusBadRequest, ID: "invalid-voyage-id"}
	INVALID_ASSIST_ID  = httpError{statusCode: http.StatusBadRequest, ID: "invalid-assist-id"}
	INVALID_RISK_ID    = httpError{statusCode: http.StatusBadRequest, ID: "invalid-risk-id"}
	IMMUTABLE_RISK     = httpError{statusCode: http.StatusBadRequest, ID: "risk-is-not-mutable"}
	RISK_TIME_REQUIRED = httpError{statusCode: http.StatusBadRequest, ID: "risk-requires-time"}
	INVALID_QSTRING    = httpError{statusCode: http.StatusBadRequest, ID: "invalid-query-string-format"}
	ENDPOINT_NOT_FOUND = httpError{statusCode: http.StatusNotFound, ID: "endpoint-not-found"}
	VOYAGE_NOT_FOUND   = httpError{statusCode: http.StatusNotFound, ID: "voyage-not-found"}
	ASSIST_NOT_FOUND   = httpError{statusCode: http.StatusNotFound, ID: "assist-not-found"}
	RISK_NOT_FOUND     = httpError{statusCode: http.StatusNotFound, ID: "risk-not-found"}
	JSON_UNMARSHAL     = httpError{statusCode: http.StatusInternalServerError, ID: "server-failed-json-parse"}
	JSON_MARSHAL       = httpError{statusCode: http.StatusInternalServerError, ID: "server-failed-json-create"}
	STORAGE_FAIL       = httpError{statusCode: http.StatusInternalServerError, ID: "failed-to-store"}
	RETRIEVAL_FAIL     = httpError{statusCode: http.StatusInternalServerError, ID: "failed-to-retrieve"}
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
