package controllers

import (
	"encoding/json"
	"identity-provider/models"
	"identity-provider/objects"

	"net/http"

	"github.com/gorilla/mux"
)

type auhtCtrl struct {
	auth *models.AuthM
}

func InitAuth(r *mux.Router, auth *models.AuthM) {
	ctrl := &auhtCtrl{auth}
	r.HandleFunc("/authorize", ctrl.authorize).Methods("POST")
}

func (ctrl *auhtCtrl) authorize(w http.ResponseWriter, r *http.Request) {
	req_body := new(objects.AuthRequest)
	err := json.NewDecoder(r.Body).Decode(req_body)
	if err != nil {
		ValidationErrorResponse(w, err.Error())
		return
	}

	data, err := ctrl.auth.Auth(req_body.Username, req_body.Password)
	if err != nil {
		BadRequest(w, "auth failed")
	} else {
		JsonSuccess(w, data)
	}
}
