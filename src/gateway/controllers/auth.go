package controllers

import (
	"encoding/json"
	"fmt"
	"gateway/controllers/responses"
	"gateway/objects"
	"gateway/utils"
	"io/ioutil"
	"strings"

	"net/http"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type Token struct {
	jwt.StandardClaims
}

func newJWKs(rawJWKS string) *keyfunc.JWKS {
	jwksJSON := json.RawMessage(rawJWKS)
	jwks, err := keyfunc.NewJSON(jwksJSON)
	if err != nil {
		panic(err)
	}
	return jwks
}

func RetrieveToken(w http.ResponseWriter, r *http.Request) *Token {
	reqToken := r.Header.Get("Authorization")
	if len(reqToken) == 0 {
		responses.TokenIsMissing(w)
		return nil
	}
	splitToken := strings.Split(reqToken, "Bearer ")
	tokenStr := splitToken[1]
	jwks := newJWKs(utils.Config.RawJWKS)
	tk := &Token{}

	token, err := jwt.ParseWithClaims(tokenStr, tk, jwks.Keyfunc)
	if err != nil || !token.Valid {
		responses.JwtAccessDenied(w)
		return nil
	}
	if time.Now().Unix()-tk.ExpiresAt > 0 {
		responses.TokenExpired(w)
		return nil
	}

	return tk
}

var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token := RetrieveToken(w, r); token != nil {
			r.Header.Set("X-User-Name", token.Subject)
			next.ServeHTTP(w, r)
		}
	})
}

type authCtrl struct {
	client *http.Client
}

func InitAuth(r *mux.Router, client *http.Client) {
	ctrl := &authCtrl{client}
	r.HandleFunc("/authorize", ctrl.authorize).Methods("POST")
}

func (ctrl *authCtrl) authorize(w http.ResponseWriter, r *http.Request) {
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/authorize", utils.Config.IdentityProviderEndpoint), r.Body)

	resp, err := ctrl.client.Do(req)
	if err != nil {
		responses.InternalError(w)
		return
	}
	if resp.StatusCode == http.StatusOK {
		data := &objects.AuthResponse{}
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, data)
		responses.JsonSuccess(w, data)
	} else {
		responses.BadRequest(w, "auth failed")
	}
}
