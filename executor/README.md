# Executor
API to execute a job. 

## Create GCP context
The GCP content holds some high level information used to help execute a job
```
gcpContext := &executor.GcpContext{
		Project:    "<project name>",
		Namespace:  "<kubernetes namespace>",
		BucketName: "<cloud storage bucket name>",
		Business:   "<business name>"}
```
## Create executor client
```
exe, err := executor.New(gcpContext)
if err != nil {
   log.Fatalf("Error creating client, %v", err)
}
```
## Run job on GKE
```
job := ...
result := exe.Execute(job)
```
see [../client/README.md](../client/README.md) for details on creating a job

## Handle result
Call to `Execute` returns a `result`, typical approach would be to handle error on result or wait for job to complete. If there is an error on the result then it was not possible to submit the job, otherise, if the job was submitted, then wait for it to complete by calling `result.Watch`, which in turn may return errors associated with running the job. Example flow:
```
	err := result.Error
	if err != nil {
		err = result.Watch()
	}

	if err != nil {
		log.Printf("%v", err)
		cxlErr := exe.Cancel(job)
		if cxlErr != nil {
			log.Printf("Error cancelling job: %v", cxlErr)
		}
	}
```

## Viewing job state
The job will transition through several states, to see current state:
```
state := job.CurrentState()
```
see [../client/job.go](../client/job.go) for details on possible states, they include: "Job Created", "Task Data Uploading", and "Job Message Published"


## Viewing job errors
To view any job errors, more specifically tasks in error:
```
if job.HasErrors() {
   for _, task := range job.TasksInError() {
     log.Printf("task errors: %v\n", task.Errors())
   }
}
```
