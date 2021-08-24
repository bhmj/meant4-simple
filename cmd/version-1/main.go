package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const invalidInput string = `{ "error":"Incorrect input"}`

func factorial(n int, ch chan *big.Int) {
	bn := big.NewInt(int64(n))
	result := big.NewInt(1)
	one := big.NewInt(1)
	// no need to recurse
	for i := 0; i < n; i++ {
		result.Mul(result, bn)
		bn.Sub(bn, one)
	}
	ch <- result
	close(ch)
}

// calculate factorials in parallel
func calculate(a, b int) (resultA, resultB *big.Int) {
	chA := make(chan *big.Int)
	chB := make(chan *big.Int)
	go factorial(a, chA)
	go factorial(b, chB)
	for resultA = range chA {
	}
	for resultB = range chB {
	}
	return resultA, resultB
}

// Query holds input values
type Query struct {
	A int `json:"a"`
	B int `jsob:"b"`
}

// safeFactorial is a wrapper that checks input
func safeFactorial(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var q Query
	// always return json
	w.Header().Set("Content-Type", "application/json")
	// decore input
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		http.Error(w, invalidInput, http.StatusBadRequest)
		return
	}
	// check values
	if q.A <= 0 || q.B <= 0 {
		http.Error(w, invalidInput, http.StatusBadRequest)
		return
	}
	// calculate
	resA, resB := calculate(q.A, q.B)
	// response
	fmt.Fprintf(w, `{"a!":"%s", "b!":"%s"}`, resA.Text(10), resB.Text(10))
}

func main() {
	router := httprouter.New()
	router.POST("/calculate", safeFactorial)

	fmt.Println("Starting factorial calculator...")
	log.Fatal(http.ListenAndServe(":8989", router))
}
