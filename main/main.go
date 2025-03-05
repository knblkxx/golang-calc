package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/knblkxx/golang-calc/calculator"
	"github.com/knblkxx/golang-calc/calculator/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	calculator.Initial()
	err := os.MkdirAll("database", os.ModePerm)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	log.Println("Starting Server")
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/calculate", handlers.CalculatorHandler)
	r.HandleFunc("/api/v1/expressions/{id}", handlers.ExpressionAnswer)
	r.HandleFunc("/api/v1/expressions", handlers.ExpressionsList)
	http.ListenAndServe(":8080", r)
}
