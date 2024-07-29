package router

import (
	"fmt"
	"net/http"

	"github.com/bignyap/verifyjwt/internal/jwt"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome !!")
}

func RegisterHandlers(mux *http.ServeMux) {

	mux.HandleFunc("/", RootHandler)
	mux.HandleFunc("GET /wele/verifyJWT/", jwt.JWTVerificationHandler)
}
