package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseObjectPath(t *testing.T) {
	objectPath, err := ParseObjectPath("bu2/web-client-session/web-client-job-1/c048mugbigp86pbtv9sg")
	if err != nil {
		t.Fatalf("Could not parse string, error: %v", err)
	}
	assert.Equal(t, objectPath.Business, "bu2", "Unexpected Business")
	assert.Equal(t, objectPath.Session, "web-client-session", "Unexpected Session")
	assert.Equal(t, objectPath.Job, "web-client-job-1", "Unexpected Job")
	assert.Equal(t, objectPath.Task, "c048mugbigp86pbtv9sg", "Unexpected Task")
}
