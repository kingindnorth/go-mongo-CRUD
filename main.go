package main

import (
	"fmt"
	"log"
	"net/http"

	"crud/router"
)

func main() {
	router := router.Router()
	log.Fatal(http.ListenAndServe(":5000", router))
	fmt.Println("server is running on port 5000...")
}
