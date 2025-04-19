package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/FairleyC/space-sim-service/internal/commodity"
	"github.com/FairleyC/space-sim-service/internal/solarSystem"
	"github.com/gorilla/mux"
)

type Handler struct {
	Router             *mux.Router
	CommodityService   commodity.CommodityService
	SolarSystemService solarSystem.SolarSystemService
	Server             *http.Server
}

func NewHandler(commodityService commodity.CommodityService, solarSystemService solarSystem.SolarSystemService) *Handler {
	h := &Handler{
		CommodityService:   commodityService,
		SolarSystemService: solarSystemService,
	}

	h.Router = mux.NewRouter()

	h.mapRoutes()

	h.Server = &http.Server{
		Addr:    ":8080",
		Handler: h.Router,
	}

	return h
}

var (
	API = "/api"
	V1  = "/v1"
)

func (h *Handler) mapRoutes() {
	h.Router.HandleFunc(withPath(V1, "/commodities"), h.GetCommodities).Methods("GET")
	h.Router.HandleFunc(withPath(V1, "/commodities/{id}"), h.GetCommodity).Methods("GET")
	h.Router.HandleFunc(withPath(V1, "/commodities"), h.PostCommodity).Methods("POST")
	h.Router.HandleFunc(withPath(V1, "/commodities/{id}"), h.DeleteCommodity).Methods("DELETE")

	h.Router.HandleFunc(withPath(V1, "/solarSystems"), h.GetSolarSystems).Methods("GET")
	h.Router.HandleFunc(withPath(V1, "/solarSystems/{id}"), h.GetSolarSystem).Methods("GET")
	h.Router.HandleFunc(withPath(V1, "/solarSystems"), h.PostSolarSystem).Methods("POST")
	h.Router.HandleFunc(withPath(V1, "/solarSystems/{id}"), h.DeleteSolarSystem).Methods("DELETE")
}

func (h *Handler) Serve() error {
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Println(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	h.Server.Shutdown(ctx)

	log.Println("Shutting down the server...")

	return nil
}

func withPath(version string, path string) string {
	return fmt.Sprintf("%s%s%s", API, version, path)
}
