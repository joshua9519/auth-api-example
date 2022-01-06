package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"log"

	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

// makeIAPRequest makes a request to an application protected by Identity-Aware
// Proxy with the given audience.
func makeIAPRequest(w io.Writer, request *http.Request) error {
	// request, err := http.NewRequest("GET", "http://example.com", nil)
	audience := "406810047214-jc75s22od84i0nbgolajp5nt7kco8cee.apps.googleusercontent.com"
	ctx := context.Background()

	// client is a http.Client that automatically adds an "Authorization" header
	// to any requests made.
	client, err := idtoken.NewClient(ctx, audience, option.WithCredentialsFile("./azure-auth.json"))
	if err != nil {
			return fmt.Errorf("idtoken.NewClient: %v", err)
	}

	response, err := client.Do(request)
	if err != nil {
			return fmt.Errorf("client.Do: %v", err)
	}
	defer response.Body.Close()
	if _, err := io.Copy(w, response.Body); err != nil {
			return fmt.Errorf("io.Copy: %v", err)
	}

	return nil
}

func main() {
	request, err := http.NewRequest("GET", "https://api.josh.cts-gcp.com/ping", nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Add("Accept", "application/json")

	log.Println("Make request without authentication")
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n\n", body)

	log.Println("Make request with authentication")
	request.Header.Del("Accept")
	if err = makeIAPRequest(os.Stdout, request); err != nil {
		log.Fatal(err)
	}
	println()
}