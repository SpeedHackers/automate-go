package main

import (
	"log"
	"net/url"
	"os"
	"regexp"

	"github.com/SpeedHackers/automate-go/openhab"
)

func main() {

	srv := &server{}
	url, err := url.Parse(os.Getenv("OH_URL"))
	if err != nil {
		log.Fatal(err)
	}
	oldbase := url.Scheme + "://" + url.Host

	srv.Client = openhab.NewClient(url.String(),
		os.Getenv("OH_USER"),
		os.Getenv("OH_PASS"),
		false)

	srv.rew = regexp.MustCompile(oldbase)
	srv.Port = "8888"

	log.Fatal(srv.Run())
}
