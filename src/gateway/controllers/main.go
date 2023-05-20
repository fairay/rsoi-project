package controllers

import (
	"fmt"
	"gateway/models"
	"gateway/utils"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func initControllers(r *mux.Router, models *models.Models) {
	r.Use(utils.LogHandler)
	r.Use(func (next http.Handler) http.Handler {
		return utils.RequestStatMiddleware(next, "none", models.Kafka.KafkaTopic, models.Kafka.Producer);
	})
	api1_r := r.PathPrefix("/api/v1/").Subrouter()

	InitAuth(api1_r, models.Client)
	api1_r_auth := api1_r.NewRoute().Subrouter()
	api1_r_auth.Use(JwtAuthentication)

	InitFlights(api1_r_auth, models.Flights)
	InitPrivileges(api1_r_auth, models.Privileges)
	InitTickets(api1_r_auth, models.Tickets)
}

func InitRouter() *mux.Router {
	router := mux.NewRouter()
	models := models.InitModels()

	initControllers(router, models)
	return router
}

func RunRouter(r *mux.Router, port uint16) error {
	c := cors.New(cors.Options{})
	handler := c.Handler(r)
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), handler)
}
