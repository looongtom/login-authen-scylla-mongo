package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"log"
	"login-user/routes"
	"net/http"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Error getting env, %v", err)
	}
	router := httprouter.New()
	router.POST("/auth/login", routes.Login)
	router.POST("/auth/register", routes.Register)
	router.GET("/auth/me", routes.GetProfile)
	router.POST("/auth/logout", routes.Logout)

	fmt.Println("Listening to port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
