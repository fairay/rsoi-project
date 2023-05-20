package controllers

import (
	"fmt"
	"net/http"
	"statistics/utils"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Ctrl struct{}

func initControllers(r *mux.Router) {
	r.Use(utils.LogHandler)

	api1_r := r.PathPrefix("/api/v1/").Subrouter()

	ctrl := &Ctrl{}
	api1_r.HandleFunc("/all", ctrl.fetch).Methods("GET")
}

func InitRouter() *mux.Router {
	router := mux.NewRouter()
	// models := models.InitModels()

	initControllers(router)
	return router
}

func RunRouter(r *mux.Router, port uint16) error {
	c := cors.New(cors.Options{})
	handler := c.Handler(r)
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), handler)
}

func (ctrl *Ctrl) fetch(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
