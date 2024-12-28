package main

import (
	"encoding/json"
	"fmt"
	"github.com/knblkxx/golang-calc/calculator"
	"io/ioutil"
	"log"
	"net/http"
)

type Request struct {
	Expression string `json:"expression"`
}
type ResponseOK struct {
	Result string `json:"result"`
}
type ResponseERROR struct {
	Error string `json:"error"`
}

func CalculatorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		var req Request
		err = json.Unmarshal(body, &req)
		result, err := calculator.Calc(req.Expression)

		defer func() {
			if r := recover(); r != nil {
				w.WriteHeader(500)
				jsonResult, _ := json.Marshal(ResponseERROR{Error: "Internal server error"})
				fmt.Fprint(w, string(jsonResult))
				log.Println("POST", req, string(jsonResult), 500)
				return
			}
		}()

		if err != nil {
			w.WriteHeader(422)
			jsonResult, _ := json.Marshal(ResponseERROR{Error: "Expression is not valid"})
			fmt.Fprint(w, string(jsonResult))
			log.Println("POST", req, string(jsonResult), 422)
			return
		}

		w.WriteHeader(200)
		jsonResult, err := json.Marshal(ResponseOK{Result: fmt.Sprintf("%f", result)})
		if err != nil {
			w.WriteHeader(500)
			jsonResult, _ = json.Marshal(ResponseERROR{Error: "Internal server error"})
			fmt.Fprint(w, string(jsonResult))
			log.Println("POST", req, string(jsonResult), 500)
			return
		}
		fmt.Fprint(w, string(jsonResult))
		log.Println("POST", req, string(jsonResult), 200)

	} else {
		w.WriteHeader(405)
		log.Println(r.Method, 405)
	}
}

func main() {
	http.HandleFunc("/api/v1/calculate", CalculatorHandler)
	http.ListenAndServe(":8080", nil)
}
