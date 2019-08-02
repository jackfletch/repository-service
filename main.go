// fetch github user repos
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
)

var client *github.Client

func main() {
	client = github.NewClient(nil)

	r := mux.NewRouter()
	r.HandleFunc("/", handler)
	r.HandleFunc("/user/{username}", userHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello on the container!")
}

var reGithubUsername = regexp.MustCompile("(?i)^([a-z\\d]+-)*[a-z\\d]+$")

func userHandler(w http.ResponseWriter, r *http.Request) {
	// fetch user repositories
	vars := mux.Vars(r)
	username := vars["username"]
	match := reGithubUsername.MatchString(username)
	if !match {
		fmt.Fprintf(w, "Invalid Github username: %s\n", username)
		fmt.Fprintf(os.Stderr, "ERROR: invalid github username: %s\n", username)
		return
	}

	repos, resp, err := fetchRepositories(username, client)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// copy headers from GitHub response
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// remove forked repositories
	i := 0
	for _, repo := range repos {
		if !repo.GetFork() {
			repos[i] = repo
			i++
		}
	}

	// sort repositories by stars
	repos = repos[:i]
	sort.Slice(repos, func(i, j int) bool {
		return repos[i].GetStargazersCount() > repos[j].GetStargazersCount()
	})

	// marshall json object to bytes
	json, err := json.MarshalIndent(&repos, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// output json
	fmt.Fprintln(w, string(json))
}

// fetchRepositories fetches all the public repositories of a user
func fetchRepositories(username string, client *github.Client) ([]*github.Repository, *github.Response, error) {
	repos, resp, err := client.Repositories.List(context.Background(), username, nil)
	return repos, resp, err
}
