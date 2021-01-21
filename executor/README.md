# Executor
API to submit a job for processing on on GCP Google Kubernetes Engine (GKE). There's a clean separation between task submission via the executor, and task execution controlled by a separate _orchestrator_ running on GKE, the executor will:

* Upload task data to a _cloud storage bucket_
* Use _cloud pub/sub_ to publish message containing location of task data 

An orchestrator process running in GKE will subscribe to the messages and process the task data. 
## Create GCP context
The `GCPContext` holds high level information used to help execute a job:
```
gcpContext := &executor.GcpContext{
		Project:    "[project name]",
		Namespace:  "[kubernetes namespace]",
		BucketName: "[cloud storage bucket name]",
		Business:   "[business name]"}
```
## Create executor client
```
exe, err := executor.New(gcpContext)
if err != nil {
   log.Fatalf("Error creating client, %v", err)
}
```
## Run job 
```
job := ...
result := exe.Execute(job)
```
see [client/README.md](../client/README.md) for details on creating a job

## Handle result
Call to `Execute` returns a `result`, typical approach would be to handle error on result or wait for job to complete. If there is an error on the result then it was not possible to submit the job, otherise, if the job was submitted, then wait for it to complete by calling `result.Watch`, which in turn may return errors associated with running the job. Example flow:
```
err := result.Error
if err != nil {
	err = result.Watch()
}

if err != nil {
	log.Printf("%v", err)
	err = exe.Cancel(job)
	if err != nil {
		log.Printf("Error cancelling job: %v", err)
	}
}
```

## Viewing job state
The job will transition through several states, to see current state:
```
state := job.CurrentState()
```
see [client/job.go](../client/job.go) for details on possible states, they include: "Job Created", "Task Data Uploading", and "Job Message Published"

## Viewing job errors
To view any job errors, more specifically tasks in error:
```
if job.HasErrors() {
   for _, task := range job.TasksInError() {
     log.Printf("Task errors: %v\n", task.Errors())
   }
}
```
