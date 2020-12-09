package kubernetes

import (
	"log"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const k8Format = "2006-01-02 15:04:05 +0000 MST"
const format = "Jan 02, 2006, 15:04:05 PM (MST)"

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

// ToString get string representation of Time
func ToString(k8Time *v1.Time) string {
	if k8Time.IsZero() {
		return ""
	}
	time, err := time.Parse(k8Format, k8Time.String())
	if err != nil {
		log.Printf("Failed to parse k8Time: %s, error: %v", k8Time.String(), err)
		return k8Time.String()
	}
	return time.Format(format)
}
