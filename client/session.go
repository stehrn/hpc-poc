package client

import "github.com/stehrn/hpc-poc/internal/utils"

// Session a session, can contain one or more Jobs
type Session interface {
	AddJob(job Job)
	Jobs() []Job
	ObjectPath() *ObjectPath
	Destroy()
}

// LocalSession a client side (local) session
type LocalSession struct {
	Name     string
	ID       string
	Business string
	jobs     []Job
}

// NewSession create new session
func NewSession(name, business string) *LocalSession {
	return &LocalSession{
		Name:     name,
		ID:       utils.GenerateID(),
		Business: business}
}

// NewJob create new job
func (s *LocalSession) NewJob(name string) *LocalJob {
	return NewJob(name, s)
}

// AddJob add job to session
func (s *LocalSession) AddJob(job Job) {
	s.jobs = append(s.jobs, job)
}

// Jobs return jobs for this session
func (s *LocalSession) Jobs() []Job {
	return s.jobs
}

// ObjectPath location for given session (parent of jobs)
func (s *LocalSession) ObjectPath() *ObjectPath {
	return ObjectPathForSession(s.Business, s.Name)
}

// Destroy destroy the session, closing all jobs associated with the session
func (s *LocalSession) Destroy() {
	for _, job := range s.jobs {
		job.Close()
	}
}
