package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
)

const invalidInput = "incorrect input"
const noNumbers = "no numbers in input"

func TestMain(t *testing.T) {
	tests := []struct {
		Input  string
		Code   int
		Output string
	}{
		// different invalid inputs
		{``, 400, invalidInput},
		{`{}`, 400, noNumbers},
		{`foo`, 400, invalidInput},
		{`{"numbers":[0,1]}`, 400, invalidInput},
		{`{"numbers":[1,-1]}`, 400, invalidInput},
		{`{"numbers":["1","1"]}`, 400, invalidInput}, // we want numbers in input
		// valid inputs
		{`{"numbers":[1,1]}`, 200, `[1,1]`},
		{`{"numbers":[10,5]}`, 200, `[3628800,120]`},
		// large values (results taken from https://onlinemschool.com/math/formula/factorial_table/)
		{`{"numbers":[19, 20]}`, 200, `[121645100408832000,2432902008176640000]`},
		{`{"numbers":[49, 50]}`, 200, `[608281864034267560872252163321295376887552831379210240000000000,30414093201713378043612608166064768844377641568960512000000000000]`},
	}

	router := httprouter.New()
	router.POST("/calculate", handleCalculate)

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
