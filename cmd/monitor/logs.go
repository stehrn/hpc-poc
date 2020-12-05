package main

import (
	"fmt"
	"net/http"
	"strings"
)

// view logs for job or pod
// url expected to be one of:
//    /logs/job/<job name>
//    /logs/pod/<pod name>
func (ctx *handlerContext) LogsHandler(w http.ResponseWriter, r *http.Request) error {
	split := strings.Split(r.URL.Path, "/")
	objectType := split[2]
	name := split[3]

	var logs string
	var err error
	if objectType == "job" {
		logs, err = ctx.client.LogsForJob(name)
	} else if objectType == "pod" {
		logs, err = ctx.client.LogsForPod(name)
	} else {
		return fmt.Errorf("Unkown object type: %v (expeded 'job' or 'pod')", objectType)
	}

	if err != nil {
		return err
	}

	_, err = w.Write([]byte(logs))
	return err
}
