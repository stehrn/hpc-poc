# Go(lang) tips

# Init module
After creating a new package, creat modules via: `go mod init`

# Pull in updated deps
If a dependeant module has been updated, pull update into package via:
`go get github.com/stehrn/hpc-poc/kubernetes@main`

This will pull update to kubernetes into `go.mod` (`go.sum` should also be updated)

# Run main
`go run main.go`