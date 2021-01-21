package client

import (
	"fmt"
	"strings"
)

// ObjectPath path to an object, used when deriving object path in cloud storage, which is a function of below attributes
type ObjectPath struct {
	Business string
	Session  string
	Job      string
	Task     string
}

// ObjectPathForSession location for given job (parent of jobs)
func ObjectPathForSession(business, session string) *ObjectPath {
	return &ObjectPath{
		Business: business,
		Session:  session,
	}
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

// Parse prase a string like bu2/web-client-session/web-client-job-1/c048mugbigp86pbtv9sg into an *ObjectPath
func ParseObjectPath(object string) (*ObjectPath, error) {
	parts := strings.Split(object, "/")
	if len(parts) != 4 {
		return nil, fmt.Errorf("Could not pase %s, expected 4 parts, got %d", object, len(parts))
	}
	return &ObjectPath{
		Business: parts[0],
		Session:  parts[1],
		Job:      parts[2],
		Task:     parts[3],
	}, nil
}

// BusinessDir directory for business (appends /)
func (p *ObjectPath) BusinessDir() string {
	return fmt.Sprintf("%s/", p.Business)
}

// SessionDir directory for session (appends /)
func (p *ObjectPath) SessionDir() string {
	return fmt.Sprintf("%s/%s/", p.Business, p.Session)
}

// JobDir directory for job (appends /)
func (p *ObjectPath) JobDir() string {
	return fmt.Sprintf("%s/%s/%s/", p.Business, p.Session, p.Job)
}

func (p *ObjectPath) String() string {
	return fmt.Sprintf("%s/%s/%s/%s", p.Business, p.Session, p.Job, p.Task)
}
