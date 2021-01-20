# Client API
API to create sessions, jobs and tasks. 

Execution of jobs is covered in [executor/README.md](../executor/README.md)

# Usage
Create a job with 2 tasks, the task data is passed to the engine when the job is executed.

```
session := NewSession("test-session", "bu1") 
job := NewJob("test-job-1", session)
job.CreateTask([]byte("ABC€"))
job.CreateTask([]byte("1234"))
```

Notes:
* Session name and business and are used to create topics/subscriptions, and cloud storage objects
* Given above point, session/business name combination should be unique otherwise expect errors due to duplicate resources 