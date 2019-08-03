A microservice for fetching a user's GitHub repositories and returning the useful information from the GitHub API.

Undocumented useful commands:

```sh
REDIS_URL=localhost:6379 go run main.go

docker run --publish 8080:8080 --name gitRepoService --rm git-repos

errcheck
godoc -http=localhost:6060
curl -s localhost:8080/user/jackfletch
```

```ps1
gci env:* | sort-object name
```
