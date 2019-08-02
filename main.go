// fetch github user repos
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handler)
	// r.HandleFunc("/user/{username:(?i)^[[a-z\\d]+-]*[a-z\\d]+$}", userHandler)
	r.HandleFunc("/user/{username}", userHandler)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello on the container!")
}

// type people struct {
// 	Number int `json:"number"`
// }

// Users struct which contains an array of users
type Users struct {
	Users []User `json:"users"`
}

// User struct which contains a name a type and a list of social links
type User struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Age    int    `json:"Age"`
	Social Social `json:"social"`
}

// Social struct which contains a list of links
type Social struct {
	Facebook string `json:"facebook"`
	Twitter  string `json:"twitter"`
}

// Repo struct that contains a repository
type Repo struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Language        string `json:"language"`
	HTMLURL         string `json:"html_url"`
	Fork            bool   `json:"fork"`
	StargazersCount int    `json:"stargazers_count"`
	OpenIssues      int    `json:"open_issues_count"`
	Forks           int    `json:"forks"`
	CreatedAt       string `json:"created_at"`
	PushedAt        string `json:"pushed_at"`
	UpdatedAt       string `json:"updated_at"`
}

var reGithubUsername = regexp.MustCompile("(?i)^([a-z\\d]+-)*[a-z\\d]+$")

func userHandler(w http.ResponseWriter, r *http.Request) {
	// fetch user repositories
	vars := mux.Vars(r)
	match := reGithubUsername.MatchString(vars["username"])
	if !match {
		fmt.Fprintf(w, "Invalid Github username: %s\n", vars["username"])
		fmt.Fprintf(os.Stderr, "ERROR: invalid github username: %s\n", vars["username"])
		return
	}
	response := getGithubUserRepos(w, vars["username"])
	data, _ := ioutil.ReadAll(response.Body)

	// unmarshall bytes to json object
	var repos []Repo
	err := json.Unmarshal(data, &repos)
	if err != nil {
		fmt.Println("error:", err)
	}

	// remove forked repositories
	i := 0
	for _, repo := range repos {
		if !repo.Fork {
			repos[i] = repo
			i++
		}
	}
	repos = repos[:i]
	sort.Slice(repos, func(i, j int) bool {
		return repos[i].StargazersCount > repos[j].StargazersCount
	})

	// marshall json object to bytes
	json, err := json.MarshalIndent(&repos, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// output json
	fmt.Fprintln(w, string(json))
}

func getGithubUserRepos(w http.ResponseWriter, username string) (response *http.Response) {
	response, err := http.Get("https://api.github.com/users/" + username + "/repos")
	if err != nil {
		fmt.Fprintf(w, "The HTTP request failed with error %s\n", err)
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	return response
}
