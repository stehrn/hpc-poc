package client

import "fmt"

// ObjectPath path to an object
type ObjectPath struct {
	Business string
	Session  string
	Job      string
	Task     string
}

// ObjectPathForJob location for given job (parent of tasks)
func ObjectPathForJob(business, session, job string) *ObjectPath {
	return &ObjectPath{
		Business: business,
		Session:  session,
		Job:      job,
	}
}

// ObjectPathForTask location for given task (parent of tasks)
func ObjectPathForTask(job *ObjectPath, task string) *ObjectPath {
	return &ObjectPath{
		Business: job.Business,
		Session:  job.Session,
		Job:      job.Job,
		Task:     task,
	}
}

// BusinessDir directory for business
func (p *ObjectPath) BusinessDir() string {
	return fmt.Sprintf("%s/", p.Business)
}

func (p *ObjectPath) String() string {
	return fmt.Sprintf("%s/%s/%s", p.Business, p.Session, p.Job)
}
