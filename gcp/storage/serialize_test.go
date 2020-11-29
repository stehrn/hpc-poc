package storage

import (
	"reflect"
	"testing"
)

func TestEncoding(t *testing.T) {

	// convert location to bytes
	location := Location{"bucket-a", "object-1"}
	bytes, err := ToBytes(location)
	if err != nil {
		t.Errorf("ToBytes() error = %v", err)
		return
	}

	// now convert back to location, and make sure its equal to expected
	got, err := ToLocation(bytes)
	if err != nil {
		t.Errorf("ToLocation() error = %v", err)
		return
	}

	if !reflect.DeepEqual(got, location) {
		t.Errorf("ToLocation() = %v, want %v", got, location)
	}
}
