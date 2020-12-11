package main

import (
	"fmt"
	"net/http"
	"strings"
)

// view logs for pod
// url: /logs/pod/<pod name>
func (ctx *handlerContext) LogsHandler(w http.ResponseWriter, r *http.Request) error {
	split := strings.Split(r.URL.Path, "/")
	name := split[3]
	if name == "" {
		return fmt.Errorf("Missing pod name in uri: %v", r.URL.Path)
	}

	logs, err := ctx.client.LogsForPod(name)

	if err != nil {
		return err
	}

	_, err = w.Write([]byte(logs))
	return err
}
