


To build image:
```
go mod init
gcloud builds submit --tag gcr.io/hpc-poc/monitor
```


subscriptions
|date/time|id|data|job|
|dddmmmttt|123|<link>|<link>|

jobs
|name|status|start time|completion time|duration|logs|




fail, re-submit


    job, _ := jobClient.Get(jobName, metav1.GetOptions{})

    if job.Status.Active > 0 {
    return "Job is still running"

    } else {
      if job.Status.Succeeded > 0 {
       return "Job Successful"
       } 
       return "Job failed"
    }
