package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestMain(t *testing.T) {
	tests := []struct {
		Input  string
		Code   int
		Output string
	}{
		// different invalid inputs
		{``, 400, invalidInput},
		{`{}`, 400, invalidInput},
		{`foo`, 400, invalidInput},
		{`{"a":0, "b":1}`, 400, invalidInput},
		{`{"a":1, "b":-1}`, 400, invalidInput},
		{`{"a":"1", "b":"1"}`, 400, invalidInput}, // we want numbers in input
		// valid inputs
		{`{"a":1, "b":1}`, 200, `{"a!":"1", "b!":"1"}`},
		{`{"a":5, "b":10}`, 200, `{"a!":"120", "b!":"3628800"}`},
		// large values (results taken from https://onlinemschool.com/math/formula/factorial_table/)
		{`{"a":19, "b":20}`, 200, `{"a!":"121645100408832000", "b!":"2432902008176640000"}`},
		{`{"a":49, "b":50}`, 200, `{"a!":"608281864034267560872252163321295376887552831379210240000000000", "b!":"30414093201713378043612608166064768844377641568960512000000000000"}`},
	}

	router := httprouter.New()
	router.POST("/calculate", safeFactorial)

	for itst, tst := range tests {

		reader := strings.NewReader(tst.Input)
		req, err := http.NewRequest("POST", "/calculate", reader)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// check http status code
		if status := rr.Code; status != tst.Code {
			t.Errorf("test #%d: handler returned wrong status code: got %d, want %d",
				itst, status, tst.Code)
		}

		// check response body
		if strings.TrimSpace(rr.Body.String()) != tst.Output {
			t.Errorf("test #%d: handler returned unexpected body: got >%s<, want >%s<",
				itst, rr.Body.String(), tst.Output)
		}
	}
}
