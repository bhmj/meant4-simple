package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const errInvalidInput string = `{ "error":"Incorrect input" }`
const errTwoNumbersExpected string = `{ "error":"Only two numbers expected in v1" }`

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

type inputQuery struct {
	Count   int   `json:"count"`
	Numbers []int `json:"numbers"`
}

// safeFactorial is a wrapper that checks input
func safeFactorial(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var q inputQuery
	// always return json
	w.Header().Set("Content-Type", "application/json")
	// decode input
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		http.Error(w, errInvalidInput, http.StatusBadRequest)
		return
	}
	if len(q.Numbers) != 2 {
		http.Error(w, errTwoNumbersExpected, http.StatusBadRequest)
		return
	}
	// check values
	if q.Numbers[0] <= 0 || q.Numbers[1] <= 0 {
		http.Error(w, errInvalidInput, http.StatusBadRequest)
		return
	}
	// calculate
	resA, resB := calculate(q.Numbers[0], q.Numbers[1])

	// short output: number of bits
	if q.Count > 0 {
		short := make([]int, 2)
		short[0] = resA.BitLen()
		short[1] = resB.BitLen()
		buf, err := json.Marshal(short)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(buf))
		return
	}

	// serialize response
	buf, err := json.Marshal([]*big.Int{resA, resB})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(buf))
}

func main() {
	fmt.Println("Factorial calculating service\nUsage: POST localhost:8989/factorial with body { \"numbers\": [num1, num2] }\nCtrl+C to stop")
	router := httprouter.New()
	router.POST("/factorial", safeFactorial)
	log.Fatal(http.ListenAndServe(":8989", router))
}
