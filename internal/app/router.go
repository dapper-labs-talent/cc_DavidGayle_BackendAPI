package app

import (
	"cc_DavidGayle_BackendAPI/internal/app/common"
	"cc_DavidGayle_BackendAPI/internal/app/core"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/handlers"
	_ "github.com/gorilla/handlers"
	mux "github.com/gorilla/mux"
	"log"
	"net/http"
)

type svcRouter struct {
	router *mux.Router
}

type SvcRouter interface {
	Start()
}

func NewSvcRouter(handlers core.UserHandler) SvcRouter {
	router := mux.NewRouter()
	router.HandleFunc("/signup", handlers.SignUp).Methods("POST")
	router.HandleFunc("/login", handlers.Login).Methods("POST")
	router.HandleFunc("/users", IsAuthorized(handlers.GetAllUsers)).Methods("GET")
	router.HandleFunc("/users", IsAuthorized(handlers.UpdateUser)).Methods("PUT")

	return &svcRouter{
		router: router,
	}
}

func (r *svcRouter) Start() {
	fmt.Println("Server started at http://localhost:8080")

	headers := []string{"X-Requested-With", "Access-Control-Allow-Origin", "Content-Type", "x-authentication-token"}
	methods := []string{"GET", "POST", "PUT"}
	origins := []string{"*"}

	handler := handlers.CORS(handlers.AllowedHeaders(headers),
		handlers.AllowedMethods(methods),
		handlers.AllowedOrigins(origins))(r.router)

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}

func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return func(respWriter http.ResponseWriter, req *http.Request) {
		authTokenKey := "X-Authentication-Token"
		if req.Header[authTokenKey] == nil {
			json.NewEncoder(respWriter).Encode(common.SetError("No Token Found"))
			return
		}

		var mySigningKey = []byte(common.SecretKey)
		token, err := jwt.Parse(
			req.Header[authTokenKey][0],
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("encountered an error while parsing token")
				}
				return mySigningKey, nil
			},
		)

		if err != nil {
			json.NewEncoder(respWriter).Encode(common.SetError("Your authentication token has been expired."))
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if email := claims["email"]; email != nil {
				req.Header.Set("Email", email.(string))
				handler.ServeHTTP(respWriter, req)
				return
			} else {
				// This should only happen if
				json.NewEncoder(respWriter).Encode(common.SetError("Your authentication token did not include your email address."))
				return
			}
		}
		json.NewEncoder(respWriter).Encode(common.SetError("Not Authorized."))
	}
}
