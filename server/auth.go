package server

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var errInvalidKey = errors.New(ErrorInvalidKey)
var errUnauthorized = errors.New(http.StatusText(http.StatusUnauthorized))

// NewAuth init a new Auth middleware
func NewAuth(realm IRealm, opts Options) *Auth {
	a := new(Auth)
	a.opts = opts
	a.realm = realm
	a.log = createLogger("auth", opts)
	return a
}

// Auth handles request authentication
type Auth struct {
	opts  Options
	log   *logrus.Entry
	realm IRealm
}

//checkRequest check if the input is valid
func (a *Auth) checkRequest(key, id, token string) error {

	if key != a.opts.Key {
		return errInvalidKey
	}

	if id == "" {
		return errUnauthorized
	}

	client := a.realm.GetClientByID(id)

	if client == nil {
		return errUnauthorized
	}

	if len(client.GetToken()) > 0 && client.GetToken() != token {
		return errUnauthorized
	}

	return nil
}

//WSHandler return a websocket handler middleware
func (a *Auth) WSHandler(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// keys := r.URL.Query()
		// key := keys.Get("key")
		// id := keys.Get("id")
		// token := keys.Get("token")

		// err := a.checkRequest(key, id, token)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusUnauthorized)
		// 	return
		// }

		handler(w, r)
	})
}

//HTTPHandler return an HTTP handler middleware
func (a *Auth) HTTPHandler(handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		key := params["key"]
		id := params["id"]
		token := params["token"]

		err := a.checkRequest(key, id, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		handler(w, r)
	})
}
