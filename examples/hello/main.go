package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/arturoeanton/go-r2-utils/examples/hello/ejemplo1"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, ejemplo1.Ejemplo1()+" 1")
	fmt.Println("Endpoint Hit: homePage")
}

func homePage2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, ejemplo1.Ejemplo1()+" 2")
	fmt.Println("Endpoint Hit: homePage 2!")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/2", homePage2)
	fmt.Println("Server on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	handleRequests()
}
