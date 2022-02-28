package main

import (
	"encoding/base64"
	"flag"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	// Parse arguments
	//
	hostname := flag.String("host", "localhost", "Router hostname")
	username := flag.String("username", "admin", "Router username")
	password := flag.String("password", "admin", "Router password")
	flag.Parse()

	reboot(hostname, username, password)
}

// Send reboot command to orbi router
//
func reboot(hostname *string, username *string, password *string) {

	// Create HTTP client with timeout
	//
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Get re-boot page
	//
	httpurl := "http://" + *hostname + "/reboot.htm"
	log.Println("GET", httpurl)
	request, err := http.NewRequest(http.MethodGet, httpurl, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(*username+":"+*password)))
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Find the form
	//
	// <form method="POST" action="/apply.cgi?/reboot_waiting.htm timestamp=488450730402957"
	//
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	val, exists := doc.Find("form").First().Attr("action")
	if exists {

		httpurl := "http://" + *hostname + strings.ReplaceAll(val, " ", "%20")
		log.Println("POST", httpurl)

		data := url.Values{}
		data.Set("submit_flag", "reboot")
		data.Set("yes", "Yes")

		request, err := http.NewRequest(http.MethodPost, httpurl, strings.NewReader(data.Encode()))
		if err != nil {
			log.Fatal(err)
		}
		request.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(*username+":"+*password)))
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		response, err := client.Do(request)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(response.Status)

		defer response.Body.Close()
	} else {
		log.Fatal("POST Form not found")
	}

}
