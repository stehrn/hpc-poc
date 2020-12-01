package storage

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// ToBytes conver Location to []byte
func ToBytes(location Location) ([]byte, error) {
	bytes, err := json.Marshal(location)
	if err != nil {
		return nil, errors.Wrap(err, "Could not convert location to bytes")
	}
	return bytes, nil
}

// ToLocation convert []byte to Location
func ToLocation(bytes []byte) (Location, error) {
	var location Location
	err := json.Unmarshal(bytes, &location)
	if err != nil {
		return Location{}, errors.Wrap(err, "Could not convert bytes to location")
	}
	return location, nil
}
