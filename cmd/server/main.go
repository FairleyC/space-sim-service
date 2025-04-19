package main

import (
	"context"
	"fmt"

	"github.com/FairleyC/space-sim-service/internal/database"
	"github.com/FairleyC/space-sim-service/internal/services/commodity"
	"github.com/FairleyC/space-sim-service/internal/services/solarSystem"
	transport "github.com/FairleyC/space-sim-service/internal/transport/http"
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

	commodityService := commodity.NewService(db)
	solarSystemService := solarSystem.NewService(db)
	httpHandler := transport.NewHandler(commodityService, solarSystemService)
	if err := httpHandler.Serve(); err != nil {
		return err
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
