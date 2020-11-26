# Go(lang) tips

# Init module
After creating a new package, creat modules via: `go mod init`

# Pull in updated deps
If a dependeant module has been updated, pull update into package via:
`go get github.com/stehrn/hpc-poc/kubernetes@main`

This will pull update to kubernetes into `go.mod` (`go.sum` should also be updated)

# Run main
`go run main.go`

# Get valid version of libs
Get version of k8 
```
kubectl version
```
Add correct libs to `go.mod` (here, version is '1.16.3'):
```
go get k8s.io/client-go@kubernetes-1.16.13
```
(run from base dir with `mod.go` in)
