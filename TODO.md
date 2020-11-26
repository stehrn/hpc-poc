
If pod cant be scheduled, need to pick this up.
Otherwise just looks like job is running

Pod:
Events:
  Type     Reason            Age                 From               Message
  ----     ------            ----                ----               -------
  Warning  FailedScheduling  76s (x11 over 14m)  default-scheduler  0/1 nodes are available: 1 Insufficient cpu.



Look into quotas: https://cloud.google.com/kubernetes-engine/quotas


require (
	cloud.google.com/go/pubsub v1.8.3
	github.com/googleapis/gnostic v0.5.3 // indirect
	github.com/pkg/errors v0.9.1
	k8s.io/client-go v0.16.13 // indirect
)


require (
	cloud.google.com/go v0.38.0
	github.com/pkg/errors v0.9.1
	github.com/stehrn/hpc-poc/kubernetes v0.0.0-20201125141723-dc86a9109e59
	k8s.io/api v0.16.13
	k8s.io/apimachinery v0.16.13
)


require (
	cloud.google.com/go v0.71.0
	cloud.google.com/go/pubsub v1.8.3
	github.com/stehrn/hpc-poc/gcp v0.0.0-20201125141723-dc86a9109e59
	github.com/stehrn/hpc-poc/kubernetes v0.0.0-20201125141723-dc86a9109e59
	// avoids issue with OpenAPIv2
    k8s.io/client-go v0.17.11

)

go install github.com/stehrn/hpc-poc/monitor