package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"sort"
	"sync"

	"github.com/julienschmidt/httprouter"
)

var (
	errInvalidInput = errors.New("incorrect input")
	errNoNumbers    = errors.New("no numbers in input")
)

type inputQuery struct {
	Numbers []int `json:"numbers"`
}

type queryParameter struct {
	Position int
	Number   int
	Result   *big.Int
}

type spanResult struct {
	HighBound int
	Result    *big.Int
}

type valuedParameters []queryParameter
type positionedParameters []queryParameter
type spanResults []spanResult

func (pp valuedParameters) Len() int           { return len(pp) }
func (pp valuedParameters) Less(i, j int) bool { return pp[i].Number < pp[j].Number }
func (pp valuedParameters) Swap(i, j int)      { pp[i], pp[j] = pp[j], pp[i] }

func (pp positionedParameters) Len() int           { return len(pp) }
func (pp positionedParameters) Less(i, j int) bool { return pp[i].Position < pp[j].Position }
func (pp positionedParameters) Swap(i, j int)      { pp[i], pp[j] = pp[j], pp[i] }

func factorialSpan(numbers []queryParameter, low, high int, target *spanResult) {
	bigOne := big.NewInt(1)
	bigHead := big.NewInt(int64(low))
	result := big.NewInt(1)
	head := low
	pos := 0
	// move pos to first number above span low bound
	for pos < len(numbers) && numbers[pos].Number < head {
		pos++
	}

	for head <= high {
		result.Mul(result, bigHead)
		for pos < len(numbers) && numbers[pos].Number == head {
			numbers[pos].Result = new(big.Int).Set(result)
			pos++
		}
		bigHead.Add(bigHead, bigOne) // big math
		head++                       // int math
	}

	*target = spanResult{HighBound: high, Result: result}
}

const Parallelism int = 8 // TODO: count CPUs or ...?

// calculateFactorials does it wisely
func calculateFactorials(vp valuedParameters) []*big.Int {
	sort.Sort(vp) // by value asc

	threads := Parallelism
	maxNumber := vp[len(vp)-1].Number
	if maxNumber < 2000 { // do not parallel if max number is not worth it
		threads = 1
	}

	results := make(spanResults, threads)

	var wg sync.WaitGroup

	lowBound := 1
	for i := 1; i <= threads; i++ {
		// spawn parallel calculations
		wg.Add(1)
		highBound := maxNumber * i / threads
		go func(low, high, pos int) { factorialSpan(vp, low, high, &results[pos-1]); wg.Done() }(lowBound, highBound, i)
		lowBound = highBound + 1
	}
	wg.Wait()

	for _, result := range results {
		for _, v := range vp {
			if v.Number > result.HighBound {
				v.Result.Mul(v.Result, result.Result)
			}
		}
	}

	sort.Sort(positionedParameters(vp)) // by position in query

	var result []*big.Int
	for _, v := range vp {
		result = append(result, v.Result)
	}
	return result
}

// handleCalculate checks input
func handleCalculate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var query inputQuery
	// always return json
	w.Header().Set("Content-Type", "application/json")
	// decode input
	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		http.Error(w, errInvalidInput.Error(), http.StatusBadRequest)
		return
	}
	// check values and create positioned parameters array
	if len(query.Numbers) == 0 {
		http.Error(w, errNoNumbers.Error(), http.StatusBadRequest)
		return
	}
	var pp []queryParameter
	for i, num := range query.Numbers {
		if num <= 0 {
			http.Error(w, errInvalidInput.Error(), http.StatusBadRequest)
			return
		}
		pp = append(pp, queryParameter{Position: i, Number: num})
	}
	// calculate
	result := calculateFactorials(pp)
	// serialize response
	buf, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(buf))
}

func main() {
	fmt.Println("Factorial calculating service\nUsage: POST localhost:8989/factorial with body { \"numbers\": [num1, num2, ..., numN] }\nCtrl+C to stop")
	router := httprouter.New()
	router.POST("/calculate", handleCalculate)
	log.Fatal(http.ListenAndServe(":8989", router))
}
