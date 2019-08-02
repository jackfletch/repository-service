A microservice for fetching a user's GitHub repositories and returning the useful information from the GitHub API.

Undocumented useful commands:

```sh
docker build -t git-repos .
docker run --publish 6060:8080 --name gitRepoService --rm git-repos
errcheck
go run main.go
godoc -http=localhost:6060
curl -s localhost:8080/user
```

```ps1
gci env:* | sort-object name
```
