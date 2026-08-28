package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var MySigningKey = []byte(os.Getenv("SECRET_KEY"))

func GetJwt() (string ,error){
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorised"] = true
	claims["client"]= "dhruv tiwari"
	claims["aud"] = "billing.jwt.go.io"
	claims["iss"] = "jwtgo.io"
	claims["exp"] = time.Now().Add(time.Minute * 1). Unix()
	
	tokenString, err := token.SignedString(MySigningKey)

	if err!= nil{
		fmt.Errorf("something went wrong : %s", err.Error())
		return "",err
	}

	return tokenString , nil
}

func Index(w http.ResponseWriter, r *http.Request) {
	validToken, err := GetJwt()
	fmt.Println(validToken)

	if err != nil {
		fmt.Println("failed to generate the jwt token")
	}

	fmt.Fprintf(w , validToken)
}	

func handleRequest(){
	http.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequest()
}