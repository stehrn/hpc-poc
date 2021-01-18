# Client API
API to create sessions, jobs and tasks. 

Execution of jobs is covered in [executor/README.md](executor/README.md)

# Usage
Create a job with 2 tasks, the task data is passed to the engine when the job is executed.

```
session := NewSession("test-session", "bu1") 
job := NewJob("Test Job 1", session1)
job.CreateTask([]byte("ABC€"))
job.CreateTask([]byte("1234"))
```

Notes:
* Business and session name are used to create topics/subscriptions, and cloud storage objects
* Given above point, business/session name combinaton should be unique otherwise expect errors due to duplicate resources 