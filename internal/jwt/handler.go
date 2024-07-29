package jwt

import (
	"fmt"
	"net/http"

	"github.com/bignyap/verifyjwt/internal/utils"
)

func JWTVerificationHandler(w http.ResponseWriter, r *http.Request) {

	tokenString, err := ExtractToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	res, err := ParseAndVerifyJWT(tokenString)
	if err != nil {
		errMsg := fmt.Sprintf("Error while verifying token %s", err)
		w.Header().Set("realm", "")
		w.Header().Set("sub", "")
		http.Error(w, errMsg, http.StatusUnauthorized)
		return
	}

	w.Header().Set("realm", res["realm"].(string))
	w.Header().Set("sub", res["sub"].(string))

	utils.ToJson(res, w, r)
}
