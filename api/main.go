package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "super secret information")
}

var MySigningKey = []byte(os.Getenv("SECRET_KEY"))

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					return nil, fmt.Errorf("Invalid signing signature method")
				}
				return MySigningKey, nil
			}, jwt.WithAudience("billing.jwt.go.io"), jwt.WithIssuer("jwtgo.io"))

			if err != nil {
				fmt.Fprintf(w, "%s", err.Error())
				return
			}

			if token != nil && token.Valid {
				endpoint(w, r)
			} else {
				fmt.Fprintf(w, "Not Authorized")
			}
		} else {
			fmt.Fprintf(w, "No authorisation token provided")
		}
	})
}

func handleRequest() {
	http.Handle("/", isAuthorized(homePage))
	log.Fatal(http.ListenAndServe(":9001", nil))
}

func main() {
	fmt.Println("server")
	handleRequest()
} 