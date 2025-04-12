package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/FairleyC/space-sim-service/internal/database"
)

// Run - is going to be responsible for
// the initialization and startup of our
// go application
func Run() error {
	fmt.Println("Starting the application...")

	db, err := database.NewDatabase(context.Background())
	if err != nil {
		fmt.Println("database.NewDatabase() error: ", err)
		return err
	}

	if err := db.Migrate(); err != nil {
		fmt.Println("database.Migrate() error: ", err)
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /commodity", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "return all defined commodities\n")
	})

	mux.HandleFunc("GET /commodity/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "return a single comment with the commodity id: %s\n", id)
	})

	mux.HandleFunc("POST /commodity", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "create a new commodity\n")
	})

	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		// separating Run() allows for us to avoid main
		// from panicking when a problem occurs and instead
		// react to the error.
		fmt.Println(err)
	}
}
