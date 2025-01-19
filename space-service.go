package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Hello, World!")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println(err.Error())
	}
}
