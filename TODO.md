
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


35.234.156.194:8082




2020/11/27 08:58:16 Pod status is: {Succeeded [{Initialized True 0001-01-01 00:00:00 +0000 UTC 2020-11-26 19:46:56 +0000 GMT PodCompleted } {Ready False 0001-01-01 00:00:00 +0000 UTC 2020-11-26 19:46:56 +0000 GMT PodCompleted } {ContainersReady False 0001-01-01 00:00:00 +0000 UTC 2020-11-26 19:46:56 +0000 GMT PodCompleted } {PodScheduled True 0001-01-01 00:00:00 +0000 UTC 2020-11-26 19:46:56 +0000 GMT  }]    10.154.0.7 10.44.1.45 [{10.44.1.45}] 2020-11-26 19:46:56 +0000 GMT [] [{engine {nil nil &ContainerStateTerminated{ExitCode:0,Signal:0,Reason:Completed,Message:,StartedAt:2020-11-26 19:46:57 +0000 GMT,FinishedAt:2020-11-26 19:46:57 +0000 GMT,ContainerID:docker://a112b2dbc72bc085cdbf77fe69b1f45f8aea0d5b1a57bcd34f5534248d10c651,}} {nil nil nil} false 0 gcr.io/hpc-poc/engine:latest docker-pullable://gcr.io/hpc-poc/engine@sha256:94b6dd8e5588909da9041e42a43de983d08db597a626502d040c6fdc91628e46 docker://a112b2dbc72bc085cdbf77fe69b1f45f8aea0d5b1a57bcd34f5534248d10c651 0xc000613229}] Burstable []}



2020/11/27 08:58:16 Pod status is: {Pending [{PodScheduled False 0001-01-01 00:00:00 +0000 UTC 2020-11-26 20:56:21 +0000 GMT Unschedulable 0/1 nodes are available: 1 Insufficient cpu.}]      [] <nil> [] [] Burstable []}




	// for i, condition := range podStatus.Conditions {
		// 	log.Printf("condition %d type: %v", i, condition.Type)
		// 	log.Printf("condition %d status: %v", i, condition.Status)
		// 	log.Printf("condition %d reason: %v", i, condition.Reason)
		// 	log.Printf("condition %d message: %v", i, condition.Message)
		// }