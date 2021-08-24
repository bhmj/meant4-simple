package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"sort"

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

type valuedParameters []queryParameter
type positionedParameters []queryParameter

func (pp valuedParameters) Len() int           { return len(pp) }
func (pp valuedParameters) Less(i, j int) bool { return pp[i].Number < pp[j].Number }
func (pp valuedParameters) Swap(i, j int)      { pp[i], pp[j] = pp[j], pp[i] }

func (pp positionedParameters) Len() int           { return len(pp) }
func (pp positionedParameters) Less(i, j int) bool { return pp[i].Position < pp[j].Position }
func (pp positionedParameters) Swap(i, j int)      { pp[i], pp[j] = pp[j], pp[i] }

func factorialUp(numbers []queryParameter) {
	bigOne := big.NewInt(1)
	bigHead := big.NewInt(1)
	result := big.NewInt(1)
	head := 1
	pos := 0

	for i := 0; i < numbers[len(numbers)-1].Number; i++ {
		result.Mul(result, bigHead)
		for pos < len(numbers) && numbers[pos].Number == head {
			numbers[pos].Result = new(big.Int).Set(result)
			pos++
		}
		bigHead.Add(bigHead, bigOne) // big math
		head++                       // int math
	}
}

// calculateFactorials does it wisely
func calculateFactorials(vp valuedParameters) []*big.Int {
	sort.Sort(vp)

	factorialUp(vp)

	sort.Sort(positionedParameters(vp))

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
	// decore input
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
	router := httprouter.New()
	router.POST("/calculate", handleCalculate)

	fmt.Println("Starting factorial calculator...")
	log.Fatal(http.ListenAndServe(":8989", router))
}
