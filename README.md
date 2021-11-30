Check cleanliness of all git repositories under the current working directory
```
nastyboii
nastyboii -v
```

Exclude repos using glob patterns
```
nastyboii --exclude 'playing-with-*' -e 'drafts'
```


# Learning Go
- https://github.com/spf13/cobra
- https://go.dev/doc/tutorial/create-module

```
go mod init github.com/felixhummel/nastyboii
go get github.com/spf13/cobra/cobra
go install github.com/spf13/cobra/cobra
```
