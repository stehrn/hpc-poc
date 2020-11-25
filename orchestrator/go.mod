module github.com/stehrn/hpc-poc/orchestrator

go 1.13

require (
	cloud.google.com/go v0.71.0
	cloud.google.com/go/pubsub v1.8.3
	github.com/stehrn/hpc-poc/gcp v0.0.0-20201125000012-c4f96d1b4c24
	github.com/stehrn/hpc-poc/kubernetes v0.0.0-20201125000012-c4f96d1b4c24
	k8s.io/api v0.19.4
	k8s.io/apimachinery v0.19.4
	k8s.io/client-go v0.0.0-20191114101535-6c5935290e33
	sigs.k8s.io/structured-merge-diff v0.0.0-20190525122527-15d366b2352e // indirect
)
