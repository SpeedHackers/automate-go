package main

import (
	"log"
	"net/http"
	"os"

	"github.com/SpeedHackers/automate-go/openhab"
)

func main() {

	srv := &server{}
	srv.Client = openhab.NewClient(os.Getenv("OH_URL"),
		os.Getenv("OH_USER"),
		os.Getenv("OH_PASS"),
		false)

	log.Fatal(http.ListenAndServe(":8888", srv.setupRoutes()))
}
