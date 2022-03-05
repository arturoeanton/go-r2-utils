package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/arturoeanton/go-r2-utils/cmd/hello/ejemplo1"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, ejemplo1.Ejemplo1())
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	fmt.Println("Server on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	handleRequests()
}
