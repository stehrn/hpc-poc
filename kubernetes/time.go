package kubernetes

import (
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Duration duration between two times
func Duration(start, end *v1.Time) string {
	if start.IsZero() || end.IsZero() {
		return ""
	}
	startStr, err := start.MarshalQueryParameter()
	if err != nil {
		return ""
	}
	endStr, err := end.MarshalQueryParameter()
	if err != nil {
		return ""
	}
	startTime, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return ""
	}
	endTime, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		return ""
	}
	return endTime.Sub(startTime).String()
}

// ToString get stirng representation of Time
func ToString(time *v1.Time) string {
	if time.IsZero() {
		return ""
	}
	return time.String()
}
